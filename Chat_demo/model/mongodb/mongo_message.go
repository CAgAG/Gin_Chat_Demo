package mongodb

import (
	"Chat_demo/conf"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"strconv"
	"time"
)

type Mongo_Trainer struct {
	Content string `bson:"content"` // 内容
	Created int64  `bson:"created"` // 创建时间
	EndTime int64  `bson:"endTime"` // 过期时间
	Read    uint   `bson:"read"`    // 已读
}

type Mongo_Full_Trainer struct {
	ID      string `bson:"_id"`     // 内容
	Content string `bson:"content"` // 内容
	Created int64  `bson:"created"` // 创建时间
	EndTime int64  `bson:"endTime"` // 过期时间
	Read    uint   `bson:"read"`    // 已读
}

type MongoResult struct {
	ID        string
	Created   int64
	Msg       string
	Content   interface{}
	Direction string
}

func MongoInsert(id string, content string, expire int64, read uint) error {
	database := conf.MongoDBName
	collection := conf.MongoDBClient.Database(database).Collection(id)
	comment := Mongo_Trainer{
		Content: content,
		Created: time.Now().Unix(),
		EndTime: time.Now().Unix() + expire,
		Read:    read,
	}
	_, err := collection.InsertOne(context.TODO(), comment)
	return err
}

func FindHis(chat_id1 string, chat_id2 string, skip_count, limit_count int) ([]MongoResult, error) {
	database := conf.MongoDBName

	var resultsSend []Mongo_Full_Trainer
	var resultsReceive []Mongo_Full_Trainer

	sendIdCollection := conf.MongoDBClient.Database(database).Collection(chat_id1)
	ReceiveIdCollection := conf.MongoDBClient.Database(database).Collection(chat_id2)

	// 在MongoDB中，可以使用skip()和limit()方法进行分页查询。skip()方法用于跳过指定数量的数据，limit()方法用于限制返回的数据量。
	// db.collection.find().skip(跳过的数量).limit(每页显示的数量)

	// 如果不知道该使用什么context，可以通过context.TODO() 产生context
	// opts := options.Find().SetSort(bson.D{{"created", -1}})
	// opts2 := options.Find().SetSkip(int64(skip_count)).SetLimit(int64(limit_count))

	opts := bson.D{}
	opts2 := options.Find()

	sendIdTimeCursor, err := sendIdCollection.Find(context.TODO(), opts, opts2)
	ReceiveIdTimeCursor, err := ReceiveIdCollection.Find(context.TODO(), opts, opts2)

	err = sendIdTimeCursor.All(context.TODO(), &resultsSend)       // sendId 对面发过来的
	err = ReceiveIdTimeCursor.All(context.TODO(), &resultsReceive) // 发给对面的, 或者说接收到的
	results, err := MergeAndSort(resultsSend, resultsReceive, chat_id1, chat_id2)
	// return results, err

	return results[skip_count:min(len(results), skip_count+limit_count)], err
}

// 读取所有未读的消息
func FindHisUnread(chat_id1 string, chat_id2 string) ([]MongoResult, error) {
	database := conf.MongoDBName

	var resultsSend []Mongo_Full_Trainer
	var resultsReceive []Mongo_Full_Trainer

	sendIdCollection := conf.MongoDBClient.Database(database).Collection(chat_id1)
	ReceiveIdCollection := conf.MongoDBClient.Database(database).Collection(chat_id2)

	opts := bson.D{{"read", 0}}
	opts2 := options.Find()

	sendIdTimeCursor, err := sendIdCollection.Find(context.TODO(), opts, opts2)
	ReceiveIdTimeCursor, err := ReceiveIdCollection.Find(context.TODO(), opts, opts2)

	err = sendIdTimeCursor.All(context.TODO(), &resultsSend)       // sendId 对面发过来的
	err = ReceiveIdTimeCursor.All(context.TODO(), &resultsReceive) // 发给对面的, 或者说接收到的
	results, err := MergeAndSort(resultsSend, resultsReceive, chat_id1, chat_id2)

	return results, err
}

func MergeAndSort(results1, results2 []Mongo_Full_Trainer, chat_id1, chat_id2 string) ([]MongoResult, error) {
	rets := make([]MongoResult, len(results1)+len(results2))
	rets_i := 0

	for _, rt := range results1 {
		ret := MongoResult{
			ID:        rt.ID,
			Direction: chat_id1,
			Created:   rt.Created,
			Content:   rt.Content,
			Msg:       strconv.Itoa(int(rt.Read)),
		}
		rets[rets_i] = ret
		rets_i++
	}
	for _, rt := range results2 {
		ret := MongoResult{
			ID:        rt.ID,
			Direction: chat_id2,
			Created:   rt.Created,
			Content:   rt.Content,
			Msg:       strconv.Itoa(int(rt.Read)),
		}
		rets[rets_i] = ret
		rets_i++
	}
	// 排序, 按创建时间
	sort.Slice(rets, func(i, j int) bool { return rets[i].Created < rets[j].Created })
	return rets, nil
}

// 修改消息 ID为已读
func Message_Read(chat_id1, chat_id2, message_id string) error {
	database := conf.MongoDBName

	sendIdCollection := conf.MongoDBClient.Database(database).Collection(chat_id1)
	ReceiveIdCollection := conf.MongoDBClient.Database(database).Collection(chat_id2)

	// 获取 object id
	object_id, _ := primitive.ObjectIDFromHex(message_id)
	opts := bson.D{{"_id", object_id}}
	// 更新操作
	opts2 := bson.D{{"$set", bson.D{{"read", 1}}}}

	var err error
	_, err = sendIdCollection.UpdateOne(context.TODO(), opts, opts2, options.Update().SetUpsert(true))
	_, err = ReceiveIdCollection.UpdateOne(context.TODO(), opts, opts2, options.Update().SetUpsert(true))
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 删除消息 ID
func Message_Del(chat_id1, chat_id2, message_id string) error {
	database := conf.MongoDBName

	sendIdCollection := conf.MongoDBClient.Database(database).Collection(chat_id1)
	ReceiveIdCollection := conf.MongoDBClient.Database(database).Collection(chat_id2)

	// 获取 object id
	object_id, _ := primitive.ObjectIDFromHex(message_id)
	opts := bson.D{{"_id", object_id}}

	var err1, err2 error
	_, err1 = sendIdCollection.DeleteOne(context.TODO(), opts)
	_, err2 = ReceiveIdCollection.DeleteOne(context.TODO(), opts)

	if err1 == nil || err2 == nil {
		return nil
	} else {
		if err1 != nil {
			return err1
		}
		return err2
	}
}
