package main

import (
	"fmt"
	"math"
	"time"
)

func doSomethingUnreliable() error {
	return fmt.Errorf("transient server error")
}

func main() {
	const maxRetries = 5
	const baseDelay = 100 * time.Millisecond
	const maxDelay = 5 * time.Second
	var err error

	fmt.Println("=== Starting Exponential Backoff Strategy ===")

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = doSomethingUnreliable()
		if err == nil {
			fmt.Println("Success!")
			return
		}

		
		backoffTime := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
		
		
		if backoffTime > maxDelay {
			backoffTime = maxDelay
		}

		fmt.Printf("Attempt %d failed, waiting %v...\n", attempt+1, backoffTime)
		
		if attempt < maxRetries-1 {
			time.Sleep(backoffTime)
		}
	}
}