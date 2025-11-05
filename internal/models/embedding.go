package models

// EmbeddingVector represents an embedding vector stored in Qdrant
type EmbeddingVector struct {
	ConversationID string    `json:"conversation_id"`
	Vector         []float32 `json:"vector"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// EmbeddingRequest represents a request to create embeddings
type EmbeddingRequest struct {
	Text string `json:"text"`
}

// EmbeddingResponse represents OpenAI's embedding response
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingUsage represents token usage in embedding request
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}