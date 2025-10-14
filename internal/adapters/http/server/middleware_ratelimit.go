package server

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// CustomRateLimiter holds the rate limiters for each client
type CustomRateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// visitor holds the rate limiter and last seen time for a client
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewCustomRateLimiter creates a new rate limiter
func NewCustomRateLimiter(r rate.Limit, b int) *CustomRateLimiter {
	rl := &CustomRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    b,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// cleanupVisitors removes old entries from the visitors map
func (rl *CustomRateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// getVisitor returns the rate limiter for the given IP
func (rl *CustomRateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// RateLimitConfig defines configuration for rate limiting
type RateLimitConfig struct {
	RequestsPerSecond float64
	Burst             int
	UseAPIKey         bool // If true, rate limit by API key; otherwise by IP
}

// rateLimitMiddleware implements rate limiting for the server
func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	// Default rate limit configuration
	// Can be made configurable via server config in the future
	config := RateLimitConfig{
		RequestsPerSecond: 10,
		Burst:             30,
		UseAPIKey:         false, // Default to IP-based rate limiting
	}

	// Use the consolidated rate limit middleware
	return CreateRateLimitMiddleware(config)(next)
}

// CreateRateLimitMiddleware creates configurable rate limiting middleware
func CreateRateLimitMiddleware(config RateLimitConfig) func(http.Handler) http.Handler {
	rl := NewCustomRateLimiter(rate.Limit(config.RequestsPerSecond), config.Burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Determine the rate limit key
			var key string
			if config.UseAPIKey {
				// Try to get API key from header
				key = r.Header.Get("X-API-Key")
				if key == "" {
					// Fall back to IP if no API key
					key = r.RemoteAddr
				}
			} else {
				// Use IP address for rate limiting
				key = r.RemoteAddr
			}

			limiter := rl.getVisitor(key)

			if !limiter.Allow() {
				response := NewTooManyRequestsResponse("rate limit exceeded", r.URL.Path)
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusTooManyRequests)
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
