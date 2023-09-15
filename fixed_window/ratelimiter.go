package fixed_window

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	rate                      int
	duration                  time.Duration
	requestProcessed          int
	firstRequestProcessedTime time.Time
	syncMutex                 *sync.Mutex
}

func NewFixedWindowRateLimiter(
	rate int,
	duration time.Duration,
) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		rate:      rate,
		duration:  duration,
		syncMutex: &sync.Mutex{},
	}
}

func (f *FixedWindowRateLimiter) AllowRequest() bool {
	f.syncMutex.Lock()

	if f.requestProcessed == 0 {
		f.firstRequestProcessedTime = time.Now()
	}

	if time.Since(f.firstRequestProcessedTime) <= f.duration {
		if f.requestProcessed < f.rate {
			f.requestProcessed++
			f.syncMutex.Unlock()
			return true
		}
		f.syncMutex.Unlock()
		return false
	}

	f.requestProcessed = 0
	f.syncMutex.Unlock()

	return f.AllowRequest()
}

func (f *FixedWindowRateLimiter) AllowMessage() {
	f.syncMutex.Lock()

	if f.requestProcessed == 0 {
		f.firstRequestProcessedTime = time.Now()
	}

	if time.Since(f.firstRequestProcessedTime) <= f.duration {
		if f.requestProcessed < f.rate {
			f.requestProcessed++
			f.syncMutex.Unlock()
		} else {
			waitTime := f.duration - time.Since(f.firstRequestProcessedTime)
			time.Sleep(waitTime)
			f.requestProcessed = 0
			f.syncMutex.Unlock()
			f.AllowMessage()
		}
		return
	}

	f.requestProcessed = 0
	f.syncMutex.Unlock()
	f.AllowMessage()

	return
}
