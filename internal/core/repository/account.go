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
