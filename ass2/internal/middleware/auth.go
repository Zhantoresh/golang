package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

const APIKey = "secret12345"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger") {
			next.ServeHTTP(w, r)
			return
		}

		apiKey := r.Header.Get("X-API-KEY")
		
		if apiKey != APIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format("2006-01-02T15:04:05")
		log.Printf("%s %s %s {request received}", timestamp, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
