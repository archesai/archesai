package llm

import (
	"context"
	"fmt"
)

// Temporary stub implementations until the actual implementations are fixed

// ClaudeLLM stub implementation
type ClaudeLLM struct{}

// NewClaudeLLM creates a new Claude LLM stub instance
func NewClaudeLLM(_ string) *ClaudeLLM {
	return &ClaudeLLM{}
}

// CreateChatCompletion implements the LLM interface for Claude (stub)
func (c *ClaudeLLM) CreateChatCompletion(_ context.Context, _ ChatCompletionRequest) (ChatCompletionResponse, error) {
	return ChatCompletionResponse{}, fmt.Errorf("claude implementation temporarily disabled")
}

// CreateChatCompletionStream implements the LLM interface for Claude streaming (stub)
func (c *ClaudeLLM) CreateChatCompletionStream(_ context.Context, _ ChatCompletionRequest) (ChatCompletionStream, error) {
	return nil, fmt.Errorf("claude streaming implementation temporarily disabled")
}

// GeminiLLM stub implementation
type GeminiLLM struct{}

// NewGeminiLLM creates a new Gemini LLM stub instance
func NewGeminiLLM(_ string, _ ...interface{}) (*GeminiLLM, error) {
	return &GeminiLLM{}, nil
}

// CreateChatCompletion implements the LLM interface for Gemini (stub)
func (g *GeminiLLM) CreateChatCompletion(_ context.Context, _ ChatCompletionRequest) (ChatCompletionResponse, error) {
	return ChatCompletionResponse{}, fmt.Errorf("gemini implementation temporarily disabled")
}

// CreateChatCompletionStream implements the LLM interface for Gemini streaming (stub)
func (g *GeminiLLM) CreateChatCompletionStream(_ context.Context, _ ChatCompletionRequest) (ChatCompletionStream, error) {
	return nil, fmt.Errorf("gemini streaming implementation temporarily disabled")
}
