package llm

import (
	"context"

	"github.com/archesai/archesai/internal/core/models"
	"github.com/archesai/archesai/internal/core/services"
)

// LLM defines the interface that all LLM providers must implement.
// This is kept here for backward compatibility and concrete implementations.
type LLM interface {
	CreateChatCompletion(
		ctx context.Context,
		req models.ChatCompletionRequest,
	) (models.ChatCompletionResponse, error)
	CreateChatCompletionStream(
		ctx context.Context,
		req models.ChatCompletionRequest,
	) (ChatCompletionStream, error)
}

// ChatCompletionStream represents a streaming response.
type ChatCompletionStream = services.ChatCompletionStream

// Ensure we can use LLM as services.LLMService
var _ services.LLMService = (LLM)(nil)
