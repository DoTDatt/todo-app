package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type window struct {
	requests []time.Time
}

var (
	store = make(map[string]*window)
	mu    sync.Mutex
)

const limit = 10
const windowSize = 10 * time.Second

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		defer mu.Unlock()

		w, exists := store[ip]

		if !exists {
			store[ip] = &window{
				requests: []time.Time{now},
			}
			c.Next()
			return
		}

		valid := []time.Time{}
		for _, t := range w.requests {
			if now.Sub(t) <= windowSize {
				valid = append(valid, t)
			}
		}
		w.requests = valid

		if len(w.requests) >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		w.requests = append(w.requests, now)

		c.Next()

	}
}

// type bucket struct {
// 	tokens     float64
// 	lastUpdate time.Time
// }

// var (
// 	store = make(map[string]*bucket)
// 	mu    sync.Mutex
// )

// const (
// 	capacity   = 10
// 	refillRate = 10.0 / 60.0
// )

// func RateLimit() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		ip := c.ClientIP()
// 		now := time.Now()

// 		mu.Lock()
// 		defer mu.Unlock()

// 		b, exists := store[ip]

// 		if !exists {
// 			store[ip] = &bucket{
// 				tokens:     capacity - 1,
// 				lastUpdate: now,
// 			}
// 			c.Next()
// 			return
// 		}

// 		elapsed := now.Sub(b.lastUpdate).Seconds()

// 		b.tokens += elapsed * refillRate

// 		if b.tokens > capacity {
// 			b.tokens = capacity
// 		}

// 		b.lastUpdate = now

// 		if b.tokens < 1 {
// 			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
// 				"error": "rate limit exceeded",
// 			})
// 			return
// 		}

// 		b.tokens--

// 		c.Next()
// 	}
// }
