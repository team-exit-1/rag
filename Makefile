.PHONY: help build run test clean setup-db setup-qdrant deps

# Variables
BINARY_NAME=rag-server
GO=go
GOFLAGS=-v
OUTPUT_DIR=./bin

help:
	@echo "RAG Server - Available targets:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Download dependencies"
	@echo "  make setup-db       - Setup PostgreSQL database"
	@echo "  make setup-qdrant   - Setup Qdrant vector database"
	@echo "  make setup-all      - Setup both databases"

build: deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build $(GOFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/server

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(OUTPUT_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	$(GO) test -v ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(OUTPUT_DIR)
	$(GO) clean

deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

setup-db:
	@echo "Setting up PostgreSQL..."
	@bash scripts/setup_postgres.sh

setup-qdrant:
	@echo "Setting up Qdrant..."
	@bash scripts/setup_qdrant.sh

setup-all: setup-db setup-qdrant
	@echo "Database setup completed!"

docker-build:
	@echo "Building Docker image..."
	docker build -t repo-rag-server:latest .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file config/.env repo-rag-server:latest