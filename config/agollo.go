package config

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"go_demo/component/cache"
	"go_demo/component/logger"
	"go_demo/info"
)

var c = &config.AppConfig{
	AppID:          info.SecretAppID,
	Cluster:        "LOCAL",
	IP:             "http://127.0.0.1:8080",
	NamespaceName:  info.SecretNameSpace,
	IsBackupConfig: false,
}

var (
	clientInstance *agollo.Client
)

func init() {
	// 自定义组件
	agollo.SetLogger(&logger.DefaultLogger{})
	agollo.SetCache(&cache.ConfigMapCacheFactory{})

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		fmt.Println("err:", err)
		panic(err)
	}
	clientInstance = &client
}

// GetApolloClient 提供实例的暴露点
func GetApolloClient() *agollo.Client {
	if clientInstance == nil {
		panic("Apollo client is not initialized")
	}
	return clientInstance
}
