package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)


func doSomethingUnreliable() error {
	if rand.Intn(10) < 7 { 
		return errors.New("temporary failure")
	}
	return nil
}

func main() {
	const maxRetries = 5
	const delay = 1 * time.Second
	var err error

	fmt.Println("=== Starting Fixed Delay Strategy ===")

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = doSomethingUnreliable()
		if err == nil {
			fmt.Printf("Attempt %d: Success!\n", attempt)
			return
		}

		fmt.Printf("Attempt %d failed, waiting %v before next retry...\n", attempt, delay)
		
		if attempt < maxRetries {
			time.Sleep(delay) 
		}
	}

	if err != nil {
		fmt.Printf("Operation failed after %d attempts: %v\n", maxRetries, err)
	}
}