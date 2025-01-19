package service

import (
	"Chat_demo/cache"
	"Chat_demo/model"
	"Chat_demo/model/mongodb"
	"Chat_demo/pkg/status"
	"Chat_demo/pkg/utils"
	"Chat_demo/serializer"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 结构体定义 ===================================================
// json 数据
// 发送消息的类型: 客户端发给服务端，服务端发给客户端
type SocketMsg struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Type    int    `json:"type"`    // 数据类型
	OpType  int    `json:"op_type"` // 操作类型
	Content string `json:"content"`
	Msg     string `json:"msg"`
}

// Message 信息转JSON (包括：发送者、接收者、内容)  ==> 历史消息
type Message struct {
	SenderID    string `json:"sender,omitempty"`
	RecipientID string `json:"recipient,omitempty"`
	Created     string `json:"created,omitempty"`
	Content     string `json:"content,omitempty"`
	Type        string `json:"type,omitempty"`
}

// 线程通信 =====================================================
// 广播类，包括广播内容和源用户
type Broadcast struct {
	From    string
	To      string
	Message []byte
	Type    int // 消息类型
	Code    int
}

// 用户类
type Client struct {
	ID             string
	CurToID        string
	Socket         *websocket.Conn
	ReceiveMessage chan *Broadcast
	Type           int
}

// 用户管理
type ClientManager struct {
	Clients             map[string]*Client
	BroadcastToAllUsers chan *Broadcast
	OnlineUsers         chan *Client
	OfflineUsers        chan *Client
}

// 函数定义 =====================================================
var Manager = ClientManager{
	Clients:             make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	BroadcastToAllUsers: make(chan *Broadcast),
	OnlineUsers:         make(chan *Client),
	OfflineUsers:        make(chan *Client),
}

func Handler(context *gin.Context) {
	uid := context.Query("uid")
	to_uid := context.Query("to_uid")

	// 验证 id 是否合法
	var count int64 = 0
	model.DB.Model(&model.User{}).Where("user_name=?", uid).First(&model.User{}).Count(&count)
	if count == 0 {
		context.JSON(http.StatusBadRequest, serializer.Response{
			Status: status.USER_NOT_EXIST,
			Msg:    status.TransCode(status.USER_NOT_EXIST),
		}) // 返回结果
		return
	}
	model.DB.Model(&model.User{}).Where("user_name=?", to_uid).First(&model.User{}).Count(&count)
	if count == 0 {
		context.JSON(http.StatusBadRequest, serializer.Response{
			Status: status.USER_NOT_EXIST,
			Msg:    status.TransCode(status.USER_NOT_EXIST),
		}) // 返回结果
		return
	}

	// 设置websocket
	// CheckOrigin防止跨站点的请求伪造
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 升级get请求为webSocket协议
	conn, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		logging.Info("ws 协议升级失败")
		logging.Info(err)
		context.JSON(http.StatusBadRequest, serializer.Response{
			Status: status.WS_CREATE_FAIL,
			Msg:    status.TransCode(status.WS_CREATE_FAIL),
		}) // 返回结果
		return
	}

	// 创建用户实例
	client := &Client{
		ID:             uid,
		CurToID:        to_uid, // 当前聊天的对象
		Socket:         conn,
		ReceiveMessage: make(chan *Broadcast),
	}
	// 用户提交到 用户管理
	Manager.OnlineUsers <- client

	go client.Read()
	go client.Write()

	context.JSON(http.StatusOK, serializer.Response{
		Status: status.WS_CREATE_SUCCESS,
		Msg:    status.TransCode(status.WS_CREATE_SUCCESS),
	}) // 返回结果
}

func (client *Client) Offline() {
	Manager.OfflineUsers <- client // 离线
}

