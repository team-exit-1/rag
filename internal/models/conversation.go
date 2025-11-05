package models

import "time"

// Conversation represents a conversation record in the system
type Conversation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Metadata  string    `json:"metadata"` // JSON string for flexible metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConversationSearchRequest represents a request to search conversations
type ConversationSearchRequest struct {
	Query  string `json:"query"`
	UserID string `json:"user_id"`
	Limit  int    `json:"limit"`
}

// ConversationSearchResult represents a search result with similarity score
type ConversationSearchResult struct {
	ConversationID string    `json:"conversation_id"`
	Score          float32   `json:"score"`
	Timestamp      time.Time `json:"timestamp"`
	Messages       []Message `json:"messages"`
}

// Message represents a single message in a conversation
type Message struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

// ConversationSaveRequest represents a request to save a conversation
type ConversationSaveRequest struct {
	ConversationID string    `json:"conversation_id"`
	Messages       []Message `json:"messages"`
	Metadata       *Metadata `json:"metadata,omitempty"`
}

// Metadata represents conversation metadata
type Metadata struct {
	Source    string `json:"source,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

// ConversationResponse represents a conversation response
type ConversationResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Metadata  string    `json:"metadata"`
	Score     float32   `json:"score,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// APIResponse represents a standard API response wrapper
type APIResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	Error    *ErrorInfo  `json:"error,omitempty"`
	Metadata Metadata    `json:"metadata"`
}

// ErrorInfo represents error details in API response
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// SearchResponse represents the response for search API
type SearchResponse struct {
	Query          string                     `json:"query"`
	Results        []ConversationSearchResult `json:"results"`
	TotalResults   int                        `json:"total_results"`
	SearchMetadata SearchMetadata             `json:"search_metadata"`
}

// SearchMetadata represents search-specific metadata
type SearchMetadata struct {
	EmbeddingModel string `json:"embedding_model"`
	VectorDB       string `json:"vector_db"`
	SearchTimeMs   int64  `json:"search_time_ms"`
}

// SaveResponse represents the response for save API
type SaveResponse struct {
	ConversationID   string `json:"conversation_id"`
	VectorsCreated   int    `json:"vectors_created"`
	MessagesStored   int    `json:"messages_stored"`
	StoredAt         string `json:"stored_at"`
	ProcessingTimeMs int64  `json:"processing_time_ms"`
}

// HealthCheckResponse represents health check response
type HealthCheckResponse struct {
	Status       string             `json:"status"`
	Timestamp    string             `json:"timestamp"`
	Version      string             `json:"version"`
	Dependencies DependenciesStatus `json:"dependencies"`
}

// DependenciesStatus represents the status of dependencies
type DependenciesStatus struct {
	Qdrant     QdrantStatus     `json:"qdrant"`
	PostgreSQL PostgreSQLStatus `json:"postgresql"`
	OpenAI     OpenAIStatus     `json:"openai"`
}

// QdrantStatus represents Qdrant dependency status
type QdrantStatus struct {
	Status         string `json:"status"`
	ResponseTimeMs int    `json:"response_time_ms,omitempty"`
	Collections    int    `json:"collections,omitempty"`
	TotalVectors   int    `json:"total_vectors,omitempty"`
	Error          string `json:"error,omitempty"`
	LastSuccess    string `json:"last_success,omitempty"`
}

// PostgreSQLStatus represents PostgreSQL dependency status
type PostgreSQLStatus struct {
	Status         string          `json:"status"`
	ResponseTimeMs int             `json:"response_time_ms,omitempty"`
	Connections    ConnectionsInfo `json:"connections,omitempty"`
	Error          string          `json:"error,omitempty"`
}

// ConnectionsInfo represents database connections info
type ConnectionsInfo struct {
	Active int `json:"active"`
	Idle   int `json:"idle"`
	Max    int `json:"max"`
}

// OpenAIStatus represents OpenAI dependency status
type OpenAIStatus struct {
	Status    string `json:"status"`
	LastCheck string `json:"last_check,omitempty"`
	Error     string `json:"error,omitempty"`
}
