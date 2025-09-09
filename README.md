# CacheFlow

A high-performance, flexible caching library for Go that supports multiple eviction policies and TTL (Time-To-Live) functionality.

## Features

- **LRU (Least Recently Used) Cache**: Memory-efficient cache with configurable capacity
- **TTL Cache**: Time-based expiration with automatic cleanup
- **Size-aware**: Tracks memory usage for intelligent eviction
- **Thread-safe operations**: Ready for concurrent applications
- **Clean interfaces**: Easy to extend and customize

## Installation

```bash
go get github.com/ChiranshuDoshi/CacheFlow
```

## Quick Start

### LRU Cache Example

```go
package main

import (
    "fmt"
    "github.com/ChiranshuDoshi/CacheFlow/lru"
)

type IntValue int64
func (i IntValue) Size() int64 { return 8 } // 8 bytes for int64

func main() {
    // Create LRU cache with 16 bytes capacity (holds 2 int64 values)
    cache := lru.New(16)
    
    // Add items
    cache.Put("key1", IntValue(100))
    cache.Put("key2", IntValue(200))
    
    // Retrieve item
    if val, ok := cache.Get("key1"); ok {
        fmt.Printf("Found: %v\n", val)
    }
}
```

### TTL Cache Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/ChiranshuDoshi/CacheFlow/ttlcache"
)

func main() {
    cache := ttlcache.New()
    
    // Add item with 5-second TTL
    cache.Put("session", IntValue(12345), 5*time.Second)
    
    // Item available immediately
    if val, ok := cache.Get("session"); ok {
        fmt.Printf("Session active: %v\n", val)
    }
    
    // Wait for expiration
    time.Sleep(6 * time.Second)
    
    // Item expired
    if _, ok := cache.Get("session"); !ok {
        fmt.Println("Session expired")
    }
}
```

## API Reference

### Cache Interface

All caches implement the basic `Cache` interface:

```go
type Cache interface {
    Get(key string) (Value, bool)
    Put(key string, value Value)
}

type Value interface {
    Size() int64
}
```

### LRU Cache

#### Constructor
- `lru.New(capacity int64)` - Creates new LRU cache with byte-based capacity

#### Methods
- `Put(key string, value cache.Value)` - Adds or updates a key-value pair
- `Get(key string) (cache.Value, bool)` - Retrieves value and marks as recently used
- `List() []map[string]cache.Value` - Returns all cached items

#### Features
- **Capacity Management**: Automatically evicts least recently used items when capacity exceeded
- **Access Tracking**: Items are moved to "most recent" position on access
- **Memory Awareness**: Uses actual byte size of values for eviction decisions

### TTL Cache

#### Constructor
- `ttlcache.New()` - Creates new TTL cache

#### Methods
- `Put(key string, value cache.Value, ttl time.Duration)` - Adds item with expiration time
- `Get(key string) (cache.Value, bool)` - Retrieves value if not expired
- `List() []map[string]cache.Value` - Returns all non-expired items

#### Features
- **Automatic Expiration**: Items expire after specified duration
- **Lazy Cleanup**: Expired items removed on access
- **Flexible TTL**: Set different expiration times per item
- **No Expiration**: Use `ttl <= 0` for permanent storage

## Advanced Usage

### Custom Value Types

Implement the `cache.Value` interface for your custom types:

```go
type User struct {
    ID   int
    Name string
    Data []byte
}

func (u User) Size() int64 {
    return int64(8 + len(u.Name) + len(u.Data)) // Approximate size
}

// Usage
cache := lru.New(1024) // 1KB capacity
user := User{ID: 1, Name: "Alice", Data: make([]byte, 100)}
cache.Put("user:1", user)
```

### Capacity Planning

For LRU cache capacity planning:

```go
// For storing N items of size S bytes each:
capacity := int64(N * S)

// Example: 100 int64 values
cache := lru.New(100 * 8) // 800 bytes
```

## Example Output

```
=== LRU Eviction Demo ===
Cache state: [map[a:10] map[b:20]]
Cache after inserting c: [map[a:10] map[c:30]]
Cache after inserting d: [map[c:30] map[d:40]]

=== TTL Expiration Demo ===
Initial Cache: [map[k1:111] map[k2:222]]
Get k1: 111 true
Sleeping 6s to let TTL expire...
After TTL, Get k1: <nil> false
Cache after TTL: []
```

## Project Structure

```
CacheFlow/
├── cache/          # Core interfaces
│   └── cache.go
├── lru/            # LRU implementation
│   └── lru.go
├── ttlcache/       # TTL implementation
│   └── ttlcache.go
├── main.go         # Demo examples
└── README.md
```

## Performance Characteristics

### LRU Cache
- **Get**: O(1) - Hash table lookup + doubly-linked list move
- **Put**: O(1) - Hash table insert + list operations
- **Space**: O(n) where n is number of items

### TTL Cache
- **Get**: O(1) - Hash table lookup with expiration check
- **Put**: O(1) - Hash table insert
- **Space**: O(n) where n is number of items (includes expired items until accessed)

## Use Cases

- **Web Session Management**: TTL cache for user sessions
- **API Response Caching**: LRU cache for frequently accessed data
- **Database Query Caching**: Combine both for optimal performance
- **Rate Limiting**: TTL cache for request tracking
- **Memory Management**: LRU cache for bounded memory usage

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


## Acknowledgments

- Inspired by Redis and Memcached caching strategies
- Built with Go's standard library for optimal performance
- Thanks to the Go community for best practices and patterns
