package models

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id        int64     `json:"id"`
	Body      string    `json:"body"`
	Complete  bool      `json:"complete" gorm:"default:false"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
	User_ID   uuid.UUID `json:"-" gorm:"not null"`
}
