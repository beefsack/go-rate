package rate

import (
	"time"
	"unsafe"
	"sync/atomic"
)

type state struct {
	cycle       int64
	permissions int64
}

// A RateLimiter limits the rate at which an action can be performed.  It
// applies neither smoothing (like one could achieve in a token bucket system)
// nor does it offer any conception of warmup, wherein the rate of actions
// granted are steadily increased until a steady throughput equilibrium is
// reached.
type RateLimiter struct {
	start    time.Time
	limit    int64
	interval int64
	state    unsafe.Pointer
}

// New creates a new rate limiter for the limit and interval.
func New(limit int, interval time.Duration) *RateLimiter {
	start := time.Now()
	newState := state{
		cycle:       0,
		permissions: int64(limit),
	}
	lim := &RateLimiter{
		start:    start,
		limit:    int64(limit),
		interval: int64(interval),
		state:    unsafe.Pointer(&newState),
	}

	return lim
}

// Wait blocks if the rate limit has been reached.  Wait offers no guarantees
// of fairness for multiple actors if the allowed rate has been temporarily
// exhausted.
func (r *RateLimiter) Wait() {
	for {
		ok, remaining := r.Try()
		if ok {
			break
		}
		time.Sleep(remaining)
	}
}

// Try returns true if under the rate limit, or false if over and the
// remaining time before the rate limit expires.
func (r *RateLimiter) Try() (ok bool, remaining time.Duration) {
	for {
		//extract previous state
		previousStatePointer := atomic.LoadPointer(&r.state)
		previousState := (*state)(previousStatePointer)
		// compute new state
		now := int64(time.Now().Sub(r.start))
		currentCycle := now / r.interval
		newState := state{
			cycle:       currentCycle,
			permissions: previousState.permissions,
		}
		// count elapsed cycles and produce more permissions if necessary
		elapsedCycles := currentCycle - previousState.cycle
		if elapsedCycles > 0 {
			permissionsToAppear := elapsedCycles * r.limit
			newState.permissions = min(previousState.permissions+permissionsToAppear, r.limit)
		}
		// try to acquire permission by atomic update
		if newState.permissions > 0 {
			newState.permissions -= 1
			if atomic.CompareAndSwapPointer((*unsafe.Pointer)(&r.state), previousStatePointer, unsafe.Pointer(&newState)) {
				return true, 0
			}
			continue
		}
		// if there is not enough permissions calculate wait duration
		nextCycleStart := (currentCycle + 1) * int64(r.interval)
		nsToNextCycle := nextCycleStart - now
		fullCyclesRequired := (-newState.permissions) / r.interval
		remaining := (fullCyclesRequired * r.interval) + nsToNextCycle
		return false, time.Duration(remaining)
	}
}

func min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
