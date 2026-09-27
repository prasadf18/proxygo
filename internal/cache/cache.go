package cache

import (
	"sync"
	"time"
)

type entry struct {
	data      []byte
	expiresAt time.Time
}

type Cache struct {
	mu      sync.Mutex
	items   map[string]entry
	ttl     time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]entry),
		ttl:   ttl,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, found := c.items[key]
	if !found {
		return nil, false
	}
	if time.Now().After(e.expiresAt) {
		delete(c.items, key)
		return nil, false
	}
	return e.data, true
}

func (c *Cache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = entry{
		data:      data,
		expiresAt: time.Now().Add(c.ttl),
	}
}