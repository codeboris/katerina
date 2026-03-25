package model

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Original    string
	Corrected   string
	Translation string
	Answer      string
	AudioURL    string
	CreatedAt   time.Time
}
