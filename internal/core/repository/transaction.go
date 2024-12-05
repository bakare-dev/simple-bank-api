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

func (repo *TransactionRepository) GetBalance(ctx context.Context, accountId string) (float64, error) {
	var result struct {
		CreditSum float64
		DebitSum  float64
	}

	err := repo.db.WithContext(ctx).
		Model(&model.Transaction{}).
		Select(
			"SUM(CASE WHEN type = ? THEN amount ELSE 0 END) AS credit_sum, "+
				"SUM(CASE WHEN type = ? THEN amount ELSE 0 END) AS debit_sum",
			model.Credit, model.Debit,
		).
		Where("account_id = ?", accountId).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	balance := result.CreditSum - result.DebitSum
	return balance, nil
}
