package leaky_bucket

import "time"

type Request struct {
	x int
}

type LeakyBucketRateLimiter struct {
	rate              int
	duration          time.Duration
	queue             chan Request
	maxBucketCapacity int
}

func NewLeakyBucketRateLimiter(rate int, maxBucketCapacity int, duration time.Duration) *LeakyBucketRateLimiter {
	return &LeakyBucketRateLimiter{
		rate:     rate,
		queue:    make(chan Request, maxBucketCapacity),
		duration: duration,
	}
}

func (l *LeakyBucketRateLimiter) AddRequest(req int) {
	l.queue <- Request{
		x: req,
	}
}

func (l *LeakyBucketRateLimiter) ConsumeRequest() Request {
	ticker := time.NewTicker(l.duration / time.Duration(l.rate))
	defer ticker.Stop()

	select {
	case <-ticker.C:
		return <-l.queue
	}
}
