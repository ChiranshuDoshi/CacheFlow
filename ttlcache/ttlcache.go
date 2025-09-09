package ttlcache

import (
	"time"

	"github.com/ChiranshuDoshi/CacheFlow/cache"
)

type item struct {
	value  cache.Value
	expiry int64
}

type TTLCache struct {
	table map[string]*item
}

// New creates a new TTL cache
func New() *TTLCache {
	return &TTLCache{
		table: make(map[string]*item),
	}
}

// Put adds a key-value pair with TTL
func (c *TTLCache) Put(key string, value cache.Value, ttl time.Duration) {
	it := &item{
		value: value,
	}
	if ttl > 0 {
		it.expiry = time.Now().Add(ttl).UnixNano()
	} else {
		it.expiry = 0 // No expiration if TTL <= 0
	}
	c.table[key] = it
}

// Get retrieves a value and respects TTL
func (c *TTLCache) Get(key string) (cache.Value, bool) {
	it, exists := c.table[key]
	if !exists {
		return nil, false
	}

	// Check if item has expired
	if it.expiry > 0 && time.Now().UnixNano() > it.expiry {
		delete(c.table, key) // Clean up expired item
		return nil, false
	}

	return it.value, true
}

// List returns current cache content, skipping expired items
func (c *TTLCache) List() []map[string]cache.Value {
	var listContent []map[string]cache.Value
	now := time.Now().UnixNano()

	for key, it := range c.table {
		// Check if item has expired
		if it.expiry > 0 && now > it.expiry {
			delete(c.table, key) // Clean up expired item
			continue
		}

		listContent = append(listContent, map[string]cache.Value{
			key: it.value,
		})
	}
	return listContent
}
