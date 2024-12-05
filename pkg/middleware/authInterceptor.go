package middleware

import (
	"net/http"
	"strings"

	"github.com/bakare-dev/simple-bank-api/pkg/response"
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
			response.Error(ctx, http.StatusUnauthorized, nil, "Authorization header is missing or invalid")
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := util.VerifyJWTToken(token)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, err, "Invalid or expired token")
			ctx.Abort()
			return
		}

		userID, userIDOk := claims["sub"].(string)
		role, roleOk := claims["role"].(string)
		if !userIDOk || !roleOk {
			response.Error(ctx, http.StatusUnauthorized, err, "Invalid or expired token")
			ctx.Abort()
			return
		}

		if ctx.Request.Method == http.MethodPost && ctx.Request.URL.Path == "/api/v1/auth/signout" {
			err := util.InvalidateToken(token)
			if err != nil {
				response.Error(ctx, http.StatusInternalServerError, err, "Failed to sign out")
				ctx.Abort()
				return
			}
			response.JSON(ctx, http.StatusOK, nil, "Sign out successful")
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("role", role)
		ctx.Next()
	}
}
