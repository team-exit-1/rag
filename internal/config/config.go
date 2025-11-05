package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port int
	Env  string // development, production

	// PostgreSQL
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string

	// Qdrant
	QdrantHost       string
	QdrantPort       int
	QdrantCollection string

	// OpenAI
	OpenAIAPIKey string
	OpenAIModel  string
	EmbeddingDim int

	// Logging
	LogLevel string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:             getEnvAsInt("PORT", 8080),
		Env:              getEnv("ENVIRONMENT", "development"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnvAsInt("POSTGRES_PORT", 5432),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       getEnv("POSTGRES_DB", "rag_db"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		QdrantHost:       getEnv("QDRANT_HOST", "localhost"),
		QdrantPort:       getEnvAsInt("QDRANT_PORT", 6334),
		QdrantCollection: getEnv("QDRANT_COLLECTION", "conversations"),
		OpenAIAPIKey:     getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:      getEnv("OPENAI_MODEL", "text-embedding-3-large"),
		EmbeddingDim:     getEnvAsInt("EMBEDDING_DIM", 3072),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

// GetPostgresDSN returns PostgreSQL connection string
func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
		c.PostgresSSLMode,
	)
}

// GetQdrantURL returns Qdrant server URL
func (c *Config) GetQdrantURL() string {
	return fmt.Sprintf("http://%s:%d", c.QdrantHost, c.QdrantPort)
}
