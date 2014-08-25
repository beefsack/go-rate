package rate_limiter

import (
	"sync"
	"time"
)

// A RateLimiter limits the rate at which an action can be performed.
type RateLimiter struct {
	limit    int
	interval time.Duration
	lock     *sync.Mutex
	times    []time.Time
}

// New creates a new rate limiter for the limit and interval.
func New(limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		interval: interval,
		lock:     &sync.Mutex{},
		times:    make([]time.Time, limit),
	}
}

// Wait will block if the rate limit has been reached.
func (r *RateLimiter) Wait() {
	r.lock.Lock()
	for len(r.times) == r.limit {
		diff := time.Now().Sub(r.times[0])
		if diff < r.interval {
			time.Sleep(r.interval - diff)
		}
		r.times = r.times[1:]
	}
	r.times = append(r.times, time.Now())
	r.lock.Unlock()
}
