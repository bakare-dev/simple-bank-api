package repository

import (
	"context"

	"github.com/bakare-dev/simple-bank-api/internal/user/model"
	repository "github.com/bakare-dev/simple-bank-api/pkg/db"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
	repository.Repository[model.Profile]
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		db:         db,
		Repository: *repository.NewRepository[model.Profile](db),
	}
}

func (repo *ProfileRepository) PartialUpdate(ctx context.Context, id string, updatedData map[string]interface{}) error {
	if err := repo.db.WithContext(ctx).Model(&model.Profile{}).Where("id = ?", id).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}
