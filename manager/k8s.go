package manager

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

type KubernetesManager struct {
	client *kubernetes.Clientset
}

type ConfigMapInfo struct {
	Name      string
	Namespace string
	Data      map[string]string
}

func NewKubernetesManager() (*KubernetesManager, error) {
	// TODO 连接k8s集群
	// 使用InClusterConfig()在集群内部获取配置
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Can't get in-cluster config: %v\n", err)
	}

	// 使用配置创建Kubernetes客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Can't create Kubernetes client: %v\n", err)
	}
	return &KubernetesManager{client: client}, nil
}

// 根据key获取ConfigMap中的value
func (km *KubernetesManager) GetConfigMapValue(namespace, name, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configMap, err := km.client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metaV1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return "", fmt.Errorf("ConfigMap %s not found in namespace %s", name, namespace)
		}
	}

	return configMap.Data[key], nil
}

// set，更新ConfigMap
func (km *KubernetesManager) SetConfigMap(namespace, name string, data map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configMap, err := km.client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metaV1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("ConfigMap %s not found in namespace %s", name, namespace)
		}
	}

	// 更新值
	configMap.Data = data
	_, err = km.client.CoreV1().ConfigMaps(namespace).Update(ctx, configMap, metaV1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (km *KubernetesManager) CreateConfigMap(namespace string, configMapInfo *ConfigMapInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configMap := &corev1.ConfigMap{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      configMapInfo.Name,
			Namespace: namespace,
		},
		Data: configMapInfo.Data,
	}
	_, err := km.client.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, metaV1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (km *KubernetesManager) DeleteConfigMap(namespace, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := km.client.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metaV1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("ConfigMap %s not found in namespace %s", name, namespace)
		}
		return fmt.Errorf("failed to delete ConfigMap %s in namespace %s: %v", name, namespace, err)
	}
	return nil
}
