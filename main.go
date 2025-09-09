package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/ChiranshuDoshi/CacheFlow/lru"
	"github.com/ChiranshuDoshi/CacheFlow/ttlcache"
)

// IntValue implements cache.Value
type IntValue int64

func (i IntValue) Size() int64 {
	return int64(unsafe.Sizeof(i))
}

func main() {
	// ---------------------------
	// PART 1: LRU EVICTION DEMO
	// ---------------------------
	fmt.Println("=== LRU Eviction Demo ===")
	// Fixed: Use capacity that can hold 2 int64 values (2 * 8 = 16 bytes)
	lruCache := lru.New(16)

	lruCache.Put("a", IntValue(10))
	lruCache.Put("b", IntValue(20))
	fmt.Println("Cache state:", lruCache.List())

	lruCache.Get("a") // access a to make it recently used

	lruCache.Put("c", IntValue(30)) // should evict b (least recently used)
	fmt.Println("Cache after inserting c:", lruCache.List())

	lruCache.Put("d", IntValue(40)) // should evict a (now least recently used)
	fmt.Println("Cache after inserting d:", lruCache.List())

	// ---------------------------
	// PART 2: TTL EXPIRATION DEMO
	// ---------------------------
	fmt.Println("\n=== TTL Expiration Demo ===")
	ttlCache := ttlcache.New()
	ttl := 5 * time.Second

	ttlCache.Put("k1", IntValue(111), ttl)
	ttlCache.Put("k2", IntValue(222), 2*ttl)
	fmt.Println("Initial Cache:", ttlCache.List())

	val, ok := ttlCache.Get("k1")
	fmt.Println("Get k1:", val, ok)

	fmt.Println("Sleeping 6s to let TTL expire...")
	time.Sleep(6 * time.Second)

	val, ok = ttlCache.Get("k1")
	fmt.Println("After TTL, Get k1:", val, ok)
	fmt.Println("Cache after TTL:", ttlCache.List())
}
