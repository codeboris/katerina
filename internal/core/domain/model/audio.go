package model

import (
	"time"

	"github.com/google/uuid"
)

type Audio struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FilePath  string
	MimeType  string
	CreatedAt time.Time
}
