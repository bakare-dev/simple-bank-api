package dto

import (
	"github.com/bakare-dev/simple-bank-api/internal/core/model"
)

type CreateAccountRequest struct {
	Pin string `json:"pin" binding:"required,min=4,max=4"`
}

func (dto *CreateAccountRequest) ToModel() (*model.Account, error) {

	return &model.Account{
		Pin: dto.Pin,
	}, nil
}

type GetAccountByAccountNumberRequest struct {
	Number string `form:"acctno" binding:"required,min=10,max=10"`
}

type GetAccountTransactionRequest struct {
	Page string `form:"page" binding:"required,min=1"`
	Size string `form:"size" binding:"required,min=1,max=50"`
}

type GetTransactionRequest struct {
	TransactionId string `form:"id" binding:"required"`
}

type TransferRequest struct {
	ToAccountId string  `json:"toAccountId" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Pin         string  `json:"pin" binding:"required,min=4,max=4"`
}
