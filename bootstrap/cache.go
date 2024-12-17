package bootstrap

import (
	"log"
	"time"

	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/dgraph-io/ristretto"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func NewRedisCache(env *Env) *cache.RedisCache {
	cmd := redis.NewClient(&redis.Options{
		Addr:     env.RedisHost + ":" + env.RedisPort,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})
	if _, err := cmd.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}
	return cache.NewRedisCache(cmd, time.Duration(env.RedisExpiration)*time.Minute)
}

func NewLocalCache(env *Env) *cache.RistrettoCache {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		log.Fatal(err)
	}
	// defer cache.Close()
	return cache.NewRistrettoCache(c)
}
