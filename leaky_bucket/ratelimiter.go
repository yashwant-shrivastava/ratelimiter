package leaky_bucket

import "time"

type Request struct {
	x int
}

type LeakyBucketRateLimiter struct {
	rate              int
	refillDuration    time.Duration
	queue             chan Request
	maxBucketCapacity int
}

func NewLeakyBucketRateLimiter(rate int, maxBucketCapacity int, refillDuration time.Duration) *LeakyBucketRateLimiter {
	return &LeakyBucketRateLimiter{
		rate:           rate,
		queue:          make(chan Request, maxBucketCapacity),
		refillDuration: refillDuration,
	}
}

func (l *LeakyBucketRateLimiter) AddRequest(req int) {
	l.queue <- Request{
		x: req,
	}
}

func (l *LeakyBucketRateLimiter) ConsumeRequest() Request {
	ticker := time.NewTicker(l.refillDuration / time.Duration(l.rate))
	defer ticker.Stop()

	select {
	case <-ticker.C:
		return <-l.queue
	}
}
