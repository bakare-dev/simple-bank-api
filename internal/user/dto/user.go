package dto

import "github.com/bakare-dev/simple-bank-api/internal/user/model"

type User struct {
	ID string `json:"id"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (dto *CreateUserRequest) ToModel() *model.User {
	return &model.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func UserResponseFromModel(user *model.User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}
