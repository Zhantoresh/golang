package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)


func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	
	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}


func CalculateBackoff(attempt int, baseDelay time.Duration) time.Duration {
	
	backoff := float64(baseDelay) * math.Pow(2, float64(attempt))
	
	return time.Duration(rand.Int63n(int64(backoff)))
}

func ExecutePayment(ctx context.Context, url string, maxRetries int, baseDelay time.Duration) error {
	client := &http.Client{}

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(`{"amount": 1000}`))

		resp, err := client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Attempt %d: Success!\n", attempt+1)
			return nil
		}

		
		if !IsRetryable(resp, err) {
			return fmt.Errorf("non-retryable error or exhausted: %v", err)
		}

		if attempt == maxRetries-1 {
			return fmt.Errorf("max retries reached")
		}

		wait := CalculateBackoff(attempt, baseDelay)
		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt+1, wait)

		
		select {
		case <-time.After(wait):
			// продолжаем цикл
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func main() {
	
	count := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if count <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable) // 503
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": "success"}`)
	}))
	defer ts.Close()

	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := ExecutePayment(ctx, ts.URL, 5, 500*time.Millisecond)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
	}
}