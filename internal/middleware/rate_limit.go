package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

func RateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
	limiter := NewRateLimiter(
		rate.Limit(requestsPerSecond),
		requestsPerSecond,
	)

	return func(c *gin.Context) {
		// Usa IP ou API Key como identificador
		key := c.ClientIP()
		if userID := c.GetHeader("X-User-ID"); userID != "" {
			key = userID
		}

		if !limiter.GetLimiter(key).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests. Please try again later.",
				"retry_after": time.Now().Add(time.Second).Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
