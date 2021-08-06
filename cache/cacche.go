package cache

import (
	"example/cache/lru"
	"sync"
)

type cache struct {
	mutex          sync.Mutex
	lru            *lru.Cache
	maxSizeInBytes int64
}

func (c *cache) add(key string, value ByteView) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.maxSizeInBytes, nil)
	}
	return c.lru.Add(key, value)
}

func (c *cache) get(key string) (ByteView, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var result ByteView
	if c.lru == nil {
		return result, false
	}
	if value, ok := c.lru.Get(key); ok {
		return value.(ByteView), ok
	} else {
		return result, false
	}
}
