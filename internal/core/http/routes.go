package routes

import (
	corerepository "github.com/bakare-dev/simple-bank-api/internal/core/repository"
	"github.com/bakare-dev/simple-bank-api/internal/core/service"
	userrepository "github.com/bakare-dev/simple-bank-api/internal/user/repository"
	"github.com/bakare-dev/simple-bank-api/pkg/mailer"
	mailerService "github.com/bakare-dev/simple-bank-api/pkg/mailer/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCoreRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepo := userrepository.NewUserRepository(db)
	accountRepo := corerepository.NewAccountRepository(db)
	transactionRepo := corerepository.NewTransactionRepository(db)

	mailer := mailer.NewMailer()

	notificationService := mailerService.NewNotificationService(*mailer)

	coreService := service.NewCoreService(*userRepo, *accountRepo, *transactionRepo, *notificationService)
	coreHandler := NewCoreHandler(coreService)

	coreRoutes := router.Group("/core")
	{
		coreRoutes.POST("/account", coreHandler.HandleCreateUser)
	}
}
