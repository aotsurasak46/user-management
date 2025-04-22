package models

import (
	"time"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromID    uint      `json:"from_id"`
	From      User      `gorm:"foreignKey:FromID"`
	ToID      uint      `json:"to_id"`
	To        User      `gorm:"foreignKey:ToID"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime"`
}

