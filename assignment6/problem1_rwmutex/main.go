package main

import (
	"fmt"
	"sync"
)

func main() {
	m := make(map[string]int)
	var mu sync.RWMutex
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			mu.Lock()
			m["key"] = v
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	mu.RLock()
	value := m["key"]
	mu.RUnlock()

	fmt.Println("Final value from RWMutex map:", value)
}