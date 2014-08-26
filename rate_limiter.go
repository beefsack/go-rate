package rate

import (
	"sync"
	"time"
)

// A RateLimiter limits the rate at which an action can be performed.
type RateLimiter struct {
	interval time.Duration
	mtx      sync.Mutex
	times    []time.Time
}

// New creates a new rate limiter for the limit and interval.
func New(limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		interval: interval,
		times:    make([]time.Time, 0, limit),
	}
}

// Wait will block if the rate limit has been reached.
func (r *RateLimiter) Wait() {
	for {
		ok, remaining := r.Try()
		if ok {
			break
		}
		time.Sleep(remaining)
	}
}

// Try will return true if under the rate limit, or false if over and the
// remaining time before the rate limit expires.
func (r *RateLimiter) Try() (ok bool, remaining time.Duration) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	now := time.Now()
	if l := len(r.times); l == cap(r.times) {
		if diff := now.Sub(r.times[0]); diff < r.interval {
			return false, r.interval - diff
		}
		copy(r.times, r.times[1:])
		r.times = r.times[:l-1]
	}
	r.times = append(r.times, now)
	return true, 0
}
