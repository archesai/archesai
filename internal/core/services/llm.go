package services

import (
	"context"

	"github.com/archesai/archesai/internal/core/models"
)

// ChatCompletionStream represents a streaming response.
type ChatCompletionStream interface {
	Recv() (models.ChatCompletionResponse, error)
	Close() error
}

// LLMService defines the interface for language model services.
type LLMService interface {
	// CreateChatCompletion creates a non-streaming chat completion.
	CreateChatCompletion(
		ctx context.Context,
		req models.ChatCompletionRequest,
	) (models.ChatCompletionResponse, error)

	// CreateChatCompletionStream creates a streaming chat completion.
	CreateChatCompletionStream(
		ctx context.Context,
		req models.ChatCompletionRequest,
	) (ChatCompletionStream, error)
}
