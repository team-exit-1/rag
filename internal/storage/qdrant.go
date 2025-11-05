package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"

	"repo-rag-server/internal/models"
)

// QdrantStore implements VectorStore using REST API
type QdrantStore struct {
	baseURL    string
	collection string
	client     *http.Client
}

// NewQdrantStore creates a new Qdrant vector store
func NewQdrantStore(baseURL string, collection string) (*QdrantStore, error) {
	return &QdrantStore{
		baseURL:    baseURL,
		collection: collection,
		client:     &http.Client{},
	}, nil
}

// CollectionExists checks if a collection exists in Qdrant
func (qs *QdrantStore) CollectionExists(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/collections", qs.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create collections list request: %w", err)
	}

	resp, err := qs.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute collections list request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("qdrant returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var listResp struct {
		Result struct {
			Collections []struct {
				Name string `json:"name"`
			} `json:"collections"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return false, fmt.Errorf("failed to decode collections list response: %w", err)
	}

	// Check if our collection exists
	for _, col := range listResp.Result.Collections {
		if col.Name == qs.collection {
			return true, nil
		}
	}

	return false, nil
}

// InitializeCollection creates the collection if it doesn't exist
func (qs *QdrantStore) InitializeCollection(ctx context.Context, vectorSize int) error {
	// First, check if collection already exists
	exists, err := qs.CollectionExists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	if exists {
		fmt.Printf("Collection '%s' already exists, skipping creation\n", qs.collection)
		return nil
	}

	// Prepare collection creation request
	createRequest := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}

	body, err := json.Marshal(createRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal collection creation request: %w", err)
	}

	// Make HTTP request to create collection
	url := fmt.Sprintf("%s/collections/%s", qs.baseURL, qs.collection)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create collection request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := qs.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute collection creation request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Printf("Successfully created collection '%s'\n", qs.collection)
	return nil
}

// SaveVector saves an embedding vector to Qdrant
func (qs *QdrantStore) SaveVector(ctx context.Context, conversationID string, vector []float32, metadata map[string]interface{}) error {
	pointID := hashConversationID(conversationID)

	// Prepare payload with metadata
	payload := make(map[string]interface{})
	payload["conversation_id"] = conversationID
	for key, value := range metadata {
		payload[key] = value
	}

	// Prepare the point
	point := map[string]interface{}{
		"id":      pointID,
		"vector":  vector,
		"payload": payload,
	}

	// Create request body
	requestBody := map[string]interface{}{
		"points": []map[string]interface{}{point},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/collections/%s/points?wait=true", qs.baseURL, qs.collection)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := qs.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// SearchVectors searches for similar vectors in Qdrant
func (qs *QdrantStore) SearchVectors(ctx context.Context, queryVector []float32, limit int) ([]models.ConversationSearchResult, error) {
	// Prepare search request
	searchRequest := map[string]interface{}{
		"vector":       queryVector,
		"limit":        limit,
		"with_payload": true,
	}

	body, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/collections/%s/points/search", qs.baseURL, qs.collection)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := qs.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qdrant returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var searchResp struct {
		Result []struct {
			ID      uint64                 `json:"id"`
			Score   float32                `json:"score"`
			Payload map[string]interface{} `json:"payload"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert results
	var searchResults []models.ConversationSearchResult
	for _, item := range searchResp.Result {
		conversationID := ""
		if val, ok := item.Payload["conversation_id"]; ok {
			if strVal, ok := val.(string); ok {
				conversationID = strVal
			}
		}

		if conversationID == "" {
			continue
		}

		result := models.ConversationSearchResult{
			Conversation: &models.Conversation{
				ID: conversationID,
			},

			Score: item.Score,
		}
		searchResults = append(searchResults, result)
	}

	return searchResults, nil
}

// DeleteVector deletes a vector from Qdrant
func (qs *QdrantStore) DeleteVector(ctx context.Context, conversationID string) error {
	pointID := hashConversationID(conversationID)

	// Prepare delete request
	deleteRequest := map[string]interface{}{
		"points": []uint64{pointID},
	}

	body, err := json.Marshal(deleteRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal delete request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/collections/%s/points/delete", qs.baseURL, qs.collection)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := qs.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// Close closes the Qdrant client connection
func (qs *QdrantStore) Close() error {
	// HTTP client doesn't need explicit closing in this case
	return nil
}

// hashConversationID converts a string ID to a uint64 hash using FNV-1a
func hashConversationID(id string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(id))
	hash := h.Sum64()
	// Ensure the hash is never 0 (Qdrant uses 0 as a special value)
	if hash == 0 {
		hash = 1
	}
	return hash
}
