package router

import (
	"Chat_demo/api"
	"Chat_demo/cache"
	"Chat_demo/service"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
	"net/http"
)

func NewRouter() *gin.Engine {
	server := gin.Default()

	// session
	store, err := cache.NewSession()
	if err != nil {
		logging.Info("session 设置失败")
		logging.Info(err)
		return nil
	}
	server.Use(sessions.Sessions("mysession", store))
	// 恢复 和 日志
	server.Use(gin.Recovery(), gin.Logger())

	rootGroup := server.Group("/")
	{
		rootGroup.GET("ping", func(context *gin.Context) {
			session := sessions.Default(context)
			user_name := session.Get("UserName")
			if us, ok := user_name.(string); ok {
				context.JSON(http.StatusOK, fmt.Sprintf("welcome: %s", us))
			} else {
				context.JSON(http.StatusOK, "ok")
			}
		})
	}

	userGroup := server.Group("/user")
	{
		// user_name, password
		userGroup.POST("/register", api.UserRegister)
		// 记录 session
		userGroup.POST("/login", api.UserLogin)
		// 修改密码
		userGroup.POST("/set_password", api.SetPassWord)
		// 退出用户
		userGroup.POST("/logout", api.Logout)
	}

	chatGroup := server.Group("/chat")
	{
		// ws://127.0.0.1:8081/chat/test?uid=1&to_uid=2
		// 注意: 使用的是 ws开头的连接, 表明使用websocket
		// 可以使用 http://www.jsons.cn/websocket/
		chatGroup.GET("/test", service.Handler)
	}

	return server
}
