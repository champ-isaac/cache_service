package utils

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashRing struct {
	replica int
	keys    []int
	hashMap map[int]string
}

func New(replica int) *HashRing {
	return &HashRing{
		replica: replica,
		hashMap: make(map[int]string),
	}
}

func hashKey(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}

func (h *HashRing) Add(node string) {
	for i := 0; i < h.replica; i++ {
		virtualKey := hashKey(node + "#" + strconv.Itoa(i))
		h.keys = append(h.keys, virtualKey)
		h.hashMap[virtualKey] = node
	}
	sort.Ints(h.keys)
}

func (h *HashRing) Get(key string) string {
	if len(h.keys) == 0 {
		return ""
	}
	hash := hashKey(key)
	idx := sort.Search(len(h.keys), func(i int) bool {
		//check which range slot hash fit in
		return h.keys[i] >= hash
	})
	if idx == len(h.keys) {
		idx = 0
	}
	return h.hashMap[h.keys[idx]]
}
