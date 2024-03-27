package cache

import (
	"practice/lru"
	"sync"
)

/*
* cache是主要的缓存的结构体，lru只是实现了lru策略的一个数据结构，主要还是这个并发包装过的cache结构体
 */

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return value, ok
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return value, ok
}
