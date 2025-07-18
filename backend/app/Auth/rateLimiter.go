package auth

import (
	"net/http"
	"sync"
	"time"
)

var (
	mu       sync.Mutex
	clients  = make(map[string]*client)
	limit    = 15
	interval = 10 * time.Second
)

type client struct {
	requests []time.Time
}

func ApplyRateLimit(key string, w http.ResponseWriter, r *http.Request, next http.Handler) {
	mu.Lock()
	c, exists := clients[key]
	if !exists {
		c = &client{}
		clients[key] = c
	}
	now := time.Now()

	// Clean old requests
	valid := []time.Time{}
	for _, t := range c.requests {
		if now.Sub(t) < interval {
			valid = append(valid, t)
		}	
	}
	c.requests = valid

	if len(c.requests) >= limit {
		mu.Unlock()
		JsResponse(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	c.requests = append(c.requests, now)
	mu.Unlock()

	next.ServeHTTP(w, r)
}
