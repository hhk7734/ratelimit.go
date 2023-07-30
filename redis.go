package ratelimit

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

func init() {
	slidingWindowLogLuaScript = redis.NewScript(slidingWindowLogLua)
}

func NewRedisRateLimit(client *redis.Client) *RedisRateLimit {
	return &RedisRateLimit{
		client: client,
	}
}

var _ RateLimit = (*RedisRateLimit)(nil)

type RedisRateLimit struct {
	client *redis.Client
}

//go:embed slidingwindolog.lua
var slidingWindowLogLua string
var slidingWindowLogLuaScript *redis.Script

func (r *RedisRateLimit) SlidingWindowLog(ctx context.Context, key string, limit int,
	window time.Duration) (int, time.Duration, error) {
	ret, err := slidingWindowLogLuaScript.Run(ctx, r.client, []string{key}, limit, window.Milliseconds(), time.Now().UnixMilli()).Int64Slice()
	if err != nil {
		return 0, 0, err
	}

	remaining := int(ret[0])
	reset := time.Duration(ret[1]) * time.Millisecond

	if remaining < 0 {
		return 0, reset, ErrLimitExceeded
	}

	return remaining, reset, nil
}
