package project_for_practice

import (
	"fmt"
	"geecache/lru"
	"sync"
)

// 实例化lru,封装mu

type cache struct {
	mu			sync.Mutex
	lru 		*lru.Cache
	cacheBytes	int64
}

func (c *cache) add(key string, value ByteView)  {
	// 需要加锁
	c.mu.Lock()
	defer c.mu.Unlock()
	// 没有就先建立
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)

}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 没有
	if c.lru == nil {
		fmt.Println("no element")
		return
	}
	// 拿到
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
