package repository

import (
	"context"

	"github.com/bakare-dev/simple-bank-api/internal/core/model"
	repository "github.com/bakare-dev/simple-bank-api/pkg/db"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
	repository.Repository[model.Account]
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db:         db,
		Repository: *repository.NewRepository[model.Account](db),
	}
}

func (repo *AccountRepository) PartialUpdate(ctx context.Context, id string, updatedData map[string]interface{}) error {
	if err := repo.db.WithContext(ctx).Model(&model.Account{}).Where("id = ?", id).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AccountRepository) GetAccountByAccountNumber(ctx context.Context, accountNumber string) (*model.Account, error) {
	var account model.Account
	if err := repo.db.WithContext(ctx).Where("number = ?", accountNumber).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}
