package redis

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/hhk7734/ratelimit.go"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
)

func init() {
	slidingWindowLogLuaScript = redis.NewScript(slidingWindowLogLua)
}

var _ ratelimit.Driver = new(Driver)

type Driver struct {
	Client *redis.Client
}

func Open(opt *redis.Options) *Driver {
	client := redis.NewClient(opt)
	return &Driver{Client: client}
}

//go:embed slidingwindolog.lua
var slidingWindowLogLua string
var slidingWindowLogLuaScript *redis.Script

func (r *Driver) SlidingWindowLog(ctx context.Context, key string, limit int,
	window time.Duration) (int, time.Duration, error) {
	ret, err := slidingWindowLogLuaScript.Run(ctx, r.Client,
		[]string{r.key(key)},
		limit,
		window.Milliseconds(),
		time.Now().UnixMilli(),
		ulid.Make().String()).Int64Slice()
	if err != nil {
		return 0, 0, err
	}

	remaining := int(ret[0])
	reset := time.Duration(ret[1]) * time.Millisecond

	if remaining < 0 {
		return 0, reset, ratelimit.ErrLimitExceeded
	}

	return remaining, reset, nil
}

func (*Driver) key(key string) string {
	return fmt.Sprintf("ratelimit:%s", key)
}
