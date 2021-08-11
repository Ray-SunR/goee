package consistenthash

import (
	"hash/crc32"
	"log"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

type Map struct {
	hashFunc Hash // hash function
	replica int // number of virtual nodes per physical node
	virtualKeys []int // virtual node keys
	keysMap map[int]string // mapping between virtual nodes to physical node
}

func New(replica int, hashFunc Hash) *Map {
	m := &Map{replica: replica, hashFunc: hashFunc, keysMap: make(map[int]string)}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) AddPhysicalNodes(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replica; i++ {
			virtualKey := key + strconv.Itoa(i)
			hashValue := int(m.hashFunc([]byte(virtualKey)))
			log.Printf("key: %s, virtual key: %s, hashValue: %d", key, virtualKey, hashValue)
			m.virtualKeys = append(m.virtualKeys, hashValue)
			m.keysMap[hashValue] = key
		}
	}
	sort.Ints(m.virtualKeys)
}

func (m *Map) GetPhysicalNode(key string) string {
	hashValue := int(m.hashFunc([]byte(key)))
	// Use this hash value to find the first key which is equal or greater than itself
	index := sort.Search(len(m.virtualKeys), func(index int) bool {
		return m.virtualKeys[index] >= hashValue
	})
	return m.keysMap[m.virtualKeys[index % len(m.virtualKeys)]]
}