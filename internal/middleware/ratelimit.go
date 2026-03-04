package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	tokens   float64
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
	rate     = 10.0              // requests per second
	capacity = 20.0              // max burst
	cleanup  = 5 * time.Minute   // cleanup interval
)

func RateLimiter() gin.HandlerFunc {
	go cleanupVisitors()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{
				lastSeen: time.Now(),
				tokens:   capacity,
			}
			v = visitors[ip]
		}

		// Refill tokens based on time elapsed
		now := time.Now()
		elapsed := now.Sub(v.lastSeen).Seconds()
		v.tokens += elapsed * rate
		if v.tokens > capacity {
			v.tokens = capacity
		}
		v.lastSeen = now

		// Check if request can proceed
		if v.tokens >= 1 {
			v.tokens--
			mu.Unlock()
			c.Next()
		} else {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
		}
	}
}

func cleanupVisitors() {
	for {
		time.Sleep(cleanup)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > cleanup {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
