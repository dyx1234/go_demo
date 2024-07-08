package cache

import (
	"context"
	"fmt"
	"github.com/apolloconfig/agollo/v4/agcache"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type ConfigMapCache struct {
	clientset *kubernetes.Clientset
	namespace string
}

func (c *ConfigMapCache) Set(key string, value interface{}, expireSeconds int) error {
	// 将 expireSeconds 参数忽略，因为 Kubernetes ConfigMap 不支持过期时间
	valueStr, ok := value.(string)
	if !ok {
		return errors.NewBadRequest("value must be a string")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configMapName := "apollo-cache" // 假定使用固定的 ConfigMap 名称
	data := map[string]string{
		key: valueStr,
	}

	// 尝试获取 ConfigMap，如果不存在则创建
	cm, err := c.clientset.CoreV1().ConfigMaps(c.namespace).Get(ctx, configMapName, metaV1.GetOptions{})
	if errors.IsNotFound(err) {
		cm = &coreV1.ConfigMap{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      configMapName,
				Namespace: c.namespace,
			},
			Data: data,
		}
		_, err = c.clientset.CoreV1().ConfigMaps(c.namespace).Create(ctx, cm, metaV1.CreateOptions{})
	} else if err != nil {
		return err
	} else {
		// ConfigMap 存在，更新数据
		cm.Data[key] = valueStr
		_, err = c.clientset.CoreV1().ConfigMaps(c.namespace).Update(ctx, cm, metaV1.UpdateOptions{})
	}

	return err
}

// 实现接口的 EntryCount 方法
func (c *ConfigMapCache) EntryCount() (entryCount int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cm, err := c.clientset.CoreV1().ConfigMaps(c.namespace).Get(ctx, "apollo-cache", metaV1.GetOptions{})
	if err != nil {
		return 0
	}
	return int64(len(cm.Data))
}

// 实现接口的 Get 方法
func (c *ConfigMapCache) Get(key string) (value interface{}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cm, err := c.clientset.CoreV1().ConfigMaps(c.namespace).Get(ctx, "apollo-cache", metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	valueStr, ok := cm.Data[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", key)
	}
	return valueStr, nil
}

// 实现接口的 Del 方法
func (c *ConfigMapCache) Del(key string) (affected bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cm, err := c.clientset.CoreV1().ConfigMaps(c.namespace).Get(ctx, "apollo-cache", metaV1.GetOptions{})
	if err != nil {
		return false
	}
	delete(cm.Data, key)
	_, err = c.clientset.CoreV1().ConfigMaps(c.namespace).Update(ctx, cm, metaV1.UpdateOptions{})
	return err == nil
}

// 实现接口的 Range 方法
func (c *ConfigMapCache) Range(f func(key, value interface{}) bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cm, err := c.clientset.CoreV1().ConfigMaps(c.namespace).Get(ctx, "apollo-cache", metaV1.GetOptions{})
	if err != nil {
		return
	}
	for key, value := range cm.Data {
		if !f(key, value) {
			break
		}
	}
}

// 实现接口的 Clear 方法
func (c *ConfigMapCache) Clear() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.clientset.CoreV1().ConfigMaps(c.namespace).Delete(ctx, "apollo-cache", metaV1.DeleteOptions{})
	if err != nil {
		fmt.Println("Failed to clear cache:", err)
	}
}

// ConfigMapCacheFactory 用于创建 ConfigMapCache 实例的工厂
type ConfigMapCacheFactory struct {
	// 可以添加一些配置参数，比如客户端配置、命名空间等
	clientset *kubernetes.Clientset
	namespace string
}

// Create 创建并返回一个实现了 CacheInterface 的 ConfigMapCache 实例
func (f *ConfigMapCacheFactory) Create() agcache.CacheInterface {
	return &ConfigMapCache{
		clientset: f.clientset,
		namespace: f.namespace,
	}
}
