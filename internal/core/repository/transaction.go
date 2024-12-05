package repository

import (
	"context"

	"github.com/bakare-dev/simple-bank-api/internal/core/model"
	repository "github.com/bakare-dev/simple-bank-api/pkg/db"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
	repository.Repository[model.Transaction]
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db:         db,
		Repository: *repository.NewRepository[model.Transaction](db),
	}
}

func (repo *TransactionRepository) PartialUpdate(ctx context.Context, id string, updatedData map[string]interface{}) error {
	if err := repo.db.WithContext(ctx).Model(&model.Transaction{}).Where("id = ?", id).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}
