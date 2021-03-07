package cache

import (
	"sync"
)

type FeedCache struct {
	mu    sync.RWMutex
	items map[int]interface{}
}

func NewFeedCache() Cache {
	return &FeedCache{
		items: make(map[int]interface{}),
	}
}

func (f *FeedCache) Read(k int) (interface{}, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	v, ok := f.items[k]
	return v, ok
}

func (f *FeedCache) Write(k int, v interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.items[k] = v
}
