package cache

import (
	"context"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

type RistrettoCache struct {
	cache *ristretto.Cache
	mu    sync.RWMutex
}

func NewRistrettoCache(cache *ristretto.Cache) *RistrettoCache {
	return &RistrettoCache{
		cache: cache,
		mu:    sync.RWMutex{},
	}
}

func (c *RistrettoCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.cache.SetWithTTL(key, value, 1, expiration)
	return nil
}

func (c *RistrettoCache) Get(ctx context.Context, key string) (string, error) {
	value, found := c.cache.Get(key)
	if !found {
		return "", nil
	}
	return value.(string), nil
}

func (c *RistrettoCache) Del(ctx context.Context, key string) error {
	c.cache.Del(key)
	return nil
}

func (c *RistrettoCache) SetNx(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	// 无锁检查
	_, found := c.cache.Get(key)
	if found {
		return false, nil
	}

	// 加锁进行设置操作
	c.mu.Lock()
	defer c.mu.Unlock()

	// 再次检查，确保在加锁前没有其他协程设置了相同的键
	_, found = c.cache.Get(key)
	if found {
		return false, nil
	}

	return c.cache.SetWithTTL(key, value, 1, expiration), nil
}
