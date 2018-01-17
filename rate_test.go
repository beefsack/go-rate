package rate

import (
	"testing"
	"time"
	"fmt"
	"sync"
)

func BenchmarkAtomicRateLimiter(b *testing.B) {
	for ng := 1; ng <= 2048; ng *= 2 {
		b.Run(fmt.Sprint(ng), func(b *testing.B) {
			limiter := New(b.N, 10)

			var wg sync.WaitGroup
			wg.Add(ng)

			n := b.N
			quota := n / ng

			for g := ng; g > 0; g-- {
				if g == 1 {
					quota = n
				}

				go func(quota int) {
					for i := 0; i < quota; i++ {
						ok, _ := limiter.Try()
						if !ok {
							b.Fatal("No enough permissions")
						}
					}
					wg.Done()
				}(quota)

				n -= quota
			}

			if n != 0 {
				b.Fatalf("Incorrect quota assignments: %v remaining", n)
			}

			b.StartTimer()
			wg.Wait()
			b.StopTimer()
		})
	}
}

func BenchmarkMutexRateLimiter(b *testing.B) {
	for ng := 1; ng <= 2048; ng *= 2 {
		b.Run(fmt.Sprint(ng), func(b *testing.B) {
			limiter := NewMutexRateLimiter(b.N, 10)

			var wg sync.WaitGroup
			wg.Add(ng)

			n := b.N
			quota := n / ng

			for g := ng; g > 0; g-- {
				if g == 1 {
					quota = n
				}

				go func(quota int) {
					for i := 0; i < quota; i++ {
						ok, _ := limiter.Try()
						if !ok {
							b.Fatal("No enough permissions")
						}
					}
					wg.Done()
				}(quota)

				n -= quota
			}

			if n != 0 {
				b.Fatalf("Incorrect quota assignments: %v remaining", n)
			}

			b.StartTimer()
			wg.Wait()
			b.StopTimer()
		})
	}
}

func TestRateLimiter_Wait_noblock(t *testing.T) {
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

func TestRateLimiter_Wait_block(t *testing.T) {
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

func TestRateLimiter_Try(t *testing.T) {
	limit := 5
	interval := time.Second * 3
	limiter := New(limit, interval)
	for i := 0; i < limit; i++ {
		if ok, _ := limiter.Try(); !ok {
			t.Fatalf("Should have allowed try on attempt %d", i)
		}
	}
	if ok, _ := limiter.Try(); ok {
		t.Fatal("Should have not allowed try on final attempt")
	}
}
