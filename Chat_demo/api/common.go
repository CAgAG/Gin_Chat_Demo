package api

import (
	"Chat_demo/pkg/status"
	"Chat_demo/serializer"
	"fmt"
)

// 返回错误信息 ErrorResponse
func ErrorResponse(err error, code int) serializer.Response {
	return serializer.Response{
		Status: code,
		Msg:    status.TransCode(code),
		Error:  fmt.Sprint(err),
	}
}
