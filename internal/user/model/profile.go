package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID          string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID      string    `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	FirstName   string    `json:"first_name" gorm:"size:255;not null"`
	LastName    string    `json:"last_name" gorm:"size:255;not null"`
	PhoneNumber string    `json:"phone_number" gorm:"size:20"`
	DateOfBirth string    `json:"dob" gorm:"size:16;not null`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	User        User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
}

func (profile *Profile) BeforeCreate(tx *gorm.DB) error {
	if profile.ID == "" {
		profile.ID = uuid.New().String()
	}
	return nil
}
