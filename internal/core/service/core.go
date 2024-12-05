package service

import (
	"context"

	corerepository "github.com/bakare-dev/simple-bank-api/internal/core/repository"
	userrepository "github.com/bakare-dev/simple-bank-api/internal/user/repository"
	mailerService "github.com/bakare-dev/simple-bank-api/pkg/mailer/service"
)

type CoreService struct {
	userRepo        userrepository.UserRepository
	accountRepo     corerepository.AccountRepository
	transactionRepo corerepository.TransactionRepository
	notificationSvc mailerService.NotificationService
}

func NewCoreService(userRepo userrepository.UserRepository, accountRepo corerepository.AccountRepository, transactionRepo corerepository.TransactionRepository, notificationSvc mailerService.NotificationService) *CoreService {
	return &CoreService{
		userRepo:        userRepo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		notificationSvc: notificationSvc,
	}
}

func (s *CoreService) CreateUser(ctx context.Context) {
	return
}
