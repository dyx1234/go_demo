package cache

// CacheInterface 自定义缓存组件接口
type CacheInterface interface {
	Set(key string, value interface{}, expireSeconds int) (err error)

	EntryCount() (entryCount int64)

	Get(key string) (value interface{}, err error)

	Del(key string) (affected bool)

	Range(f func(key, value interface{}) bool)

	Clear()
}

// CacheFactory 缓存组件工厂接口
type CacheFactory interface {
	//Create 创建缓存组件
	Create() CacheInterface
}
