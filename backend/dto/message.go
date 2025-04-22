package dto

import (
	"time"
)

type MessageRequest struct {
    To      uint   `json:"to"`      
    Content string `json:"content"` 
	TempID 	string `json:"tempId"`
}


type MessageResponse struct {
	ID        uint          `json:"ID"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at,omitempty"`
	FromID    uint          `json:"from_id"`
    From      UserResponse  `json:"from"`
	ToID      uint          `json:"to_id"`
	To        UserResponse  `json:"to"`
	Content   string        `json:"content"`
	Timestamp time.Time     `json:"timestamp"`
}
