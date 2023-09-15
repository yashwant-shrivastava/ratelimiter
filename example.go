package main

import (
	"fmt"
	"github.com/yashwant-shrivastava/ratelimiter/fixed_window"
	"github.com/yashwant-shrivastava/ratelimiter/leaky_bucket"
	"github.com/yashwant-shrivastava/ratelimiter/sliding_window_log"
	"github.com/yashwant-shrivastava/ratelimiter/token_bucket"
	"time"
)

func TokenBucketRateLimiterExample() {
	burstCapacity := 10
	refillRate := 5
	refillInterval := 1 * time.Second
	refillBurstLimitInterval := 10 * time.Second

	limiter := token_bucket.NewTokenBucketEventRateLimiter(
		burstCapacity,
		refillRate,
		refillInterval,
		refillBurstLimitInterval,
	)

	// Simulate event messages
	for i := 0; i < 150; i++ {
		limiter.AllowMessage()
		fmt.Println("Request", i+1, "allowed.", time.Now().Unix())

	}

	// Simulate client requests
	for i := 0; i < 150; i++ {
		if limiter.AllowRequest() {
			fmt.Println("Request", i+1, "allowed.", time.Now().Unix())
		} else {
			fmt.Println("Request", i+1, "denied.", time.Now().Unix())
		}
		time.Sleep(time.Millisecond * 100) // to simulate real request scenario
	}
}

func LeakyBucketRateLimiterExample() {
	rate := 1
	maxBucketCapacity := 10
	refillRate := time.Second

	limiter := leaky_bucket.NewLeakyBucketRateLimiter(rate, maxBucketCapacity, refillRate)

	go func() {
		for true {
			req := limiter.ConsumeRequest()
			fmt.Println("Request received ", req)
		}
	}()
	// Simulate event messages
	for i := 0; i < 150; i++ {
		limiter.AddRequest(i)
	}
}

func FixedWindowRateLimiterExample() {
	rate := 10
	windowDuration := time.Second

	limiter := fixed_window.NewFixedWindowRateLimiter(rate, windowDuration)

	//Simulate request messages
	for i := 0; i < 150; i++ {
		if limiter.AllowRequest() {
			fmt.Println("Request allowed ", i, time.Now().Unix())
		} else {
			fmt.Println("Request discarded", i, time.Now().Unix())
		}
		time.Sleep(time.Millisecond * 200)
	}

	// Simulate event messages
	for i := 0; i < 150; i++ {
		limiter.AllowMessage()
		fmt.Println("Request allowed ", i, time.Now().Unix())
	}
}

func SlidingWindowLogRateLimiter() {
	windowDuration := time.Second
	maxRequest := 10

	limiter := sliding_window_log.NewSlidingWindowLogRateLimiter(windowDuration, maxRequest)
	//Simulate request messages
	for i := 0; i < 150; i++ {
		if limiter.AllowRequest() {
			fmt.Println("Request allowed ", i, time.Now().Unix())
		} else {
			fmt.Println("Request discarded", i, time.Now().Unix())
		}
		time.Sleep(time.Millisecond * 10)
	}

}

func main() {
	//TokenBucketRateLimiterExample()
	//LeakyBucketRateLimiterExample()
	//FixedWindowRateLimiterExample()
	SlidingWindowLogRateLimiter()
}
