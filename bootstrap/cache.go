package bootstrap

import (
	_redis "github.com/LXJ0000/go-backend/redis"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRedisCache(env *Env) _redis.Cache {
	cmd := redis.NewClient(&redis.Options{
		Addr: env.RedisAddr,
	})
	return _redis.NewCache(cmd, time.Duration(env.RedisExpiration)*time.Minute)
}
