package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu sync.Mutex
	Entries map[string]*CacheEntry
}

type CacheEntry struct {
	createAt time.Time
	val []byte
}

func NewCache(t time.Duration) *Cache {
	c := &Cache{
		Entries: map[string]*CacheEntry{},
	}
	go c.reapLoop(t)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	if _, ok := c.Entries[key]; !ok {
		 c.Entries[key] = &CacheEntry{
			createAt: time.Now(),
			val: val, 
		}
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if _, ok := c.Entries[key]; ok {
		return c.Entries[key].val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(d time.Duration) bool {
	tick := time.NewTicker(d)
	for {
		<-tick.C
		c.mu.Lock()
		for k, _ := range c.Entries {	
			delete(c.Entries, k)
		}
		c.mu.Unlock()
	}
}
