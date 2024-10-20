package ratelimitutil

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var LuaSlideWindow string

type RedisSlideWindow struct {
	cmd    redis.Cmdable
	window time.Duration
	limit  int
}

func NewRedisSlideWindow(cmd redis.Cmdable, window time.Duration, limit int) Limiter {
	return &RedisSlideWindow{
		cmd:    cmd,
		window: window,
		limit:  limit,
	}
}

func (r *RedisSlideWindow) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, LuaSlideWindow, []string{key}, r.window, r.limit, time.Now().UnixMilli()).Bool()
}
