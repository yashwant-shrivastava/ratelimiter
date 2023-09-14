package token_bucket

import (
	"sync"
	"time"
)

// TokenBucketRequestRateLimiter to either allow or discard the request
type TokenBucketRequestRateLimiter struct {
	rate                     int
	tokens                   chan struct{}
	intervalLimit            int
	burstLimit               int
	refillInterval           time.Duration
	refillBurstLimitInterval time.Duration
	mu                       sync.Mutex
}

func NewTokenBucketRequestRateLimiter(
	burstLimit,
	rate int,
	refillInterval,
	refillBurstLimitInterval time.Duration,
) *TokenBucketRequestRateLimiter {
	limiter := &TokenBucketRequestRateLimiter{
		tokens:                   make(chan struct{}, burstLimit),
		intervalLimit:            burstLimit,
		burstLimit:               burstLimit,
		refillInterval:           refillInterval,
		rate:                     rate,
		refillBurstLimitInterval: refillBurstLimitInterval,
	}

	go limiter.refillTokens()
	go limiter.refillBurstLimit()
	return limiter
}

func (rl *TokenBucketRequestRateLimiter) refillTokens() {
	ticker := time.NewTicker(rl.refillInterval / time.Duration(rl.rate))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			rl.tokens <- struct{}{}
			rl.mu.Unlock()
		}
	}
}

func (rl *TokenBucketRequestRateLimiter) refillBurstLimit() {
	ticker := time.NewTicker(rl.refillBurstLimitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for i := len(rl.tokens); i < rl.burstLimit; i++ {
				rl.tokens <- struct{}{}
			}
		}
	}
}

func (rl *TokenBucketRequestRateLimiter) AllowRequest() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func (rl *TokenBucketRequestRateLimiter) AllowMessage() bool {
	select {
	case <-rl.tokens:
		return true
	}
}
