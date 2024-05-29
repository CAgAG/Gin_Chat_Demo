package status

import "net/http"

const (
	NET_SUCCESS     = http.StatusOK
	NET_BAD_REQUIST = http.StatusBadRequest
	NET_FAIL        = http.StatusNotFound

	USER_REGISTER_FAIL        = 41000
	USER_REGISTER_SUCCESS     = 41001
	USER_LOGIN_FAIL           = 41002
	USER_LOGIN_SUCCESS        = 41003
	USER_SET_PASSWORD_FAIL    = 41004
	USER_SET_PASSWORD_SUCCESS = 41005
	USER_NOT_EXIST            = 41006
	USER_EXIST                = 41007
	USER_SET_SESSION_FAIL     = 41008
	USER_SET_SESSION_SUCCESS  = 41009
	USER_INCORRECT_PASSWORD   = 41010
	USER_CORRECT_PASSWORD     = 41011
	USER_NOT_LOGIN            = 41012
	USER_HAVE_LOGINED         = 41013
	USER_LOGOUT_FAIL          = 41014
	USER_LOGOUT_SUCCESS       = 41015

	WS_PARSE_FAIL          = 40000
	WS_PARSE_SUCCESS       = 40001
	WS_LINK_FAIL           = 40002
	WS_LINK_SUCCESS        = 40003
	WS_HIS_FAIL            = 40004
	WS_HIS_SUCCESS         = 40005
	WS_SET_READ_FAIL       = 40006
	WS_SET_READ_SUCCESS    = 40007
	WS_SEND_FAIL           = 40008
	WS_SEND_SUCCESS        = 40009
	WS_RECEIVE_FAIL        = 40010
	WS_RECEIVE_SUCCESS     = 40011
	WS_DEL_MESSAGE_FAIL    = 40012
	WS_DEL_MESSAGE_SUCCESS = 40013
	WS_CREATE_FAIL         = 40014
	WS_CREATE_SUCCESS      = 40015
	WS_LINK_OUT_FAIL       = 40016
	WS_LINK_OUT_SUCCESS    = 40017
	WS_USER_ONLINE         = 40018
	WS_USER_OFFLINE        = 40019
	WS_SEND_FORBID         = 40020
	WS_SEND_RECOVER        = 40021

	WS_TYPE_MESSAGE_TEXT     = 1
	WS_TYPE_MESSAGE_HIS_TEXT = 2

	WS_TYPE_OP_TEXT_SEND           = 1
	WS_TYPE_OP_TEXT_GET_HIS        = 2
	WS_TYPE_OP_TEXT_NOT_READ_HIS   = 3
	WS_TYPE_OP_SET_MESSAGE_ID_READ = 4
	WS_TYPE_OP_DEL_MESSAGE_ID      = 5
	WS_TYPE_OP_TEXT_ALL_CHAT       = 6

	CODE_NOT_FOUND = 44444
)

var MsgFlags = map[int]string{
	NET_SUCCESS:     "请求成功",
	NET_FAIL:        "资源不存在",
	NET_BAD_REQUIST: "错误连接",

	USER_REGISTER_FAIL:        "用户 注册失败",
	USER_REGISTER_SUCCESS:     "用户 注册成功",
	USER_LOGIN_FAIL:           "用户 登录失败",
	USER_LOGIN_SUCCESS:        "用户 登录成功",
	USER_SET_PASSWORD_FAIL:    "用户 设置密码失败",
	USER_SET_PASSWORD_SUCCESS: "用户 设置密码成功",
	USER_NOT_EXIST:            "用户 不存在",
	USER_EXIST:                "用户 已存在",
	USER_SET_SESSION_FAIL:     "用户 设置session失败",
	USER_SET_SESSION_SUCCESS:  "用户 设置session成功",
	USER_INCORRECT_PASSWORD:   "用户 密码错误",
	USER_CORRECT_PASSWORD:     "用户 密码正确",
	USER_NOT_LOGIN:            "用户 未登录",
	USER_HAVE_LOGINED:         "用户 已登录",
	USER_LOGOUT_FAIL:          "用户 登出失败",
	USER_LOGOUT_SUCCESS:       "用户 登出成功",

	WS_CREATE_FAIL:         "websocket 创建失败",
	WS_CREATE_SUCCESS:      "websocket 创建成功",
	WS_PARSE_FAIL:          "websocket 数据解析失败",
	WS_PARSE_SUCCESS:       "websocket 数据解析成功",
	WS_LINK_FAIL:           "websocket 连接失败",
	WS_LINK_SUCCESS:        "websocket 连接成功",
	WS_HIS_FAIL:            "websocket 请求历史消息失败",
	WS_HIS_SUCCESS:         "websocket 请求历史消息成功",
	WS_SET_READ_FAIL:       "websocket 设置消息为已读失败",
	WS_SET_READ_SUCCESS:    "websocket 设置消息为已读成功",
	WS_SEND_FAIL:           "websocket 发送消息失败",
	WS_SEND_SUCCESS:        "websocket 发送消息成功",
	WS_RECEIVE_FAIL:        "websocket 接收消息失败",
	WS_RECEIVE_SUCCESS:     "websocket 接收消息成功",
	WS_DEL_MESSAGE_FAIL:    "websocket 删除消息失败",
	WS_DEL_MESSAGE_SUCCESS: "websocket 删除消息成功",
	WS_LINK_OUT_FAIL:       "websocket 断开连接失败",
	WS_LINK_OUT_SUCCESS:    "websocket 断开连接成功",
	WS_USER_ONLINE:         "websocket 用户在线",
	WS_USER_OFFLINE:        "websocket 用户离线",
	WS_SEND_FORBID:         "websocket 用户聊天封禁",
	WS_SEND_RECOVER:        "websocket 用户聊天解封",

	CODE_NOT_FOUND: "未知状态码",
}

// 获取状态码对应信息
func TransCode(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[CODE_NOT_FOUND]
}
