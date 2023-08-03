package memory

import (
	"context"
	"sync"
	"time"

	"github.com/hhk7734/ratelimit.go"
)

func Open() *memoryDriver {
	return &memoryDriver{
		store: make(map[string][]time.Time),
	}
}

var _ ratelimit.Driver = new(memoryDriver)

type memoryDriver struct {
	mu    sync.Mutex
	store map[string][]time.Time
}

func (r *memoryDriver) SlidingWindowLog(ctx context.Context, key string, limit int, window time.Duration) (int, time.Duration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// remove expired
	now := time.Now()
	times := r.store[key]
	for i := len(times) - 1; i >= 0; i-- {
		if times[i].Before(now.Add(-window)) {
			times = times[i+1:]
			break
		}
	}

	// check limit
	prevCount := len(times)
	if prevCount >= limit {
		reset := window - now.Sub(times[0])
		return 0, reset, ratelimit.ErrLimitExceeded
	}

	// add new
	r.store[key] = append(times, now)

	remaining := limit - prevCount - 1
	reset := window - now.Sub(r.store[key][0])
	return remaining, reset, nil
}
