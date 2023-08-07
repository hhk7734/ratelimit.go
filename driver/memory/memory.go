package memory

import (
	"context"
	"sync"
	"time"

	"github.com/hhk7734/ratelimit.go"
)

var _ ratelimit.Driver = new(Driver)

type Driver struct {
	Store map[string][]time.Time

	mu sync.Mutex
}

func Open() *Driver {
	return &Driver{
		Store: make(map[string][]time.Time),
	}
}

func (r *Driver) SlidingWindowLog(ctx context.Context, key string, limit int, window time.Duration) (int, time.Duration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// remove expired
	now := time.Now()
	times := r.Store[key]
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
	r.Store[key] = append(times, now)

	remaining := limit - prevCount - 1
	reset := window - now.Sub(r.Store[key][0])
	return remaining, reset, nil
}
