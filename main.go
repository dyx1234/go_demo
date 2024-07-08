package main

import (
	"github.com/gin-gonic/gin"
	"go_demo/config"
	"go_demo/manager"
	"go_demo/server"
)

func main() {

	manager.NewKubernetesManager()
	config.InitApolloClient()

	// 创建 Gin 实例
	r := gin.Default()

	// 设置路由
	server.SetupRoutes(r)

	// 运行 Web 服务
	r.Run(":8082")

}
