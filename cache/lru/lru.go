package lru

import (
	"container/list"
	"fmt"
)

type Cache struct {
	maxBytes  int64
	usedBytes int64
	data      *list.List
	cache     map[string]*list.Element

	OnEvicted func(key string, value Value)
}
type entry struct {
	key   string
	value Value
}

type Value interface {
	SizeInBytes() int64
}

type OverLimitError struct {
	maxSize    int64
	objectSize int64
	keySize    int64
}

func (error *OverLimitError) Error() string {
	return fmt.Sprintf("Cannot insert value into cache. Max size: %d, object size: %d", error.maxSize, error.objectSize)
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		data:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (cache *Cache) Get(key string) (Value, bool) {
	if ele, ok := cache.cache[key]; ok {
		kv := ele.Value.(*entry)
		cache.data.MoveToFront(ele)
		return kv.value, true
	}
	return nil, false
}

func (cache *Cache) RemoveOldest() {
	ele := cache.data.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		cache.data.Remove(cache.data.Back())
		delete(cache.cache, kv.key)
		if cache.OnEvicted != nil {
			cache.OnEvicted(kv.key, kv.value)
		}
		cache.usedBytes -= (int64)(len(kv.key)) + kv.value.SizeInBytes()
	}
}

func (cache *Cache) Len() int {
	return cache.data.Len()
}

func (cache *Cache) Add(key string, value Value) error {
	if (int64)(len(key))+value.SizeInBytes() > cache.maxBytes {
		return &OverLimitError{
			objectSize: value.SizeInBytes(),
			keySize:    (int64)(len(key)),
			maxSize:    cache.maxBytes,
		}
	}
	// Remove items to make room
	for cache.usedBytes+value.SizeInBytes()+(int64)(len(key)) > cache.maxBytes {
		cache.RemoveOldest()
	}

	if ele, ok := cache.cache[key]; ok {
		cache.data.MoveToFront(ele)
		kv := ele.Value.(*entry)
		kv.value = value
		cache.usedBytes -= kv.value.SizeInBytes()
		cache.usedBytes += value.SizeInBytes()
	} else {
		newElement := cache.data.PushFront(&entry{key: key, value: value})
		cache.cache[key] = newElement
		cache.usedBytes += (int64)(len(key)) + value.SizeInBytes()
	}
	return nil
}
