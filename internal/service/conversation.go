package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"repo-rag-server/internal/models"
	"repo-rag-server/internal/storage"
)

// ConversationService handles conversation business logic
type ConversationService struct {
	conversationStore storage.ConversationStore
	vectorStore       storage.VectorStore
	embeddingProvider storage.EmbeddingProvider
}

// NewConversationService creates a new conversation service
func NewConversationService(
	conversationStore storage.ConversationStore,
	vectorStore storage.VectorStore,
	embeddingProvider storage.EmbeddingProvider,
) *ConversationService {
	return &ConversationService{
		conversationStore: conversationStore,
		vectorStore:       vectorStore,
		embeddingProvider: embeddingProvider,
	}
}

// SaveConversation saves a new conversation and its embedding
func (cs *ConversationService) SaveConversation(ctx context.Context, req *models.ConversationSaveRequest) (*models.ConversationResponse, error) {
	// Generate a new ID
	conversationID := uuid.New().String()

	// Create embedding from the question
	embedding, err := cs.embeddingProvider.Embed(ctx, req.Question)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	// Save conversation to PostgreSQL
	now := time.Now()
	conversation := &models.Conversation{
		ID:        conversationID,
		UserID:    req.UserID,
		Question:  req.Question,
		Answer:    req.Answer,
		Metadata:  req.Metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := cs.conversationStore.SaveConversation(ctx, conversation); err != nil {
		return nil, fmt.Errorf("failed to save conversation: %w", err)
	}

	// Save embedding to Qdrant
	metadata := map[string]interface{}{
		"user_id": req.UserID,
		"created_at": now.Unix(),
	}

	if err := cs.vectorStore.SaveVector(ctx, conversationID, embedding, metadata); err != nil {
		// Log error but continue - we've already saved to PostgreSQL
		fmt.Printf("warning: failed to save vector to qdrant: %v\n", err)
	}

	return &models.ConversationResponse{
		ID:        conversationID,
		UserID:    req.UserID,
		Question:  req.Question,
		Answer:    req.Answer,
		Metadata:  req.Metadata,
		CreatedAt: now,
	}, nil
}

// SearchConversations searches for similar conversations
func (cs *ConversationService) SearchConversations(ctx context.Context, req *models.ConversationSearchRequest) ([]*models.ConversationResponse, error) {
	// Set default limit
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// Create embedding from the query
	queryEmbedding, err := cs.embeddingProvider.Embed(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to create query embedding: %w", err)
	}

	// Search in Qdrant
	searchResults, err := cs.vectorStore.SearchVectors(ctx, queryEmbedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	if len(searchResults) == 0 {
		return []*models.ConversationResponse{}, nil
	}

	// Extract conversation IDs from search results
	conversationIDs := make([]string, 0, len(searchResults))
	scoreMap := make(map[string]float32)

	for _, result := range searchResults {
		// We need to get the conversation ID from the Qdrant search
		// For now, we'll need to modify the search results structure
		// This is a limitation that needs to be addressed
		scoreMap[""] = result.Score
	}

	// Get conversations from PostgreSQL
	conversations, err := cs.conversationStore.GetConversationsByIDs(ctx, conversationIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	// Convert to response format with scores
	var responses []*models.ConversationResponse
	for _, conv := range conversations {
		responses = append(responses, &models.ConversationResponse{
			ID:        conv.ID,
			UserID:    conv.UserID,
			Question:  conv.Question,
			Answer:    conv.Answer,
			Metadata:  conv.Metadata,
			Score:     scoreMap[conv.ID],
			CreatedAt: conv.CreatedAt,
		})
	}

	return responses, nil
}

// GetConversation retrieves a single conversation by ID
func (cs *ConversationService) GetConversation(ctx context.Context, id string) (*models.ConversationResponse, error) {
	conversation, err := cs.conversationStore.GetConversation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	if conversation == nil {
		return nil, nil
	}

	return &models.ConversationResponse{
		ID:        conversation.ID,
		UserID:    conversation.UserID,
		Question:  conversation.Question,
		Answer:    conversation.Answer,
		Metadata:  conversation.Metadata,
		CreatedAt: conversation.CreatedAt,
	}, nil
}