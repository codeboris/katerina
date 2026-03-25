package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/codeboris/katerina/internal/core/domain/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Save(ctx context.Context, user *model.User) error
	UpdateVoicePreference(ctx context.Context, userID uuid.UUID, voice string) error
}

type ConversationRepository interface {
	Save(ctx context.Context, conv *model.Conversation) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Conversation, error)
}

type AudioRepository interface {
	Save(ctx context.Context, audio *model.Audio) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Audio, error)
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, token string) (uuid.UUID, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	BlacklistToken(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}
