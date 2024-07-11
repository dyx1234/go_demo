package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_demo/component/cache"
	"go_demo/server"
	"os"
)

var port = ":8082"

func main() {

	// 初始化 Kubernetes 管理器
	_, err := cache.NewKubernetesManager()
	if err != nil {
		fmt.Printf("Failed to initialize Kubernetes manager: %v\n", err)
		os.Exit(1)
	}

	// 创建 Gin 实例
	r := gin.Default()

	// 设置路由
	server.SetupRoutes(r)

	// 运行 Web 服务
	err = r.Run(port)
	if err != nil {
		fmt.Printf("Failed to start the server on port %s: %v\n", port, err)
		os.Exit(1)
	}

}
