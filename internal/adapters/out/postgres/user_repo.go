package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/codeboris/katerina/internal/core/domain/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, voice_preference, created_at FROM users WHERE email = $1`, email)

	var u model.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.VoicePreference, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, voice_preference, created_at FROM users WHERE id = $1`, id)

	var u model.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.VoicePreference, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Save(ctx context.Context, u *model.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, password_hash, created_at) VALUES ($1, $2, $3, $4)`,
		u.ID, u.Email, u.PasswordHash, u.CreatedAt)
	return err
}

func (r *UserRepository) UpdateVoicePreference(ctx context.Context, userID uuid.UUID, voice string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET voice_preference = $1 WHERE id = $2`, voice, userID)
	return err
}
