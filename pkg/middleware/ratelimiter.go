package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ClientRateLimiter struct {
	clients map[string]*rate.Limiter
	mu      sync.Mutex
	rps     int
	burst   int
}

func NewClientRateLimiter(rps, burst int) *ClientRateLimiter {
	return &ClientRateLimiter{
		clients: make(map[string]*rate.Limiter),
		rps:     rps,
		burst:   burst,
	}
}

func (c *ClientRateLimiter) getLimiter(ip string) *rate.Limiter {
	c.mu.Lock()
	defer c.mu.Unlock()

	if limiter, exists := c.clients[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(c.rps), c.burst)
	c.clients[ip] = limiter

	go func() {
		time.Sleep(1 * time.Minute)
		c.mu.Lock()
		delete(c.clients, ip)
		c.mu.Unlock()
	}()

	return limiter
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (c *ClientRateLimiter) GinMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx.Request)
		limiter := c.getLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		ctx.Next()
	}
}
