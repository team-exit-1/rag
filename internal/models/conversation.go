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
	Conversation *Conversation `json:"conversation"`
	Score        float32       `json:"score"`
}

// ConversationSaveRequest represents a request to save a conversation
type ConversationSaveRequest struct {
	UserID   string `json:"user_id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Metadata string `json:"metadata"`
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