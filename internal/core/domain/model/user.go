package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Email           string
	PasswordHash    string
	VoicePreference string
	CreatedAt       time.Time
}
