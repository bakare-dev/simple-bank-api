package model

import (
	"time"

	user "github.com/bakare-dev/simple-bank-api/internal/user/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	ID        string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	Number    string    `json:"number" gorm:"type:varchar(15);not null;uniqueIndex"`
	Pin       string    `json:"pin" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	User      user.User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
}

func (account *Account) BeforeCreate(tx *gorm.DB) error {
	if account.ID == "" {
		account.ID = uuid.New().String()
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(account.Pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Pin = string(hashedPin)

	return nil
}
