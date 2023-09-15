package sliding_window_log

import (
	"sync"
	"time"
)

type SlidingWindowLogRateLimiter struct {
	windowSize     time.Duration
	maxRequests    int
	syncMutex      *sync.Mutex
	requestHistory []time.Time
}

func NewSlidingWindowLogRateLimiter(
	windowSize time.Duration,
	maxRequests int,
) *SlidingWindowLogRateLimiter {
	return &SlidingWindowLogRateLimiter{
		windowSize:     windowSize,
		maxRequests:    maxRequests,
		requestHistory: make([]time.Time, 0),
		syncMutex:      &sync.Mutex{},
	}
}

func (s *SlidingWindowLogRateLimiter) AllowRequest() bool {
	s.syncMutex.Lock()
	defer s.syncMutex.Unlock()

	indexToStart := 0
	for index, history := range s.requestHistory {
		if time.Since(history) >= s.windowSize {
			indexToStart = index
		}
	}

	s.requestHistory = s.requestHistory[indexToStart:]
	if len(s.requestHistory) > s.maxRequests {
		return false
	}

	s.requestHistory = append(s.requestHistory, time.Now())
	return true
}
