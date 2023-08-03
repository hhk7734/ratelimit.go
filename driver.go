package ratelimit

import (
	"context"
	"time"
)

type Driver interface {
	SlidingWindowLog(ctx context.Context, key string, limit int, window time.Duration) (int, time.Duration, error)
}
