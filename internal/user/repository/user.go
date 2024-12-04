package repository

import (
	"github.com/bakare-dev/simple-bank-api/internal/user/model"
	repository "github.com/bakare-dev/simple-bank-api/pkg/db"
	"gorm.io/gorm"
)

type UserRepository struct {
	repository.Repository[model.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository: *repository.NewRepository[model.User](db),
	}
}
