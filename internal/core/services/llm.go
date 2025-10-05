package services

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// ChatCompletionStream represents a streaming response.
type ChatCompletionStream interface {
	Recv() (valueobjects.ChatCompletionResponse, error)
	Close() error
}

// LLMService defines the interface for language model services.
type LLMService interface {
	// CreateChatCompletion creates a non-streaming chat completion.
	CreateChatCompletion(
		ctx context.Context,
		req valueobjects.ChatCompletionRequest,
	) (valueobjects.ChatCompletionResponse, error)

	// CreateChatCompletionStream creates a streaming chat completion.
	CreateChatCompletionStream(
		ctx context.Context,
		req valueobjects.ChatCompletionRequest,
	) (ChatCompletionStream, error)
}
