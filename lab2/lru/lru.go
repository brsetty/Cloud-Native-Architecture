package lru

import (
	"errors"
	"fmt"
)

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Put(interface{}, interface{}) error
}

type lruCache struct {
	size      int
	remaining int
	cache     map[string]string
	queue     []string
}

func NewCache(size int) Cacher {
	return &lruCache{size: size, remaining: size, cache: make(map[string]string), queue: make([]string, 0)}
}

func (lru *lruCache) Get(key interface{}) (interface{}, error) {

	kstr := key.(string)

	//check if it is in map
	val, found := lru.cache[kstr]
	if !found {
		return nil, errors.New(" key not found in cache")
	}
	// we delete the key from queue
	lru.qDel(kstr)
	// we insert the key again at the head of the queue
	lru.queue = append(lru.queue, kstr)
	fmt.Printf("get called for %s", kstr)
	fmt.Println(lru.queue)
	return val, nil

}

func (lru *lruCache) Put(key, val interface{}) error {
	// Your code here....
	kstr := key.(string)
	valstr := val.(string)

	//check for remaining size and pop head
	if lru.remaining == 0 {
		delete(lru.cache, lru.queue[0])
		lru.qDel(lru.queue[0])
		lru.remaining++
	}
	//insert the new key to the tail of the queue
	lru.cache[kstr] = valstr // map[key] = val;
	lru.queue = append(lru.queue, kstr)
	lru.remaining--
	fmt.Printf(" put called for %s %s\n", kstr, valstr)
	fmt.Println(lru.queue)
	fmt.Println(lru.cache[valstr])
	return nil

}

// Delete element from queue
func (lru *lruCache) qDel(ele string) {
	for i := 0; i < len(lru.queue); i++ {
		if lru.queue[i] == ele {
			oldlen := len(lru.queue)
			copy(lru.queue[i:], lru.queue[i+1:])
			lru.queue = lru.queue[:oldlen-1]
			break
		}
	}
}
