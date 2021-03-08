package cache

type CacheReader interface {
	Read(k int) (interface{}, bool)
}

type CacheWriter interface {
	Write(k int, v interface{})
}

type Cache interface {
	CacheReader
	CacheWriter
}
