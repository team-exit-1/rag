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

	// Create personal_info table
	createPersonalInfoTableSQL := `
	CREATE TABLE IF NOT EXISTS personal_info (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		category VARCHAR(50) NOT NULL,
		importance VARCHAR(20) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_personal_info_user_id ON personal_info(user_id);
	CREATE INDEX IF NOT EXISTS idx_personal_info_category ON personal_info(category);
	CREATE INDEX IF NOT EXISTS idx_personal_info_user_created ON personal_info(user_id, created_at DESC);

	CREATE OR REPLACE FUNCTION update_personal_info_updated_at()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS personal_info_updated_at_trigger ON personal_info;
	CREATE TRIGGER personal_info_updated_at_trigger
	BEFORE UPDATE ON personal_info
	FOR EACH ROW
	EXECUTE FUNCTION update_personal_info_updated_at();
	`

	_, err = db.ExecContext(ctx, createPersonalInfoTableSQL)
	if err != nil {
		return fmt.Errorf("failed to run personal_info migrations: %w", err)
	}

	return nil
}

// MigrateQdrant initializes Qdrant collection for vector storage
func MigrateQdrant(qdrantStore *QdrantStore, vectorSize int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := qdrantStore.InitializeCollection(ctx, vectorSize); err != nil {
		return fmt.Errorf("failed to initialize Qdrant collection: %w", err)
	}

	return nil
}
