package cache

type Value interface {
	Size() int64
}

type Cache interface {
	Get(key string) (Value, bool)
	Put(key string, value Value)
}