// 连接的客户端 给服务器(也就是当前go程序)的消息读取
func (client *Client) Read() {
	defer func() {
		logging.Info("用户读关闭: " + client.ID)
		client.Offline()
	}()

	for {
		// 监听消息
		client.Socket.PongHandler()
		receive_msg := new(SocketMsg)
		err2 := client.Socket.ReadJSON(&receive_msg)
		// msg_type, msg_content, err2 := client.Socket.ReadMessage()
		// receive_msg.Type = msg_type
		// receive_msg.Content = string(msg_content)

		if err2 != nil {
			logging.Info("数据错误: ")
			logging.Info(err2)
			client.SystemMessage("数据错误", status.WS_PARSE_FAIL)
			break
		}
		// {
		//    "from":"1",
		//    "to":"2",
		//    "type": 1,
		//    "op_type":5,
		//    "content":"6655b69eb80e7e538122fc6b"
		// }

		if receive_msg.OpType == status.WS_TYPE_OP_TEXT_SEND { // 用户的操作类型，1 表示想要发送文本信息
			// 骚扰拦截
			send_link := fmt.Sprintf("%sto%s", receive_msg.From, receive_msg.To)
			send_count, _ := cache.RedisClient.Get(send_link).Result()
			target_reply_count, _ := cache.RedisClient.Get(fmt.Sprintf("%sto%s", receive_msg.To, receive_msg.From)).Result()
			logging.Info(fmt.Sprintf("send count: %s, target reply count: %s", send_count, target_reply_count))
			if target_reply_count == "" && send_count >= "5" {
				client.SystemMessage("骚扰拦截", status.WS_SEND_FORBID)
				continue
			}

			// 缓存发送记录
			if send_count == "" {
				cache.RedisClient.Set(send_link, 1, 3*time.Hour)
			} else {
				send_count_int, _ := strconv.Atoi(send_count)
				cache.RedisClient.Set(send_link, send_count_int+1, 3*A_Day) // 保留 3天
			}
			Manager.BroadcastToAllUsers <- &Broadcast{
				From:    receive_msg.From,
				To:      receive_msg.To,
				Message: []byte(receive_msg.Content),
				Type:    status.WS_TYPE_MESSAGE_TEXT, // 文本类型消息
				Code:    status.WS_RECEIVE_SUCCESS,
			}
			logging.Info("发送文本消息: ")
			logging.Info(fmt.Sprintf("<<< From %s To %s: %s", receive_msg.From, receive_msg.To, receive_msg.Content))
		} else if receive_msg.OpType == status.WS_TYPE_OP_TEXT_GET_HIS { // 获取文本历史消息
			from_uid := receive_msg.From
			to_uid := receive_msg.To

			mContent := strings.Split(receive_msg.Content, "-")
			var skip_message_count int
			var limit_message_count int
			var err error
			if len(mContent) != 2 {
				skip_message_count = 0
				limit_message_count = 999
			} else {
				skip_message_count, err = strconv.Atoi(mContent[0])
				if err != nil {
					skip_message_count = 0
				}
				limit_message_count, err = strconv.Atoi(mContent[1])
				if err != nil {
					limit_message_count = 999
				}
			}

			hisMessage, err := mongodb.FindHis(utils.Format_Uid(from_uid, to_uid), utils.Format_Uid(to_uid, from_uid), skip_message_count, limit_message_count)
			for _, hm := range hisMessage {
				client.HisTextMessage(hm.Direction, fmt.Sprintf("%s\n%v", hm.ID, hm.Content), status.WS_HIS_SUCCESS)
			}
		} else if receive_msg.OpType == status.WS_TYPE_OP_TEXT_NOT_READ_HIS { // 拉取所有未读的数据
			from_uid := receive_msg.From
			to_uid := receive_msg.To

			hisMessage, _ := mongodb.FindHisUnread(utils.Format_Uid(from_uid, to_uid), utils.Format_Uid(to_uid, from_uid))
			for _, hm := range hisMessage {
				client.HisTextMessage(hm.Direction, fmt.Sprintf("%s\n%v", hm.ID, hm.Content), status.WS_HIS_SUCCESS)
			}
		} else if receive_msg.OpType == status.WS_TYPE_OP_SET_MESSAGE_ID_READ { // 将某一个消息 ID 修改为已读
			from_uid := receive_msg.From
			to_uid := receive_msg.To
			message_id := receive_msg.Content

			err := mongodb.Message_Read(utils.Format_Uid(from_uid, to_uid), utils.Format_Uid(to_uid, from_uid), message_id)
			if err != nil {
				client.SystemMessage("修改失败"+err.Error(), status.WS_SET_READ_FAIL)
			} else {
				client.SystemMessage("修改成功", status.WS_SET_READ_SUCCESS)
			}
		} else if receive_msg.OpType == status.WS_TYPE_OP_DEL_MESSAGE_ID { // 删除某一个消息 ID
			from_uid := receive_msg.From
			to_uid := receive_msg.To
			message_id := receive_msg.Content

			err := mongodb.Message_Del(utils.Format_Uid(from_uid, to_uid), utils.Format_Uid(to_uid, from_uid), message_id)
			if err != nil {
				client.SystemMessage("删除失败"+err.Error(), status.WS_DEL_MESSAGE_FAIL)
			} else {
				client.SystemMessage("删除成功", status.WS_DEL_MESSAGE_SUCCESS)
			}
		} else if receive_msg.OpType == status.WS_TYPE_OP_TEXT_ALL_CHAT { // 群体文本聊天
			// 缓存发送记录
			cache_message_his := utils.Format_Uid(receive_msg.From, "ALL")
			cache.RedisClient.Incr(cache_message_his)
			_, _ = cache.RedisClient.Expire(cache_message_his, 3*A_Day).Result() // 缓存 3天

			Manager.BroadcastToAllUsers <- &Broadcast{
				From:    receive_msg.From,
				To:      "ALL",
				Message: []byte(receive_msg.Content),
				Type:    status.WS_TYPE_MESSAGE_TEXT, // 文本类型消息
				Code:    status.WS_RECEIVE_SUCCESS,
			}
			logging.Info("群体消息")
		}

	}
}

