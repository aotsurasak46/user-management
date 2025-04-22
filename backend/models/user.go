package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name             string    `json:"name"`
	Email            string    `gorm:"unique;not null" json:"email"`
	Password         string    `json:"-"`
	Role             string    `json:"role" gorm:"default:user"`
	MessagesSent     []Message `gorm:"foreignKey:FromID"`
	MessagesReceived []Message `gorm:"foreignKey:ToID"`
}
