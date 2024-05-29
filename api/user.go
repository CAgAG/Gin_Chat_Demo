package api

import (
	"Chat_demo/pkg/status"
	"Chat_demo/service"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
	"net/http"
)

func UserRegister(context *gin.Context) {
	var userBaseServer service.UserBaseService

	// 处理 HTTP请求参数绑定到结构体
	if err := context.ShouldBind(&userBaseServer); err == nil {
		res := userBaseServer.Register()
		context.JSON(http.StatusOK, res) // 返回结果
	} else {
		context.JSON(http.StatusBadRequest, ErrorResponse(err, status.USER_REGISTER_FAIL))
		logging.Info(err)
	}
}

func UserLogin(context *gin.Context) {
	var userBaseServer service.UserBaseService

	// 处理 HTTP请求参数绑定到结构体
	if err := context.ShouldBind(&userBaseServer); err == nil {
		res := userBaseServer.Login(context)
		context.JSON(http.StatusOK, res) // 返回结果
	} else {
		context.JSON(http.StatusBadRequest, ErrorResponse(err, status.USER_LOGIN_FAIL))
		logging.Info(err)
	}
}

func SetPassWord(context *gin.Context) {
	var userBaseServer service.UserSetPasswordService

	// 处理 HTTP请求参数绑定到结构体
	if err := context.ShouldBind(&userBaseServer); err == nil {
		res := userBaseServer.SetPassWord(context)
		context.JSON(http.StatusOK, res) // 返回结果
	} else {
		context.JSON(http.StatusBadRequest, ErrorResponse(err, status.USER_SET_PASSWORD_FAIL))
		logging.Info(err)
	}
}

func Logout(context *gin.Context) {
	var userBaseServer service.UserBaseService

	// 处理 HTTP请求参数绑定到结构体
	if err := context.ShouldBind(&userBaseServer); err == nil {
		res := userBaseServer.Logout(context)
		context.JSON(http.StatusOK, res) // 返回结果
	} else {
		context.JSON(http.StatusBadRequest, ErrorResponse(err, status.USER_LOGOUT_FAIL))
		logging.Info(err)
	}
}
