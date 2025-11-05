#!/bin/bash

# setup_qdrant.sh - Initialize Qdrant vector database for RAG server

set -e

# Configuration
QDRANT_HOST=${QDRANT_HOST:-localhost}
QDRANT_PORT=${QDRANT_PORT:-6334}
QDRANT_COLLECTION=${QDRANT_COLLECTION:-conversations}
EMBEDDING_DIM=${EMBEDDING_DIM:-3072}

QDRANT_URL="http://$QDRANT_HOST:$QDRANT_PORT"

echo "Qdrant setup script"
echo "===================="
echo "Qdrant URL: $QDRANT_URL"
echo "Collection: $QDRANT_COLLECTION"
echo "Embedding dimension: $EMBEDDING_DIM"
echo ""

# Check if Qdrant is running
echo "Checking Qdrant connection..."
if ! curl -s "$QDRANT_URL/health" > /dev/null; then
    echo "Error: Cannot connect to Qdrant at $QDRANT_URL"
    echo "Make sure Qdrant is running on $QDRANT_HOST:$QDRANT_PORT"
    exit 1
fi

echo "✓ Connected to Qdrant"

# Check if collection exists
echo "Checking if collection '$QDRANT_COLLECTION' exists..."
COLLECTION_EXISTS=$(curl -s "$QDRANT_URL/collections/$QDRANT_COLLECTION" | grep -q '"name":"'$QDRANT_COLLECTION'"' && echo "true" || echo "false")

if [ "$COLLECTION_EXISTS" = "true" ]; then
    echo "✓ Collection '$QDRANT_COLLECTION' already exists"
else
    echo "Creating collection '$QDRANT_COLLECTION'..."

    curl -X PUT "$QDRANT_URL/collections/$QDRANT_COLLECTION" \
        -H "Content-Type: application/json" \
        -d "{
            \"vectors\": {
                \"size\": $EMBEDDING_DIM,
                \"distance\": \"Cosine\"
            }
        }"

    echo ""
    echo "✓ Collection '$QDRANT_COLLECTION' created successfully"
fi

echo ""
echo "Qdrant setup completed successfully!"