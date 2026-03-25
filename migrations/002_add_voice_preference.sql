-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS voice_preference VARCHAR(100) NOT NULL DEFAULT 'en_US-lessac-medium';

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS voice_preference;
