package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Migrate creates all necessary tables
func Migrate(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create conversations table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS conversations (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		question TEXT NOT NULL,
		answer TEXT,
		metadata JSONB,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
	CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_conversations_user_created ON conversations(user_id, created_at DESC);

	CREATE OR REPLACE FUNCTION update_conversations_updated_at()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS conversations_updated_at_trigger ON conversations;
	CREATE TRIGGER conversations_updated_at_trigger
	BEFORE UPDATE ON conversations
	FOR EACH ROW
	EXECUTE FUNCTION update_conversations_updated_at();
	`

	_, err := db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
