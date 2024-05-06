package cache

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestCacheHSet(t *testing.T) {
	cmd := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	_, err := cmd.Ping(context.Background()).Result()
	require.NoError(t, err)

	mp := map[string]interface{}{
		"name":  "Jannan",
		"age":   18,
		"email": "122@qq.com",
	}
	err = cmd.HSet(context.Background(), "user", mp).Err()
	require.NoError(t, err)

	res, err := cmd.HGetAll(context.Background(), "user").Result()
	require.NoError(t, err)

	t.Log(res)
}

func TestCacheHMSet(t *testing.T) {
	cmd := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	_, err := cmd.Ping(context.Background()).Result()
	require.NoError(t, err)

	mp := map[string]string{
		"name":  "Jannan",
		"age":   "18",
		"email": "122@qq.com",
	}
	err = cmd.HMSet(context.Background(), "user1", mp).Err()
	require.NoError(t, err)

	res, err := cmd.HGetAll(context.Background(), "user").Result()
	require.NoError(t, err)

	t.Log(res)
}
