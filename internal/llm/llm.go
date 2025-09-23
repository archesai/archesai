package llm

import (
	"context"
)

// Role represents the role of a message participant.
type Role string

// Role constants define the different types of roles in a conversation.
const (
	// RoleSystem represents system messages.
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleFunction  Role = "function"
	RoleTool      Role = "tool"
)

// Provider represents different LLM providers.
type Provider string

// Provider constants define supported AI providers.
const (
	// OpenAI represents the OpenAI provider.
	OpenAI          Provider = "OPEN_AI"
	Azure           Provider = "AZURE"
	AzureAD         Provider = "AZURE_AD"
	CloudflareAzure Provider = "CLOUDFLARE_AZURE"
	Gemini          Provider = "GEMINI"
	Claude          Provider = "CLAUDE"
	Ollama          Provider = "OLLAMA"
	DeepSeek        Provider = "DEEPSEEK"
)

// Message represents a single message in a chat conversation.
type Message struct {
	Role      Role       `json:"role"`
	Content   string     `json:"content"`
	Name      string     `json:"name,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ChatCompletionRequest represents a generic request for chat completion.
type ChatCompletionRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             int32     `json:"top_p,omitempty"`
	N                int32     `json:"n,omitempty"`
	Stop             []string  `json:"stop,omitempty"`
	MaxTokens        int32     `json:"max_tokens,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
	User             string    `json:"user,omitempty"`
	Tools            []Tool    `json:"tools,omitempty"`
	Stream           bool      `json:"stream,omitempty"`
}

// ChatCompletionResponse represents a generic response from chat completion.
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// LLM defines the interface that all LLM providers must implement.
type LLM interface {
	CreateChatCompletion(
		ctx context.Context,
		req ChatCompletionRequest,
	) (ChatCompletionResponse, error)
	CreateChatCompletionStream(
		ctx context.Context,
		req ChatCompletionRequest,
	) (ChatCompletionStream, error)
}

// ChatCompletionStream represents a streaming response.
type ChatCompletionStream interface {
	Recv() (ChatCompletionResponse, error)
	Close() error
}

// StreamDelta represents a delta in a streaming response.
type StreamDelta struct {
	Role       Role       `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// StreamChoice represents a choice in a streaming response.
type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        StreamDelta `json:"delta"`
	FinishReason string      `json:"finish_reason"`
}

// Tool represents a function that can be called by the LLM.
type Tool struct {
	Type     string    `json:"type"`
	Function *Function `json:"function,omitempty"`
}

// Function represents a function definition.
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolCall represents a tool/function call from the LLM.
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction represents the function details in a tool call.
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
