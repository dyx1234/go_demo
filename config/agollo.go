package config

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"go_demo/component/cache"
	"go_demo/component/logger"
	"go_demo/info"
	"sync"
)

var c = &config.AppConfig{
	AppID:          info.SecretAppID,
	Cluster:        "LOCAL",
	IP:             "http://127.0.0.1:8080",
	NamespaceName:  info.SecretNameSpace,
	IsBackupConfig: false,
}

// 单例模式, 只初始化一次
var (
	clientInstance *agollo.Client
	once           sync.Once
)

func InitApolloClient() {
	once.Do(func() {
		// 初始化 Apollo 客户端
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
	})
}

// 提供实例的暴露点
func GetApolloClient() *agollo.Client {
	if clientInstance == nil {
		panic("Apollo client is not initialized")
	}
	return clientInstance
}
