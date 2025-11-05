package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"repo-rag-server/internal/api"
	"repo-rag-server/internal/config"
	"repo-rag-server/internal/service"
	"repo-rag-server/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize PostgreSQL connection
	postgresStore, err := storage.NewPostgresStore(cfg.GetPostgresDSN())
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	defer postgresStore.Close()

	// Run migrations
	log.Println("Running database migrations...")
	if err := storage.Migrate(postgresStore.GetDB()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Initialize Qdrant connection
	qdrantStore, err := storage.NewQdrantStore(cfg.GetQdrantURL(), cfg.QdrantCollection)
	if err != nil {
		log.Fatalf("Failed to initialize Qdrant: %v", err)
	}
	defer qdrantStore.Close()

	// Initialize OpenAI embedding provider
	embeddingProvider := storage.NewOpenAIEmbeddingProvider(
		cfg.OpenAIAPIKey,
		cfg.OpenAIModel,
		cfg.EmbeddingDim,
	)

	// Initialize services
	conversationService := service.NewConversationService(
		postgresStore,
		qdrantStore,
		embeddingProvider,
	)

	// Setup Gin router
	router := api.Router(conversationService)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting RAG server on %s", addr)

	// Run server in a goroutine
	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down RAG server...")
}
