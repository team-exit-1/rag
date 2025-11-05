package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"repo-rag-server/internal/models"
)

// PostgresStore implements ConversationStore
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &PostgresStore{db: db}, nil
}

// SaveConversation saves a new conversation to PostgreSQL
func (ps *PostgresStore) SaveConversation(ctx context.Context, conv *models.Conversation) error {
	query := `
		INSERT INTO conversations (id, user_id, question, answer, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			answer = EXCLUDED.answer,
			metadata = EXCLUDED.metadata,
			updated_at = EXCLUDED.updated_at
	`

	_, err := ps.db.ExecContext(
		ctx,
		query,
		conv.ID,
		conv.UserID,
		conv.Question,
		conv.Answer,
		conv.Metadata,
		conv.CreatedAt,
		conv.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save conversation: %w", err)
	}

	return nil
}

// GetConversation retrieves a conversation by ID from PostgreSQL
func (ps *PostgresStore) GetConversation(ctx context.Context, id string) (*models.Conversation, error) {
	query := `
		SELECT id, user_id, question, answer, metadata, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`

	conv := &models.Conversation{}
	err := ps.db.QueryRowContext(ctx, query, id).Scan(
		&conv.ID,
		&conv.UserID,
		&conv.Question,
		&conv.Answer,
		&conv.Metadata,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return conv, nil
}

// GetConversationsByIDs retrieves multiple conversations by IDs from PostgreSQL
func (ps *PostgresStore) GetConversationsByIDs(ctx context.Context, ids []string) ([]*models.Conversation, error) {
	if len(ids) == 0 {
		return []*models.Conversation{}, nil
	}

	query := `
		SELECT id, user_id, question, answer, metadata, created_at, updated_at
		FROM conversations
		WHERE id = ANY($1)
		ORDER BY created_at DESC
	`

	rows, err := ps.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*models.Conversation
	for rows.Next() {
		conv := &models.Conversation{}
		err := rows.Scan(
			&conv.ID,
			&conv.UserID,
			&conv.Question,
			&conv.Answer,
			&conv.Metadata,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating conversations: %w", err)
	}

	return conversations, nil
}

// Close closes the database connection
func (ps *PostgresStore) Close() error {
	return ps.db.Close()
}