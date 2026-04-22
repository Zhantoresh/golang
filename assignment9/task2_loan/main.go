package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"net/http/httptest"
	"github.com/google/uuid"
)

type CachedResponse struct {
	StatusCode int
	Body       []byte
	Completed  bool
}

type MemoryStore struct {
	mu   sync.Mutex
	data map[string]*CachedResponse
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]*CachedResponse)}
}


func IdempotencyMiddleware(store *MemoryStore, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
			return
		}

		store.mu.Lock()
		cached, exists := store.data[key]
		if exists {
			if !cached.Completed {
				store.mu.Unlock()
				http.Error(w, "Duplicate request in progress", http.StatusConflict)
				return
			}
			store.mu.Unlock()
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(cached.StatusCode)
			w.Write(cached.Body)
			return
		}

		
		store.data[key] = &CachedResponse{Completed: false}
		store.mu.Unlock()

		
		next(w, r)
	}
}


var store = NewMemoryStore()

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	
	fmt.Println("Processing started...")
	
	time.Sleep(2 * time.Second)

	response := map[string]interface{}{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": uuid.New().String(),
	}
	body, _ := json.Marshal(response)

	
	store.mu.Lock()
	store.data[key] = &CachedResponse{
		StatusCode: http.StatusOK,
		Body:       body,
		Completed:  true,
	}
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func main() {
	handler := IdempotencyMiddleware(store, PaymentHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	fmt.Println("Simulating Double-Click Attack...")
	key := uuid.New().String()
	var wg sync.WaitGroup

	
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			req, _ := http.NewRequest("POST", server.URL, nil)
			req.Header.Set("Idempotency-Key", key)
			
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("Request %d error: %v\n", id, err)
				return
			}
			fmt.Printf("Request %d: Status %d\n", id, resp.StatusCode)
		}(i)
	}

	wg.Wait()
	
	
	fmt.Println("\nRequest after completion:")
	req, _ := http.NewRequest("POST", server.URL, nil)
	req.Header.Set("Idempotency-Key", key)
	resp, _ := http.DefaultClient.Do(req)
	fmt.Printf("Final Request: Status %d\n", resp.StatusCode)
}