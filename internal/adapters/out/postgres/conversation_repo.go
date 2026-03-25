package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/codeboris/katerina/internal/core/domain/model"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Save(ctx context.Context, c *model.Conversation) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO conversations (id, user_id, original, corrected, translation, answer, audio_url, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.UserID, c.Original, c.Corrected, c.Translation, c.Answer, c.AudioURL, c.CreatedAt)
	return err
}

func (r *ConversationRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Conversation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, original, corrected, translation, answer, audio_url, created_at
		 FROM conversations WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Conversation
	for rows.Next() {
		var c model.Conversation
		if err := rows.Scan(&c.ID, &c.UserID, &c.Original, &c.Corrected, &c.Translation, &c.Answer, &c.AudioURL, &c.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &c)
	}
	return result, rows.Err()
}
