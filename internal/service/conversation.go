package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"refo-rag-server/internal/models"
	"refo-rag-server/internal/storage"
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
func (cs *ConversationService) SaveConversation(ctx context.Context, req *models.ConversationSaveRequest) (*models.SaveResponse, error) {
	// Use provided conversation ID or generate a new one
	conversationID := req.ConversationID
	if conversationID == "" {
		conversationID = uuid.New().String()
	}

	// Combine messages into a single text for embedding
	var textToEmbed string
	for _, msg := range req.Messages {
		textToEmbed += msg.Content + " "
	}

	// Create embedding from the combined messages
	embedding, err := cs.embeddingProvider.Embed(ctx, textToEmbed)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	// Save conversation to PostgreSQL
	now := time.Now()

	metadataStr := "{}"
	if req.Metadata != nil {
		metadataBytes, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataStr = string(metadataBytes)
	}

	conversation := &models.Conversation{
		ID:        conversationID,
		Question:  textToEmbed,
		Metadata:  metadataStr,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := cs.conversationStore.SaveConversation(ctx, conversation); err != nil {
		return nil, fmt.Errorf("failed to save conversation: %w", err)
	}

	// Save embedding to Qdrant
	metadata := map[string]interface{}{
		"created_at": now.Unix(),
	}

	if err := cs.vectorStore.SaveVector(ctx, conversationID, embedding, metadata); err != nil {
		// Log error but continue - we've already saved to PostgreSQL
		fmt.Printf("warning: failed to save vector to qdrant: %v\n", err)
	}

	return &models.SaveResponse{
		ConversationID:   conversationID,
		VectorsCreated:   1,
		MessagesStored:   len(req.Messages),
		StoredAt:         now.UTC().Format(time.RFC3339),
		ProcessingTimeMs: 0, // Will be set by handler
	}, nil
}

// SearchConversations searches for similar conversations
func (cs *ConversationService) SearchConversations(ctx context.Context, req *models.ConversationSearchRequest) ([]models.ConversationSearchResult, error) {
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
		return []models.ConversationSearchResult{}, nil
	}

	// Extract conversation IDs from search results
	conversationIDs := make([]string, 0, len(searchResults))
	scoreMap := make(map[string]float32)

	for _, result := range searchResults {
		conversationIDs = append(conversationIDs, result.ConversationID)
		scoreMap[result.ConversationID] = result.Score
	}

	// Get conversations from PostgreSQL
	conversations, err := cs.conversationStore.GetConversationsByIDs(ctx, conversationIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	// Convert to response format with scores and messages
	var responses []models.ConversationSearchResult
	for _, conv := range conversations {
		// Create message array from question and answer
		messages := []models.Message{}
		if conv.Question != "" {
			messages = append(messages, models.Message{
				Role:    "user",
				Content: conv.Question,
			})
		}
		if conv.Answer != "" {
			messages = append(messages, models.Message{
				Role:    "assistant",
				Content: conv.Answer,
			})
		}

		responses = append(responses, models.ConversationSearchResult{
			ConversationID: conv.ID,
			Score:          scoreMap[conv.ID],
			Timestamp:      conv.CreatedAt,
			Messages:       messages,
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