// 服务端 发送信息给 连接的客户端
func (client *Client) Write() {
	defer func() {
		logging.Info("用户写关闭: " + client.ID)
		client.Offline()
	}()

	for {
		select {
		case message, ok := <-client.ReceiveMessage:
			if !ok { // 监听到 关闭管道
				_ = client.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// logging.Info(fmt.Sprintf("用户%s: 接收到消息", client.ID))
			if message.Type == status.WS_TYPE_MESSAGE_TEXT { // 接收到 文本类型
				// logging.Info("消息: " + string(message.Message))
				socket_msg := &SocketMsg{
					From:    message.From,
					To:      message.To,
					Content: string(message.Message),
					Type:    message.Type,
					OpType:  message.Code, // 服务端返回给客户端 操作
					Msg:     status.TransCode(message.Code),
				}
				json_msg, err := json.Marshal(socket_msg)
				if err != nil {
					logging.Info("json序列化失败")
					socket_msg.Content = ""
					socket_msg.OpType = message.Code - 1
					socket_msg.Type = status.WS_TYPE_MESSAGE_TEXT
					socket_msg.Msg = status.TransCode(socket_msg.OpType)
					json_msg, _ = json.Marshal(socket_msg)
				}
				err = client.Socket.WriteMessage(websocket.TextMessage, json_msg) // 发送信息给 连接的客户端
				if err != nil {
					logging.Info("发送失败")
					logging.Info(err)
					continue
				}
			}
			if message.Type == status.WS_TYPE_MESSAGE_HIS_TEXT { // 接收到 历史信息且是文本类型
				// logging.Info("消息: " + string(message.Message))
				socket_msg := &SocketMsg{
					From:    message.From,
					To:      message.To,
					Content: string(message.Message),
					Type:    message.Type,
					OpType:  message.Code, // 服务端返回给客户端 操作
					Msg:     status.TransCode(message.Code),
				}
				json_msg, err := json.Marshal(socket_msg)
				if err != nil {
					logging.Info("json序列化失败")
					socket_msg.Content = ""
					socket_msg.OpType = message.Code - 1
					socket_msg.Type = status.WS_TYPE_MESSAGE_TEXT
					socket_msg.Msg = status.TransCode(socket_msg.OpType)
					json_msg, _ = json.Marshal(socket_msg)
				}
				err = client.Socket.WriteMessage(websocket.TextMessage, json_msg) // 发送信息给 连接的客户端
				if err != nil {
					logging.Info("发送失败")
					logging.Info(err)
					continue
				}
			}

		}
	}
}

func (client *Client) SystemMessage(message string, code int) {
	broadcast := &Broadcast{From: "System", To: client.ID, Message: []byte(message), Type: status.WS_TYPE_MESSAGE_TEXT, Code: code}
	client.ReceiveMessage <- broadcast // 服务端 发给 客户端
}

func (client *Client) HisTextMessage(direction, message string, code int) {
	ret := strings.Split(direction[1:len(direction)-1], "SEP")
	broadcast := &Broadcast{From: ret[0], To: ret[1], Message: []byte(message), Type: status.WS_TYPE_MESSAGE_HIS_TEXT, Code: code}
	client.ReceiveMessage <- broadcast // 服务端 发给 客户端
}
