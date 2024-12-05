package repository

import (
	"context"

	"github.com/bakare-dev/simple-bank-api/internal/user/model"
	repository "github.com/bakare-dev/simple-bank-api/pkg/db"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
	repository.Repository[model.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db:         db,
		Repository: *repository.NewRepository[model.User](db),
	}
}

func (repo *UserRepository) PartialUpdate(ctx context.Context, id string, updatedData map[string]interface{}) error {
	if err := repo.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Select("password").
		Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}
