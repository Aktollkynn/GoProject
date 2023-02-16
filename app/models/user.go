package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID            string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Addresses     []Address
	FirstName     string `gorm:"size:100;not null"`
	LastName      string `gorm:"size:100;not null"`
	Email         string `gorm:"size:100;not null"`
	Password      string `gorm:"size:255;not null"`
	RememberToken string `gorm:"size:255;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
