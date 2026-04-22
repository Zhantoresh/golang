package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
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


func (m *MemoryStore) Get(key string) (*CachedResponse, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	resp, exists := m.data[key]
	return resp, exists
}


func (m *MemoryStore) StartProcessing(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.data[key]; exists {
		return false 
	}
	
	m.data[key] = &CachedResponse{Completed: false}
	return true
}


func (m *MemoryStore) Finish(key string, status int, body []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if resp, exists := m.data[key]; exists {
		resp.StatusCode = status
		resp.Body = body
		resp.Completed = true
	}
}


func IdempotencyMiddleware(store *MemoryStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		
		
		if key == "" {
			http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
			return
		}

		
		if cached, exists := store.Get(key); exists {
			if cached.Completed {
				
				fmt.Println("Returning cached response for key:", key)
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
			} else {
				
				http.Error(w, "Duplicate request in progress", http.StatusConflict)
			}
			return
		}

		
		if !store.StartProcessing(key) {
			http.Error(w, "Conflict during processing", http.StatusConflict)
			return
		}

		
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		
		store.Finish(key, recorder.Code, recorder.Body.Bytes())
		
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	})
}

func main() {
	store := NewMemoryStore()
	
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing business logic...")
		time.Sleep(100 * time.Millisecond) 
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status": "created"}`))
	})

	protectedHandler := IdempotencyMiddleware(store, handler)

	
	server := httptest.NewServer(protectedHandler)
	defer server.Close()

	fmt.Println("=== Part 2: Idempotency Example ===")
	
	testKey := "unique-request-123"

	
	req1, _ := http.NewRequest("POST", server.URL, nil)
	req1.Header.Set("Idempotency-Key", testKey)
	resp1, _ := http.DefaultClient.Do(req1)
	fmt.Println("First request status:", resp1.StatusCode)

	
	req2, _ := http.NewRequest("POST", server.URL, nil)
	req2.Header.Set("Idempotency-Key", testKey)
	resp2, _ := http.DefaultClient.Do(req2)
	fmt.Println("Second request status:", resp2.StatusCode)
}