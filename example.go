package main

import (
	"fmt"
	"github.com/yashwant-shrivastava/ratelimiter/token_bucket"
	"time"
)

func main() {
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
		if limiter.AllowMessage() {
			fmt.Println("Request", i+1, "allowed.", time.Now().Unix())
		}
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
