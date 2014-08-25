package rate

import (
	"sync"
	"time"
)

// A rateLimiter limits the rate at which an action can be performed.
type rateLimiter struct {
	limit               int
	interval            time.Duration
	waitMutex, tryMutex *sync.Mutex
	wLock               *sync.Mutex
	times               []time.Time
}

// New creates a new rate limiter for the limit and interval.
func New(limit int, interval time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:     limit,
		interval:  interval,
		waitMutex: &sync.Mutex{},
		tryMutex:  &sync.Mutex{},
		times:     make([]time.Time, limit),
	}
}

// Wait will block if the rate limit has been reached.
func (r *rateLimiter) Wait() {
	r.waitMutex.Lock()
	defer r.waitMutex.Unlock()
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
func (r *rateLimiter) Try() (ok bool, remaining time.Duration) {
	r.tryMutex.Lock()
	defer r.tryMutex.Unlock()
	if len(r.times) == r.limit {
		diff := time.Now().Sub(r.times[0])
		if diff < r.interval {
			return false, r.interval - diff
		}
		r.times = r.times[1:]
	}
	r.times = append(r.times, time.Now())
	return true, 0
}
