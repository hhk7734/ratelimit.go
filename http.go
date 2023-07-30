package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func NewHttpRateLimit(ratelimit RateLimit) *HttpRateLimit {
	return &HttpRateLimit{
		ratelimit: ratelimit,
	}
}

type HttpRateLimit struct {
	ratelimit RateLimit
}

func (r *HttpRateLimit) SlidingWindowLog(ctx context.Context, w http.ResponseWriter, key string, limit int,
	window time.Duration) error {
	remaining, reset, err := r.ratelimit.SlidingWindowLog(ctx, key, limit, window)
	if err != nil && !errors.Is(err, ErrLimitExceeded) {
		return err
	}

	w.Header().Set("RateLimit-Policy",
		fmt.Sprintf("%d;w=%d;policy=\"sliding window log\"", limit, int(window.Seconds())))
	w.Header().Set("RateLimit",
		fmt.Sprintf("limit=%d, remaining=%d, reset=%d", limit, remaining, int(reset.Seconds())))
	return err
}
