package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Id       uuid.UUID `gorm:"primaryKey"`
}
