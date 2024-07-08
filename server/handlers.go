package server

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/gin-gonic/gin"
	"go_demo/config"
	"go_demo/info"
	"net/http"
)

// UserCredentials 用于表示用户登录凭据的结构体
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Handlers 包含登录处理逻辑
type Handlers struct {
	client *agollo.Client
}

// Login 处理登录请求
func (h Handlers) Login(c *gin.Context) {
	var user UserCredentials
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var client *agollo.Client

	h.client = config.GetApolloClient()
	if h.client == nil {
		config.InitApolloClient()
	} else {
		client = h.client
	}

	cache := (*client).GetConfigCache(info.SecretNameSpace)

	expectedUsername, _ := cache.Get("username")
	expectedPassword, _ := cache.Get("password")

	// 这里添加验证逻辑，例如检查用户名和密码
	if user.Username != expectedUsername || user.Password != expectedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// 如果验证成功，可以设置token或其他响应数据
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
