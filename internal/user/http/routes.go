package routes

import (
	userrepository "github.com/bakare-dev/simple-bank-api/internal/user/repository"
	"github.com/bakare-dev/simple-bank-api/internal/user/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepo := userrepository.NewUserRepository(db)

	userService := service.NewUserService(*userRepo)

	userRoutes := router.Group("/auth")
	{
		userRoutes.POST("/signup", func(ctx *gin.Context) {
			handleCreateUser(ctx, userService)
		})
	}
}
