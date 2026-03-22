package main

import (
	"fmt"
	"sync"
)

/*
The final value is not always 1000 because counter++ is not an atomic operation,
so concurrent goroutines create a data race and some increments are lost.
*/

func main() {
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	fmt.Println("Broken counter:", counter)
}