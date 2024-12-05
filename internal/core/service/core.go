package service

import (
	"context"
	"net/http"
	"strconv"

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
		"id":            account.ID,
		"accountNumber": account.Number,
		"name":          profile.FirstName + " " + profile.LastName,
	}

	return accountData, http.StatusOK, nil
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
		"id":            account.ID,
		"accountNumber": account.Number,
		"name":          profile.FirstName + " " + profile.LastName,
	}

	return accountData, http.StatusOK, nil
}

func (s *CoreService) GetAccountBalance(ctx context.Context, userId string) (interface{}, int, *string) {
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

	balance, err := s.transactionRepo.GetBalance(ctx, account.ID)

	if err != nil {
		return nil, http.StatusInternalServerError, ptr("Try again later")
	}

	return map[string]float64{"balance": balance}, http.StatusOK, nil
}

func (s *CoreService) GetTransactions(ctx context.Context, userId string, page string, size string) (interface{}, int, *string) {
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
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, http.StatusBadRequest, ptr("Invalid page number")
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return nil, http.StatusBadRequest, ptr("Invalid size number")
	}

	transactions, err := s.transactionRepo.GetAll(ctx, pageInt, sizeInt, map[string]interface{}{"account_id": account.ID})

	if err != nil {
		return nil, http.StatusInternalServerError, ptr("Try again later")
	}

	if len(transactions) == 0 {
		return nil, http.StatusNotFound, ptr("No transactions found")
	}

	count, err := s.transactionRepo.Count(ctx, nil, map[string]interface{}{"account_id": account.ID})

	if err != nil {
		return nil, http.StatusInternalServerError, ptr("Try again later")

	}

	return map[string]interface{}{
		"transactions": transactions,
		"count":        count,
		"currentPage":  pageInt,
		"currentSize":  sizeInt,
	}, http.StatusOK, nil
}

func (s *CoreService) GetTransaction(ctx context.Context, userId string, transactionId string) (interface{}, int, *string) {
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

	transaction, err := s.transactionRepo.FindByID(ctx, transactionId)

	if err != nil {
		return nil, http.StatusInternalServerError, ptr("Try again later")
	}

	if transaction == nil {
		return nil, http.StatusNotFound, ptr("Transaction not found")
	}

	if transaction.AccountID != account.ID {
		return nil, http.StatusForbidden, ptr("Insufficient permissions")
	}

	return transaction, http.StatusOK, nil
}

func (s *CoreService) Transfer(ctx context.Context, userId string, toAccountId string, amount float64, pin string) (int, *string) {
	user, err := s.userRepo.FindByID(ctx, userId)

	if err != nil {
		return http.StatusNotFound, ptr("User not found")
	}

	account, err := s.accountRepo.Get(ctx, map[string]interface{}{"user_id": user.ID})

	if err != nil {
		return http.StatusInternalServerError, ptr("Try again later")
	}

	if account == nil {
		return http.StatusNotFound, ptr("sender account not found")
	}

	if err := util.CheckPassword(account.Pin, pin); err != nil {
		return http.StatusUnauthorized, ptr("invalid pin")
	}

	toAccount, err := s.accountRepo.FindByID(ctx, toAccountId)

	if err != nil {
		return http.StatusInternalServerError, ptr("Try again later")
	}

	if toAccount == nil {
		return http.StatusNotFound, ptr("receiver account not found")
	}

	balance, err := s.transactionRepo.GetBalance(ctx, account.ID)

	if err != nil {
		return http.StatusInternalServerError, ptr("Try again later")
	}

	if balance < amount {
		return http.StatusUnprocessableEntity, ptr("Insufficient balance")
	}

	debittransaction := &model.Transaction{
		AccountID: account.ID,
		Type:      model.Debit,
		Amount:    amount,
		Status:    model.Successful,
	}

	credittransaction := &model.Transaction{
		AccountID: toAccount.ID,
		Type:      model.Credit,
		Amount:    amount,
		Status:    model.Successful,
	}

	createdebittransaction, err := s.transactionRepo.Create(ctx, debittransaction)

	if err != nil {
		return http.StatusInternalServerError, ptr("Try again later")
	}

	_, err = s.transactionRepo.Create(ctx, credittransaction)

	if err != nil {
		s.transactionRepo.PartialUpdate(ctx, createdebittransaction.ID, map[string]interface{}{"status": model.Failed})
		return http.StatusInternalServerError, ptr("Try again later")
	}

	return http.StatusNotImplemented, ptr("Not implemented")
}

func ptr(msg string) *string {
	return &msg
}
