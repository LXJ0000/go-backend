package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Lua(ctx context.Context, luaPath string, key []string, args ...string) (int, error)
}

type cache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewCache(cmd redis.Cmdable, expiration time.Duration) Cache {
	return &cache{cmd: cmd, expiration: expiration}
}

func (c *cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.cmd.Set(ctx, key, value, expiration).Err()
}
func (c *cache) Get(ctx context.Context, key string) (string, error) {
	return c.cmd.Get(ctx, key).Result()
}

func (c *cache) Del(ctx context.Context, key string) error {
	return c.cmd.Del(ctx, key).Err()
}

func (c *cache) Lua(ctx context.Context, luaPath string, key []string, args ...string) (int, error) {
	return c.cmd.Eval(ctx, luaPath, key, args).Int()
}
