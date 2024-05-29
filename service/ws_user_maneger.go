package service

import (
	"Chat_demo/model/mongodb"
	"Chat_demo/pkg/status"
	"Chat_demo/pkg/utils"
	logging "github.com/sirupsen/logrus"
)

func (c_manager *ClientManager) Start() {
	for {
		select {
		case client, ok := <-c_manager.OnlineUsers:
			if !ok {
				client.SystemMessage("连接失败", status.WS_LINK_FAIL)
				continue
			}
			if _, ok := Manager.Clients[client.ID]; !ok {
				Manager.Clients[client.ID] = client
			}
			logging.Info("已接入一个新的连接: " + client.ID)
			client.SystemMessage("连接成功", status.WS_LINK_SUCCESS)

		case client, ok := <-c_manager.OfflineUsers:
			if !ok {
				client.SystemMessage("断开连接失败", status.WS_LINK_OUT_FAIL)
				continue
			}
			logging.Info("已断开一个连接: " + client.ID)
			if _, ok := Manager.Clients[client.ID]; ok {
				// client.Offline()
				close(client.ReceiveMessage) // 关闭 chan
				_ = client.Socket.Close()    // 关闭 socket
				delete(Manager.Clients, client.ID)
			}
			// client.SystemMessage("断开连接成功", status.WS_LINK_OUT_SUCCESS)
		case broadcast, ok := <-c_manager.BroadcastToAllUsers:
			if !ok {
				continue
			}

			// 群体在线聊天, 不存储
			if broadcast.To == "ALL" {
				for _, other_client := range Manager.Clients {
					if other_client.ID == broadcast.From {
						continue
					}
					select {
					case other_client.ReceiveMessage <- broadcast:

					}
				}
				continue
			}

			// 2个用户聊天
			to_client, online_flag := Manager.Clients[broadcast.To]
			from_client, online_flag2 := Manager.Clients[broadcast.From]

			if online_flag { // 接收方在线
				select {
				case to_client.ReceiveMessage <- broadcast:
					// 插入 mongodb, 1 => 表明接收方【已在线】接收到消息
					err := InsertMsg(broadcast, 1, int64(3*A_Day_INT))
					if err != nil {
						logging.Info("Insert mongodb err: ")
						logging.Info(err)
					}
				default: // 接收方 正好在退出中
					// 插入 mongodb, 0 => 表明接收方【未在线】接收到消息
					err := InsertMsg(broadcast, 0, int64(3*A_Day_INT))
					if err != nil {
						logging.Info("Insert mongodb err: ")
						logging.Info(err)
					}
					online_flag = false
				}
			} else {
				// 插入 mongodb, 0 => 表明接收方【未在线】接收到消息
				err := InsertMsg(broadcast, 0, int64(3*A_Day_INT))
				if err != nil {
					logging.Info("Insert mongodb err: ")
					logging.Info(err)
				}
			}
			// ============================================
			if online_flag && online_flag2 { // 发送方在线
				from_client.SystemMessage("接收方在线", status.WS_USER_ONLINE)
			}
			if !online_flag && online_flag2 {
				from_client.SystemMessage("接收方不在线", status.WS_USER_OFFLINE)
			}

		}

	}
}

// mongodb 插入一条数据
func InsertMsg(broadcast *Broadcast, read uint, expire int64) error {
	id := utils.Format_Uid(broadcast.From, broadcast.To)
	content := string(broadcast.Message)
	err := mongodb.MongoInsert(id, content, expire, read)
	return err
}
