package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) 
	
	const maxRetries = 5
	baseDelay := 100 * time.Millisecond
	maxDelay := 5 * time.Second

	fmt.Println("=== Starting Exponential Backoff + Full Jitter ===")

	for attempt := 0; attempt < maxRetries; attempt++ {
		
		fmt.Printf("Executing attempt %d...\n", attempt+1)

		
		backoff := float64(baseDelay) * math.Pow(2, float64(attempt))
		if backoff > float64(maxDelay) {
			backoff = float64(maxDelay)
		}

		
		sleepTime := time.Duration(rand.Int63n(int64(backoff)))

		fmt.Printf("Attempt %d failed, waiting %v (jittered)...\n", attempt+1, sleepTime)
		
		if attempt < maxRetries-1 {
			time.Sleep(sleepTime)
		}
	}
}