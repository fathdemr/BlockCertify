package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateWindow struct {
	count   int
	resetAt time.Time
}

// RateLimit allows at most maxRequests per client IP within the given window.
// State is in-memory; suitable for a single-instance deployment.
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {

	var (
		mu      sync.Mutex
		clients = make(map[string]*rateWindow)
	)

	// Drop expired entries periodically so the map doesn't grow unbounded.
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			now := time.Now()
			for ip, w := range clients {
				if now.After(w.resetAt) {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		w, ok := clients[ip]
		if !ok || time.Now().After(w.resetAt) {
			w = &rateWindow{resetAt: time.Now().Add(window)}
			clients[ip] = w
		}
		w.count++
		exceeded := w.count > maxRequests
		mu.Unlock()

		if exceeded {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			return
		}

		c.Next()
	}
}
