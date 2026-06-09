package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitorBucket struct {
	count     int
	resetTime time.Time
}

func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := map[string]*visitorBucket{}

	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()

		for range ticker.C {
			mu.Lock()
			now := time.Now()
			for key, bucket := range visitors {
				if now.After(bucket.resetTime) {
					delete(visitors, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		key := c.ClientIP()
		now := time.Now()

		mu.Lock()
		bucket, ok := visitors[key]
		if !ok || now.After(bucket.resetTime) {
			bucket = &visitorBucket{resetTime: now.Add(window)}
			visitors[key] = bucket
		}

		bucket.count++
		remaining := limit - bucket.count
		resetTime := bucket.resetTime
		limited := bucket.count > limit
		mu.Unlock()

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if limited {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests, please try again later"})
			return
		}

		c.Next()
	}
}
