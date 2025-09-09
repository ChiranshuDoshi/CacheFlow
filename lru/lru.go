package lru

import (
	"container/list"

	"github.com/ChiranshuDoshi/CacheFlow/cache"
)

type item struct {
	key   string
	value cache.Value
	size  int64
}

type LRUCache struct {
	capacity int64
	size     int64
	ls       *list.List
	table    map[string]*list.Element
}

// New creates a new LRU cache with given capacity (in bytes)
func New(capacity int64) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		size:     0,
		ls:       list.New(),
		table:    make(map[string]*list.Element),
	}
}

// Put adds a key-value pair
func (c *LRUCache) Put(key string, value cache.Value) {
	if entry := c.table[key]; entry != nil {
		// Key already exists, update the value
		it := entry.Value.(*item)
		c.size += value.Size() - it.size
		it.value = value
		it.size = value.Size()
		c.ls.MoveToBack(entry) // Mark as most recently used
	} else {
		// New key, add to cache
		it := &item{
			key:   key,
			value: value,
			size:  value.Size(),
		}
		c.table[key] = c.ls.PushBack(it)
		c.size += it.size
	}
	c.evictLRU()
}

// Get retrieves a value and marks it as recently used
func (c *LRUCache) Get(key string) (cache.Value, bool) {
	entry := c.table[key]
	if entry == nil {
		return nil, false
	}
	it := entry.Value.(*item)
	c.ls.MoveToBack(entry) // Mark as most recently used
	return it.value, true
}

// evictLRU removes least recently used items if over capacity
func (c *LRUCache) evictLRU() {
	for c.size > c.capacity {
		front := c.ls.Front()
		if front == nil {
			return
		}
		it := front.Value.(*item)
		c.ls.Remove(front)
		delete(c.table, it.key)
		c.size -= it.size
	}
}

// List returns current cache content
func (c *LRUCache) List() []map[string]cache.Value {
	var listContent []map[string]cache.Value
	for key, entry := range c.table {
		it := entry.Value.(*item)
		listContent = append(listContent, map[string]cache.Value{
			key: it.value,
		})
	}
	return listContent
}
