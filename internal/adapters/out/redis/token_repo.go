package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, ttl time.Duration) error {
	return r.client.Set(ctx, refreshKey(token), userID.String(), ttl).Err()
}

func (r *TokenRepository) GetRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	val, err := r.client.Get(ctx, refreshKey(token)).Result()
	if err != nil {
		return uuid.Nil, fmt.Errorf("token not found: %w", err)
	}
	return uuid.Parse(val)
}

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.client.Del(ctx, refreshKey(token)).Err()
}

func (r *TokenRepository) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	return r.client.Set(ctx, blacklistKey(token), "1", ttl).Err()
}

func (r *TokenRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	n, err := r.client.Exists(ctx, blacklistKey(token)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func refreshKey(token string) string   { return "refresh:" + token }
func blacklistKey(token string) string { return "blacklist:" + token }
