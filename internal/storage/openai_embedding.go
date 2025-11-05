package storage

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// OpenAIEmbeddingProvider implements EmbeddingProvider using OpenAI API
type OpenAIEmbeddingProvider struct {
	client    *openai.Client
	model     openai.EmbeddingModel
	dimension int
}

// NewOpenAIEmbeddingProvider creates a new OpenAI embedding provider
func NewOpenAIEmbeddingProvider(apiKey string, model string, dimension int) *OpenAIEmbeddingProvider {
	client := openai.NewClient(apiKey)
	return &OpenAIEmbeddingProvider{
		client:    client,
		model:     openai.EmbeddingModel(model),
		dimension: dimension,
	}
}

// Embed converts text to a vector using OpenAI
func (oaep *OpenAIEmbeddingProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	resp, err := oaep.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: oaep.model,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned from openai")
	}

	return resp.Data[0].Embedding, nil
}

// EmbedBatch converts multiple texts to vectors using OpenAI
func (oaep *OpenAIEmbeddingProvider) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	resp, err := oaep.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: texts,
		Model: oaep.model,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create batch embeddings: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned from openai")
	}

	// Sort embeddings by index to ensure correct order
	embeddings := make([][]float32, len(texts))
	for _, data := range resp.Data {
		if data.Index < len(embeddings) {
			embeddings[data.Index] = data.Embedding
		}
	}

	return embeddings, nil
}