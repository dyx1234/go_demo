## agollo demo

1. 学习go的开发最佳实践
2. 在阅读完源码后通过使用agollo加深对源码的理解
3. 发现可以自定义组件，于是自定义了使用configmap的缓存组件
```
		agollo.SetCache(&cache.ConfigMapCacheFactory{})
```
4. 自定义缓存组件：使用client-go实现了cache接口定义的方法，增强对client-go的理解
