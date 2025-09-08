-- +goose Up
-- +goose StatementBegin
-- Enable pgvector extension for AI embeddings
CREATE EXTENSION IF NOT EXISTS vector;

-- Add example vector column to demonstrate functionality
-- This can be used for semantic search, similarity matching, etc.
ALTER TABLE artifact ADD COLUMN embedding vector(1536); -- OpenAI embedding dimension
CREATE INDEX ON artifact USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove vector column and index
DROP INDEX IF EXISTS artifact_embedding_idx;
ALTER TABLE artifact DROP COLUMN IF EXISTS embedding;

-- Drop pgvector extension
DROP EXTENSION IF EXISTS vector;
-- +goose StatementEnd
