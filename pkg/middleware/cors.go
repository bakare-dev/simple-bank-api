package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	}
}

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		ctx.Writer.Header().Set("X-Frame-Options", "DENY")
		ctx.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		ctx.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		log.Printf("[%s] %s %s %d %s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.ClientIP(),
			ctx.Writer.Status(),
			time.Since(startTime),
		)
	}
}
