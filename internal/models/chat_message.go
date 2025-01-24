package models

import (
	"time"
)

type ChatMessage struct {
	ID             int       `json:"id" db:"id"`
	ConsultationID int       `json:"consultation_id" db:"consultation_id"`
	SenderType     string    `json:"sender_type" db:"sender_type"`
	SenderID       int       `json:"sender_id" db:"sender_id"`
	Message        string    `json:"message" db:"message"`
	SentAt         time.Time `json:"sent_at" db:"sent_at"`
	IsRead         bool      `json:"is_read" db:"is_read"`
}