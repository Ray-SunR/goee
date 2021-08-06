package cache

import "sync"

// The primary benefit of having functional interface (https://karthikkaranth.me/blog/functions-implementing-interfaces-in-go/) is to
// allow developers to use anymonous function or struct (they don't necessarily need to create a struct then implement
// Getter interface. They can use a function or anyounous function or a struct that implements the Getter interface. That
// gives them more flexibility.
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name   string
	cache  cache
	getter Getter
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group, 0)
)

func NewGroup(name string, maxSizeInBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter cannot be nil")
	}

	mu.Lock()
	defer mu.Unlock()
	groups[name] = &Group{name: name, cache: cache{maxSizeInBytes: maxSizeInBytes}, getter: getter}
	return groups[name]
}

func GetGroup(name string) *Group {
	if groups == nil {
		return nil
	}

	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (group *Group) Get(key string) (ByteView, error) {
	if v, ok := group.cache.get(key); ok {
		return v, nil
	}

	byteView, err := group.load(key)
	if err != nil {
		return ByteView{}, err
	}

	// This byteview is a already a cloned copy.
	err = group.cache.add(key, byteView)
	return byteView, err
}

func (group *Group) load(key string) (ByteView, error) {
	return group.loadLocally(key)
}

func (group *Group) loadLocally(key string) (ByteView, error) {
	bytes, err := group.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	byteView := ByteView{bytes: cloneBytes(bytes)}
	group.populateCache(key, byteView)
	// Make a deep copy of the slice to avoid data being modified outside.
	// Note that we have to clone here because bytes is a slice and slice passes by reference.
	return byteView, nil
}

func (group *Group) populateCache(key string, byteView ByteView) {
	err := group.cache.add(key, byteView)
	if err != nil {
		return
	}
}