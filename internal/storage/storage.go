package storage

import (
	"context"

	"refo-rag-server/internal/models"
)

// ConversationStore defines the interface for storing conversations
type ConversationStore interface {
	// SaveConversation saves a new conversation to the database
	SaveConversation(ctx context.Context, conversation *models.Conversation) error

	// GetConversation retrieves a conversation by ID
	GetConversation(ctx context.Context, id string) (*models.Conversation, error)

	// GetConversationsByIDs retrieves multiple conversations by IDs
	GetConversationsByIDs(ctx context.Context, ids []string) ([]*models.Conversation, error)

	// Close closes the database connection
	Close() error
}

// PostgresStoreInterface defines the interface for PostgreSQL operations
type PostgresStoreInterface interface {
	ConversationStore
	PersonalInfoStore
	// Ping checks the database connection
	Ping(ctx context.Context) error
}

// QdrantStoreInterface defines the interface for Qdrant operations
type QdrantStoreInterface interface {
	VectorStore
	// CollectionExists checks if a collection exists
	CollectionExists(ctx context.Context) (bool, error)
}

// VectorStore defines the interface for storing and searching vectors
type VectorStore interface {
	// SaveVector saves an embedding vector with metadata
	SaveVector(ctx context.Context, conversationID string, vector []float32, metadata map[string]interface{}) error

	// SearchVectors searches for similar vectors
	SearchVectors(ctx context.Context, queryVector []float32, limit int) ([]models.ConversationSearchResult, error)

	// DeleteVector deletes a vector by conversation ID
	DeleteVector(ctx context.Context, conversationID string) error

	// Close closes the vector store connection
	Close() error
}

// EmbeddingProvider defines the interface for text embedding services
type EmbeddingProvider interface {
	// Embed converts text to a vector
	Embed(ctx context.Context, text string) ([]float32, error)

	// EmbedBatch converts multiple texts to vectors
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}

// PersonalInfoStore defines the interface for storing personal information
type PersonalInfoStore interface {
	// SavePersonalInfo saves personal information to the database
	SavePersonalInfo(ctx context.Context, personalInfo *models.PersonalInfo) error

	// GetPersonalInfo retrieves personal information by ID
	GetPersonalInfo(ctx context.Context, id string) (*models.PersonalInfo, error)

	// GetPersonalInfoByUser retrieves all personal information for a user
	GetPersonalInfoByUser(ctx context.Context, userID string) ([]*models.PersonalInfo, error)

	// UpdatePersonalInfo updates existing personal information
	UpdatePersonalInfo(ctx context.Context, personalInfo *models.PersonalInfo) error

	// DeletePersonalInfo deletes personal information by ID
	DeletePersonalInfo(ctx context.Context, id string) error

	// Close closes the store
	Close() error
}
