package lru

import "testing"

type String string

func (d String) SizeInBytes() int64 {
	return (int64)(len(d))
}

func (d String) String() string {
	return string(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(1024), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)

	onEvictedCalled := false
	lru := New(int64(cap), func(key string, value Value) {
		onEvictedCalled = true
		t.Logf("key: %s, value: %s get removed!", key, value)
		if key != k1 || value.(String).String() != v1 {
			t.Fatalf("Unexpected entry gets removed, expected key: %s, value: %s to be removed but got key: %s, value: %s", k1, v1, key, value)
		}
	})
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if !onEvictedCalled {
		t.Fatalf("OnEvicted callback not called!")
	}
	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}
