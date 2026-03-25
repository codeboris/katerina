package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/codeboris/katerina/internal/core/domain/model"
)

type AudioRepository struct {
	db *sql.DB
}

func NewAudioRepository(db *sql.DB) *AudioRepository {
	return &AudioRepository{db: db}
}

func (r *AudioRepository) Save(ctx context.Context, a *model.Audio) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO audio_files (id, user_id, file_path, mime_type, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		a.ID, a.UserID, a.FilePath, a.MimeType, a.CreatedAt)
	return err
}

func (r *AudioRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Audio, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, file_path, mime_type, created_at FROM audio_files WHERE id = $1`, id)

	var a model.Audio
	if err := row.Scan(&a.ID, &a.UserID, &a.FilePath, &a.MimeType, &a.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("audio not found")
		}
		return nil, err
	}
	return &a, nil
}
