package dto

import (
	"github.com/bakare-dev/simple-bank-api/internal/user/model"
)

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

type CreateProfileRequest struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	PhoneNumber string `json:"phoneNumber"`
	DOB         string `json:"dob" binding:"required"`
}

func (dto *CreateProfileRequest) ToModel() *model.Profile {
	return &model.Profile{
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		PhoneNumber: dto.PhoneNumber,
		DateOfBirth: dto.DOB,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (dto *LoginRequest) ToModel() *model.User {
	return &model.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type ActivateUserRequest struct {
	ID  string `json:"userId" binding:"required"`
	Otp string `json:"otp" binding:"required,min=6,max=6"`
}

type ResendOtpRequest struct {
	ID string `json:"userId" binding:"required"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}
