package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 直接使用 Handlers 的 Login 方法作为路由处理器
	r.POST("/login", Handlers{}.Login)
}
