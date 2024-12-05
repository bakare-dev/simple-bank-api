package routes

import (
	userrepository "github.com/bakare-dev/simple-bank-api/internal/user/repository"
	"github.com/bakare-dev/simple-bank-api/internal/user/service"
	"github.com/bakare-dev/simple-bank-api/pkg/mailer"
	mailerService "github.com/bakare-dev/simple-bank-api/pkg/mailer/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepo := userrepository.NewUserRepository(db)
	profileRepo := userrepository.NewProfileRepository(db)
	mailer := mailer.NewMailer()
	notificationService := mailerService.NewNotificationService(*mailer)
	userService := service.NewUserService(*userRepo, *profileRepo, *notificationService)
	userHandler := NewUserHandler(userService)

	userRoutes := router.Group("/auth")
	{
		userRoutes.POST("/signup", userHandler.HandleCreateUser)
		userRoutes.POST("/profile", userHandler.HandleCreateProfile)
		userRoutes.POST("/signin", userHandler.HandleLogin)
		userRoutes.GET("/profile", userHandler.HandleGetProfile)
		userRoutes.PUT("/activate", userHandler.HandleActivateUser)
		userRoutes.PUT("/profile", userHandler.HandleUpdateProfile)
		userRoutes.POST("/resend-otp", userHandler.HandleResendOtp)
		userRoutes.PUT("/deactivate", userHandler.HandleDeactivateUser)
		userRoutes.POST("/change-password", userHandler.HandleChangePassword)
	}
}
