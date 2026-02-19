package utils

import "testing"

func TestHashRing_Add_N_Get(t *testing.T) {
	ring := New(5)
	ring.Add("redis-A")
	ring.Add("redis-B")
	ring.Add("redis-C")
	ring.Add("redis-D")

	keys := []string{"user0", "user1", "user2", "user3", "user4", "user5", "user6", "user7", "user44", "user99"}
	for _, key := range keys {
		t.Log(key, "->", ring.Get(key))
	}
}
