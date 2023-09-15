package sliding_window_counter

import (
	"sync"
	"time"
)

type SlidingWindowCounterRateLimiter struct {
	windowSize  time.Duration
	interval    time.Duration
	maxRequests int
	counters    []int
	lastUpdated int
	mu          sync.Mutex
}

func NewSlidingWindowCounterRateLimiter(windowSize time.Duration, interval time.Duration, maxRequests int) *SlidingWindowCounterRateLimiter {
	intervals := int(windowSize.Seconds() / interval.Seconds())
	return &SlidingWindowCounterRateLimiter{
		windowSize:  windowSize,
		interval:    interval,
		maxRequests: maxRequests,
		counters:    make([]int, intervals),
	}
}

func (rl *SlidingWindowCounterRateLimiter) AllowRequest() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	truncated := now.Truncate(rl.interval)
	intervalSeconds := rl.interval.Seconds()
	subTime := now.Sub(truncated)
	currentInterval := int(subTime / time.Duration(intervalSeconds))

	// Remove expired intervals and update counters
	for i := rl.lastUpdated; i < currentInterval; i++ {
		rl.counters[i%len(rl.counters)] = 0
	}
	rl.lastUpdated = currentInterval

	// Check if the current interval has exceeded the limit
	if rl.counters[currentInterval%len(rl.counters)] < rl.maxRequests {
		rl.counters[currentInterval%len(rl.counters)]++
		return true
	}
	return false
}
