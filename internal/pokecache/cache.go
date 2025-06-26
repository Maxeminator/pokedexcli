package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{}
	c.cacheMap = make(map[string]cacheEntry)
	c.interval = interval
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ce := cacheEntry{}
	ce.createdAt = time.Now()
	ce.val = val
	c.cacheMap[key] = ce
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cacheMap[key]
	if ok {
		return entry.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()

		for key, entry := range c.cacheMap {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.cacheMap, key)
			}
		}

		c.mu.Unlock()
	}
}
