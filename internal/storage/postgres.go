package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"refo-rag-server/internal/models"
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

	rows, err := ps.db.QueryContext(ctx, query, pq.Array(ids))
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

// GetDB returns the database connection
func (ps *PostgresStore) GetDB() *sql.DB {
	return ps.db
}

// Ping checks the database connection
func (ps *PostgresStore) Ping(ctx context.Context) error {
	return ps.db.PingContext(ctx)
}

// SavePersonalInfo saves a new personal information entry to PostgreSQL
func (ps *PostgresStore) SavePersonalInfo(ctx context.Context, personalInfo *models.PersonalInfo) error {
	query := `
		INSERT INTO personal_info (id, user_id, content, category, importance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			category = EXCLUDED.category,
			importance = EXCLUDED.importance,
			updated_at = EXCLUDED.updated_at
	`

	_, err := ps.db.ExecContext(
		ctx,
		query,
		personalInfo.ID,
		personalInfo.UserID,
		personalInfo.Content,
		personalInfo.Category,
		personalInfo.Importance,
		personalInfo.CreatedAt,
		personalInfo.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save personal info: %w", err)
	}

	return nil
}

// GetPersonalInfo retrieves a personal information entry by ID from PostgreSQL
func (ps *PostgresStore) GetPersonalInfo(ctx context.Context, id string) (*models.PersonalInfo, error) {
	query := `
		SELECT id, user_id, content, category, importance, created_at, updated_at
		FROM personal_info
		WHERE id = $1
	`

	personalInfo := &models.PersonalInfo{}
	err := ps.db.QueryRowContext(ctx, query, id).Scan(
		&personalInfo.ID,
		&personalInfo.UserID,
		&personalInfo.Content,
		&personalInfo.Category,
		&personalInfo.Importance,
		&personalInfo.CreatedAt,
		&personalInfo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get personal info: %w", err)
	}

	return personalInfo, nil
}

// GetPersonalInfoByUser retrieves all personal information entries for a user from PostgreSQL
func (ps *PostgresStore) GetPersonalInfoByUser(ctx context.Context, userID string) ([]*models.PersonalInfo, error) {
	query := `
		SELECT id, user_id, content, category, importance, created_at, updated_at
		FROM personal_info
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := ps.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query personal info: %w", err)
	}
	defer rows.Close()

	var personalInfoList []*models.PersonalInfo
	for rows.Next() {
		personalInfo := &models.PersonalInfo{}
		err := rows.Scan(
			&personalInfo.ID,
			&personalInfo.UserID,
			&personalInfo.Content,
			&personalInfo.Category,
			&personalInfo.Importance,
			&personalInfo.CreatedAt,
			&personalInfo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan personal info: %w", err)
		}
		personalInfoList = append(personalInfoList, personalInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating personal info: %w", err)
	}

	return personalInfoList, nil
}

// UpdatePersonalInfo updates an existing personal information entry in PostgreSQL
func (ps *PostgresStore) UpdatePersonalInfo(ctx context.Context, personalInfo *models.PersonalInfo) error {
	query := `
		UPDATE personal_info
		SET content = $1, category = $2, importance = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := ps.db.ExecContext(
		ctx,
		query,
		personalInfo.Content,
		personalInfo.Category,
		personalInfo.Importance,
		personalInfo.UpdatedAt,
		personalInfo.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update personal info: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("personal info not found")
	}

	return nil
}

// DeletePersonalInfo deletes a personal information entry from PostgreSQL
func (ps *PostgresStore) DeletePersonalInfo(ctx context.Context, id string) error {
	query := `DELETE FROM personal_info WHERE id = $1`

	result, err := ps.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete personal info: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("personal info not found")
	}

	return nil
}
