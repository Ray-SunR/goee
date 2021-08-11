package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.AddPhysicalNodes("6", "4", "2")

	testCases := map[string]string{
		"20":  "2",
		"21": "2",
		"22": "2",
		"40": "4",
		"41": "4",
		"42": "4",
		"60": "6",
		"61": "6",
		"62": "6",
	}

	for k, v := range testCases {
		if hash.GetPhysicalNode(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	hash.AddPhysicalNodes("8")

	// 27 should now map to 8.
	testCases["80"] = "8"
	testCases["81"] = "8"
	testCases["82"] = "8"

	for k, v := range testCases {
		if hash.GetPhysicalNode(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}