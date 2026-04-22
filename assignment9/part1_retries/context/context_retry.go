package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fmt.Println("=== Starting Context-Aware Retry ===")

	for attempt := 1; ; attempt++ {
		fmt.Printf("Attempt %d in progress...\n", attempt)

		
		if ctx.Err() != nil {
			fmt.Printf("Retry stopped: %v\n", ctx.Err())
			return
		}

		
		time.Sleep(500 * time.Millisecond)

	
		select {
		case <-time.After(800 * time.Millisecond):
			// продолжаем цикл
		case <-ctx.Done():
			fmt.Printf("Retry cancelled during wait: %v\n", ctx.Err())
			return
		}
	}
}