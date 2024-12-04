package middleware

import (
	"net/http"
	"strings"

	"github.com/bakare-dev/simple-bank-api/pkg/util"
	"github.com/gin-gonic/gin"
)

type AuthInterceptor struct{}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

func (a *AuthInterceptor) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing or invalid"})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := util.VerifyPASETOToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		userID, userIDOk := claims["sub"].(string)
		role, roleOk := claims["role"].(string)
		if !userIDOk || !roleOk {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token payload is invalid"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("role", role)
		ctx.Next()
	}
}
