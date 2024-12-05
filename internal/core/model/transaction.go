package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	Credit TransactionType = "credit"
	Debit  TransactionType = "debit"
)

type TransactionStatus string

const (
	Processing TransactionStatus = "processing"
	Successful TransactionStatus = "successful"
	Failed     TransactionStatus = "failed"
)

type Transaction struct {
	ID        string            `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	AccountID string            `json:"account_id" gorm:"type:uuid;not null"`
	Amount    float64           `json:"amount" gorm:"type:decimal(20,2);not null"`
	Type      TransactionType   `json:"type" gorm:"type:enum('credit', 'debit');not null"`
	Status    TransactionStatus `json:"status" gorm:"type:enum('processing', 'successful', 'failed');not null"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Account Account `json:"account" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AccountID"`
}

func (t *Transaction) BeforeCreate() error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}
