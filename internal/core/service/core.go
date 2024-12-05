package service

import (
	"context"
	"net/http"

	"github.com/bakare-dev/simple-bank-api/internal/core/model"
	corerepository "github.com/bakare-dev/simple-bank-api/internal/core/repository"
	userrepository "github.com/bakare-dev/simple-bank-api/internal/user/repository"
	mailerService "github.com/bakare-dev/simple-bank-api/pkg/mailer/service"
	"github.com/bakare-dev/simple-bank-api/pkg/util"
)

type CoreService struct {
	userRepo        userrepository.UserRepository
	profileRepo     userrepository.ProfileRepository
	accountRepo     corerepository.AccountRepository
	transactionRepo corerepository.TransactionRepository
	notificationSvc mailerService.NotificationService
}

func NewCoreService(userRepo userrepository.UserRepository, profileRepo userrepository.ProfileRepository, accountRepo corerepository.AccountRepository, transactionRepo corerepository.TransactionRepository, notificationSvc mailerService.NotificationService) *CoreService {
	return &CoreService{
		userRepo:        userRepo,
		profileRepo:     profileRepo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		notificationSvc: notificationSvc,
	}
}

func (s *CoreService) CreateAccount(ctx context.Context, userId string, account *model.Account) (int, *string) {
	user, err := s.userRepo.FindByID(ctx, userId)

	if err != nil {
		return http.StatusNotFound, ptr("User not found")
	}

	accountExists, err := s.accountRepo.Get(ctx, map[string]interface{}{"user_id": user.ID})

	if err != nil {
		return http.StatusInternalServerError, ptr("Try again later")
	}

	if accountExists != nil {
		return http.StatusConflict, ptr("Account already exists")
	}

	account.UserID = user.ID

	account.Number = util.GenerateAccountNumber()

	for {
		accountExists, err := s.accountRepo.GetAccountByAccountNumber(ctx, account.Number)

		if err != nil {
			return http.StatusInternalServerError, ptr("Try again later")
		}

		if accountExists != nil {
			account.Number = util.GenerateAccountNumber()
		} else {
			break
		}
	}

	_, err = s.accountRepo.Create(ctx, account)

	if err != nil {
		return http.StatusInternalServerError, ptr("Try again later")
	}

	return http.StatusCreated, nil
}

func (s *CoreService) GetAccount(ctx context.Context, userId string) (interface{}, int, *string) {
	user, err := s.userRepo.FindByID(ctx, userId)

	if err != nil {
		return nil, http.StatusNotFound, ptr("User not found")
	}

	account, err := s.accountRepo.Get(ctx, map[string]interface{}{"user_id": user.ID})

	if err != nil {
		return nil, http.StatusInternalServerError, ptr("Try again later")
	}

	if account == nil {
		return nil, http.StatusNotFound, ptr("account not found")
	}

	accountData := map[string]interface{}{
		"accountNumber": account.Number,
	}

	return accountData, http.StatusCreated, nil
}

func (s *CoreService) GetAccountByAccountNumber(ctx context.Context, accountNumber string) (interface{}, int, *string) {
	account, err := s.accountRepo.GetAccountByAccountNumber(ctx, accountNumber)

	if err != nil {
		return nil, http.StatusInternalServerError, ptr("Try again later")
	}

	if account == nil {
		return nil, http.StatusNotFound, ptr("Account not found")
	}

	arg := map[string]any{
		"user_id": account.UserID,
	}

	profile, err := s.profileRepo.Get(ctx, arg)
	if err != nil {
		return nil, http.StatusInternalServerError, ptr("Try again later")
	}

	if profile == nil {
		return nil, http.StatusNotFound, ptr("Profile not found")
	}

	accountData := map[string]interface{}{
		"accountNumber": account.Number,
		"name":          profile.FirstName + " " + profile.LastName,
	}

	return accountData, http.StatusOK, nil
}

func (s *CoreService) GetTransactions(ctx context.Context, userId string) (interface{}, int, *string) {
	return nil, http.StatusNotImplemented, ptr("Not implemented")
}

func (s *CoreService) GetTransaction(ctx context.Context, userId string, transactionId string) (interface{}, int, *string) {
	return nil, http.StatusNotImplemented, ptr("Not implemented")
}

func (s *CoreService) GetAccountBalance(ctx context.Context, userId string, transaction *model.Transaction) (int, *string) {
	return http.StatusNotImplemented, ptr("Not implemented")
}

func (s *CoreService) Transfer(ctx context.Context, userId string, transaction *model.Transaction) (int, *string) {
	return http.StatusNotImplemented, ptr("Not implemented")
}

func ptr(msg string) *string {
	return &msg
}
