package helper

import "sync"

type Cache[K any, V any] struct {
	cache sync.Map
}

func NewCache[K any, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		cache: sync.Map{},
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	value, ok := c.cache.Load(key)
	if !ok {
		var zeroV V // zero value of type V
		return zeroV, false
	}

	return value.(V), true
}

func (c *Cache[K, V]) Exists(key K) bool {
	_, ok := c.cache.Load(key)
	return ok
}

func (c *Cache[K, V]) Add(key K, value V) bool {
	c.cache.Store(key, value)
	_, ok := c.cache.Load(key)
	return ok
}

func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, 0)
	c.cache.Range(func(key, _ interface{}) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

func (c *Cache[K, V]) Delete(key K) bool {
	c.cache.Delete(key)
	_, ok := c.cache.Load(key)
	return !ok
}
