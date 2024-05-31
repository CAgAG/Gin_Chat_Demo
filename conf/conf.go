package conf

import (
	"Chat_demo/model"
	"context"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

var (
	MongoDBClient *mongo.Client

	AppMode    string
	HttpPort   string
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string

	MongoDBName string
	MongoDBAddr string
	MongoDBUser string
	MongoDBPwd  string
	MongoDBPort string
)

func Init() {
	// 读取配置
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		panic(err)
	}
	LoadServer(file)
	LoadMysql(file)
	LoadMongoDB(file)

	if AppMode == "release" {
		logging.SetLevel(logging.ErrorLevel)
		gin.SetMode(gin.ReleaseMode)
	}
	logfile, _ := os.OpenFile("./logs/logrus.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	logging.SetOutput(logfile) // 输出到文件里
	logging.Info("start logging")

	NewMongoDB()
	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	model.NewDatabase(path)
}

func NewMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBAddr + ":" + MongoDBPort)
	var err error
	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	err = MongoDBClient.Ping(context.TODO(), nil)
	if err != nil {
		logging.Info(err)
	}
}

func LoadMongoDB(file *ini.File) {
	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBUser = file.Section("MongoDB").Key("MongoDBUser").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}
