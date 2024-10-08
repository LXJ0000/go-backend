package bootstrap

import (
	"log"
	"log/slog"
	"time"

	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func NewRedisCache(env *Env) cache.RedisCache {
	cmd := redis.NewClient(&redis.Options{
		Addr:     env.RedisAddr,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})
	if _, err := cmd.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}
	return cache.NewRedisCache(cmd, time.Duration(env.RedisExpiration)*time.Minute)
}

func NewLocalCache(env *Env) cache.LocalCache {
	slog.Error("LocalCache is nil")
	return cache.LocalCache{} // TODO
}
