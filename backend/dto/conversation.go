package dto

import (
	"time"
)

type ConversationResponse struct {
	User        UserResponse `json:"user"`
	LastMessage string       `json:"last_message"`
	Timestamp   time.Time    `json:"timestamp"`
}
