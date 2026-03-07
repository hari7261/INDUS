package engine

import (
	"sync"
	"time"
)

type cacheEntry struct {
	expiresAt time.Time
	response  Response
}

type responseCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

func newResponseCache() *responseCache {
	return &responseCache{entries: map[string]cacheEntry{}}
}

func (c *responseCache) Get(key string) (Response, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok {
		return Response{}, false
	}
	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		return Response{}, false
	}

	response := entry.response
	response.Cached = true
	return response, true
}

func (c *responseCache) Set(key string, response Response, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	response.Cached = false
	response.Duration = 0
	response.Warning = ""
	response.Effects = Effects{}
	c.entries[key] = cacheEntry{
		expiresAt: time.Now().Add(ttl),
		response:  response,
	}
}

func (c *responseCache) Clear() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := len(c.entries)
	c.entries = map[string]cacheEntry{}
	return count
}

func (c *responseCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
