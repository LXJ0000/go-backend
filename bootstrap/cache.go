package bootstrap

import (
	_redis "github.com/LXJ0000/go-backend/redis"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log"
	"time"
)

func NewRedisCache(env *Env) _redis.Cache {
	cmd := redis.NewClient(&redis.Options{
		Addr: env.RedisAddr,
	})
	if _, err := cmd.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}
	return _redis.NewCache(cmd, time.Duration(env.RedisExpiration)*time.Minute)
}
