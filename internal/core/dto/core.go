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
