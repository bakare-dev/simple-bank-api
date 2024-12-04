package middleware

import (
	"strings"

	"github.com/bakare-dev/simple-bank-api/pkg/config"
	"github.com/gin-gonic/gin"
)

func RouteProtectionMiddleware(authInterceptor *AuthInterceptor) gin.HandlerFunc {
	unprotectedRoutes := config.Settings.Server.UnprotectedRoutes

	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		isUnprotected := false
		for _, route := range unprotectedRoutes {
			if strings.HasPrefix(path, route) {
				isUnprotected = true
				break
			}
		}

		if isUnprotected {
			ctx.Next()
			return
		}

		authInterceptor.Middleware()(ctx)
	}
}
