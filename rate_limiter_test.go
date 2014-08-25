package rate_limiter

import (
	"testing"
	"time"
)

func TestWaitUnderLimit(t *testing.T) {
	start := time.Now()
	limit := 5
	interval := time.Second * 3
	limiter := New(limit, interval)
	for i := 0; i < limit; i++ {
		limiter.Wait()
	}
	if time.Now().Sub(start) >= interval {
		t.Error("The limiter blocked when it shouldn't have")
	}
}

func TestWaitOverLimit(t *testing.T) {
	start := time.Now()
	limit := 5
	interval := time.Second * 3
	limiter := New(limit, interval)
	for i := 0; i < limit+1; i++ {
		limiter.Wait()
	}
	if time.Now().Sub(start) < interval {
		t.Error("The limiter didn't block when it should have")
	}
}
