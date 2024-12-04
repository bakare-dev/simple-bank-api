package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bakare-dev/simple-bank-api/internal/user/model"
	userrepository "github.com/bakare-dev/simple-bank-api/internal/user/repository"
)

type UserService struct {
	userRepo userrepository.UserRepository
}

func NewUserService(userRepo userrepository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, int, error) {
	userExists, err := s.userRepo.Get(ctx, map[string]any{"email": user.Email})

	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to check existing user")
	}

	if userExists != nil {
		return nil, http.StatusConflict, fmt.Errorf("user with this email already exists")
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to create user: %v", err)
	}
	createdUser.Password = ""

	return createdUser, http.StatusCreated, nil
}
