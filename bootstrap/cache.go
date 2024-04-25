package bootstrap

import (
	"github.com/LXJ0000/go-backend/cache"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log"
	"log/slog"
	"time"
)

func NewRedisCache(env *Env) cache.Cache {
	cmd := redis.NewClient(&redis.Options{
		Addr: env.RedisAddr,
	})
	if _, err := cmd.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}
	return cache.NewCache(cmd, time.Duration(env.RedisExpiration)*time.Minute)
}

func NewLocalCache(env *Env) cache.LocalCache {
	slog.Error("LocalCache is nil")
	return cache.LocalCache{} // TODO
}
