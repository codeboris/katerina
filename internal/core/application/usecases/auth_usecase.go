package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/codeboris/katerina/internal/core/domain/model"
	"github.com/codeboris/katerina/internal/core/ports"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthConfig struct {
	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthUseCase struct {
	userRepo  ports.UserRepository
	tokenRepo ports.TokenRepository
	cfg       AuthConfig
}

func NewAuthUseCase(userRepo ports.UserRepository, tokenRepo ports.TokenRepository, cfg AuthConfig) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo, tokenRepo: tokenRepo, cfg: cfg}
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return uc.generateTokens(ctx, user.ID)
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*AuthOutput, error) {
	blacklisted, err := uc.tokenRepo.IsBlacklisted(ctx, refreshToken)
	if err != nil || blacklisted {
		return nil, ErrInvalidToken
	}
	userID, err := uc.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err := uc.tokenRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("delete token: %w", err)
	}
	if err := uc.tokenRepo.BlacklistToken(ctx, refreshToken, uc.cfg.RefreshTTL); err != nil {
		return nil, fmt.Errorf("blacklist token: %w", err)
	}
	return uc.generateTokens(ctx, userID)
}

func (uc *AuthUseCase) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(uc.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}
	return uuid.Parse(sub)
}

// SeedMockUser inserts the demo user if it does not yet exist.
func (uc *AuthUseCase) SeedMockUser(ctx context.Context) error {
	if _, err := uc.userRepo.FindByEmail(ctx, "demo@katerina.local"); err == nil {
		return nil // already seeded
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return uc.userRepo.Save(ctx, &model.User{
		ID:           uuid.New(),
		Email:        "demo@katerina.local",
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	})
}

func (uc *AuthUseCase) generateTokens(ctx context.Context, userID uuid.UUID) (*AuthOutput, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(uc.cfg.AccessTTL).Unix(),
		"iat": time.Now().Unix(),
	}).SignedString([]byte(uc.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshToken := uuid.New().String()
	if err := uc.tokenRepo.StoreRefreshToken(ctx, userID, refreshToken, uc.cfg.RefreshTTL); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &AuthOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
