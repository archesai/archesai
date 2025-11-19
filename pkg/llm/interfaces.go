package llm

import (
	"context"
)

// ChatCompletionStream represents a streaming response.
type ChatCompletionStream interface {
	Recv() (ChatCompletionResponse, error)
	Close() error
}

// LLMService defines the interface for language model services.
type LLMService interface {
	// CreateChatCompletion creates a non-streaming chat completion.
	CreateChatCompletion(
		ctx context.Context,
		req ChatCompletionRequest,
	) (ChatCompletionResponse, error)

	// CreateChatCompletionStream creates a streaming chat completion.
	CreateChatCompletionStream(
		ctx context.Context,
		req ChatCompletionRequest,
	) (ChatCompletionStream, error)
}
