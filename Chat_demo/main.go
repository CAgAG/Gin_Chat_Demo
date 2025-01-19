package main

import (
	"Chat_demo/cache"
	"Chat_demo/conf"
	"Chat_demo/router"
	"Chat_demo/service"
	"fmt"
	logging "github.com/sirupsen/logrus"
)

func main() {
	var err error

	// mysql, mongodb
	conf.Init()
	// redis
	cache.Init()
	go service.Manager.Start()

	// 开启服务、路由
	server := router.NewRouter()
	if server == nil {
		logging.Info("路由启动失败")
		return
	}

	fmt.Println("http://127.0.0.1:" + conf.HttpPort + "/ping")
	err = server.Run(":" + conf.HttpPort)

	if err != nil {
		logging.Info(err)
		panic(err)
	}
}
