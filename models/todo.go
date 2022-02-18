package models

import (
	"github.com/google/uuid"
)

type Todo struct {
	Id       int64     `json:"id"`
	Body     string    `json:"body"`
	Complete bool      `json:"complete" gorm:"default:false"`
	User_ID  uuid.UUID `json:"-" gorm:"not null"`
}
