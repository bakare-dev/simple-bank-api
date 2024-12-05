package model

import (
	"time"

	"github.com/bakare-dev/simple-bank-api/pkg/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string
type UserStatus string

const (
	UserRoleAdmin          UserRole   = "admin"
	UserRoleCustomer       UserRole   = "customer"
	UserStatusActivated    UserStatus = "activated"
	UserStatusDeactivated  UserStatus = "deactivated"
	UserStatusNotActivated UserStatus = "not_activated"
)

type User struct {
	ID       string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email    string `json:"email" gorm:"unique;not null;index:idx_user_email"`
	Password string `json:"password" gorm:"not null"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Role UserRole `json:"role" gorm:"not null;default:'customer'"`

	Status UserStatus `json:"status" gorm:"not null;default:'not_activated'"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if user.Role == "" {
		user.Role = UserRoleCustomer
	}
	if user.Status == "" {
		user.Status = UserStatusNotActivated
	}
	return nil
}
