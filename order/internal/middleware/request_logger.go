package middleware

import (
	"log"
	"net/http"
	"time"
)

// RequestLogger creates middleware for logging request execution time
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time of request processing
		startTime := time.Now()

		// Log the start of the request
		log.Printf("⏱️ Request started: %s %s", r.Method, r.URL.Path)

		// Pass control to the next handler
		next.ServeHTTP(w, r)

		// Calculate request execution time
		duration := time.Since(startTime)

		// Log request completion with execution time
		log.Printf("✅ Request completed: %s %s, execution time: %v", r.Method, r.URL.Path, duration)
	})
}
