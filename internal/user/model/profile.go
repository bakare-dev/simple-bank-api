package model

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID      uint      `json:"user_id" gorm:"not null;uniqueIndex"`
	FirstName   string    `json:"first_name" gorm:"size:255;not null"`
	LastName    string    `json:"last_name" gorm:"size:255;not null"`
	PhoneNumber string    `json:"phone_number" gorm:"size:20"`
	ProfilePic  string    `json:"profile_pic" gorm:"size:255"`
	DateOfBirth time.Time `json:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
