package model

import (
	"errors"
	"fmt"
	"time"

	user "github.com/bakare-dev/simple-bank-api/internal/user/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountType string

const (
	Savings AccountType = "savings"
	Current AccountType = "current"
)

type Account struct {
	ID        string      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    string      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	Number    string      `json:"number" gorm:"type:varchar(15);not null;uniqueIndex"`
	Type      AccountType `json:"type" gorm:"type:enum('savings', 'current');not null"`
	Pin       string      `json:"pin" gorm:"not null"`
	CreatedAt time.Time   `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time   `json:"updatedAt" gorm:"autoUpdateTime"`
	User      user.User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
}

func (account *Account) BeforeCreate(tx *gorm.DB) error {
	if account.ID == "" {
		account.ID = uuid.New().String()
	}

	if account.Type != Savings && account.Type != Current {
		return errors.New("invalid account type, must be 'savings' or 'current'")
	}

	if account.Number == "" {
		account.Number = generateAccountNumber()
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(account.Pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Pin = string(hashedPin)

	return nil
}

func generateAccountNumber() string {
	return fmt.Sprintf("%010d", time.Now().UnixNano()%1e10)
}
