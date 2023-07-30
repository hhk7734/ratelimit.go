package ratelimit

import (
	"context"
	"time"
)

type RateLimit interface {
	SlidingWindowLog(ctx context.Context, key string, limit int, window time.Duration) (int, time.Duration, error)
}
