package main

import (
	"container/list"
	"sync"
)

type LRU[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
}
type entry[K comparable, V any] struct {
	key   K
	value V
}
type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	ll       list.List
	items    map[K]*list.Element
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		ll:       list.List{},
		items:    make(map[K]*list.Element),
	}
}

func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	element, ok := c.items[key]
	if !ok {
		return value, false
	}

	c.ll.MoveToBack(element)

	item := element.Value.(entry[K, V])

	return item.value, true
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity <= 0 {
		return
	}

	if element, ok := c.items[key]; ok {
		item := element.Value.(entry[K, V])
		item.value = value
		element.Value = item
		return
	}

	if len(c.items) >= c.capacity {
		oldest := c.ll.Front()
		item := oldest.Value.(entry[K, V])

		delete(c.items, item.key)
		c.ll.Remove(oldest)
	}

	element := c.ll.PushBack(entry[K, V]{
		key:   key,
		value: value,
	})

	c.items[key] = element
}
