package repository

import (
	"github.com/bakare-dev/simple-bank-api/internal/user/model"
	repository "github.com/bakare-dev/simple-bank-api/pkg/db"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	repository.Repository[model.Profile]
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		Repository: *repository.NewRepository[model.Profile](db),
	}
}
