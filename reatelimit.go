package ratelimit

import (
	"context"
	"time"
)

func New(d Driver) *RateLimit {
	return &RateLimit{d: d}
}

type RateLimit struct {
	d Driver
}

func (r *RateLimit) SlidingWindowLog(ctx context.Context, key string, limit int, window time.Duration) (int, time.Duration, error) {
	return r.d.SlidingWindowLog(ctx, key, limit, window)
}
