// Package llm provides interfaces and implementations for interacting with language
package llm

import (
	"context"
	"fmt"
)

// ChatPersona represents a chat agent/assistant persona.
type ChatPersona struct {
	Name         string
	SystemPrompt string
	Model        string
	Temperature  float64
}

// ChatSession manages a conversation with an LLM.
type ChatSession struct {
	Messages []Message
	Persona  *ChatPersona
	Context  map[string]any
}

// ChatClient provides a high-level interface for chat conversations.
type ChatClient interface {
	// NewSession creates a new chat session
	NewSession(persona *ChatPersona) *ChatSession

	// SendMessage sends a message and returns the response
	SendMessage(
		ctx context.Context,
		session *ChatSession,
		content string,
	) (*Message, error)

	// SendMessageStream sends a message and returns a streaming response
	SendMessageStream(
		ctx context.Context,
		session *ChatSession,
		content string,
	) (ChatCompletionStream, error)
}

// DefaultChatClient implements ChatClient using an LLM provider.
type DefaultChatClient struct {
	llm Service
}

// NewChatClient creates a new chat client.
func NewChatClient(llm Service) ChatClient {
	return &DefaultChatClient{llm: llm}
}

// NewSession creates a new chat session.
func (c *DefaultChatClient) NewSession(persona *ChatPersona) *ChatSession {
	session := &ChatSession{
		Messages: []Message{},
		Persona:  persona,
		Context:  make(map[string]any),
	}

	// Add system message if persona has a system prompt
	if persona != nil && persona.SystemPrompt != "" {
		session.Messages = append(session.Messages, Message{
			Role:    RoleSystem,
			Content: persona.SystemPrompt,
		})
	}

	return session
}

// SendMessage sends a message and returns the response.
func (c *DefaultChatClient) SendMessage(
	ctx context.Context,
	session *ChatSession,
	content string,
) (*Message, error) {
	// Add user message to session
	userMsg := Message{
		Role:    RoleUser,
		Content: content,
	}
	session.Messages = append(session.Messages, userMsg)

	// Prepare request
	req := ChatCompletionRequest{
		Model:    session.Persona.Model,
		Messages: session.Messages,
	}

	if session.Persona.Temperature > 0 {
		req.Temperature = session.Persona.Temperature
	}

	// Get response from LLM
	resp, err := c.llm.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("chat completion error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices received")
	}

	// Extract assistant message
	assistantMsg := resp.Choices[0].Message

	// Add response to session history
	session.Messages = append(session.Messages, assistantMsg)

	return &assistantMsg, nil
}

// SendMessageStream sends a message and returns a streaming response.
func (c *DefaultChatClient) SendMessageStream(
	ctx context.Context,
	session *ChatSession,
	content string,
) (ChatCompletionStream, error) {
	// Add user message to session
	userMsg := Message{
		Role:    RoleUser,
		Content: content,
	}
	session.Messages = append(session.Messages, userMsg)

	// Prepare request
	req := ChatCompletionRequest{
		Model:    session.Persona.Model,
		Messages: session.Messages,
		Stream:   true,
	}

	if session.Persona.Temperature > 0 {
		req.Temperature = session.Persona.Temperature
	}

	// Get streaming response from LLM
	stream, err := c.llm.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("chat completion stream error: %w", err)
	}

	return stream, nil
}

// AddAssistantMessage adds an assistant message to the session (useful after streaming).
func (session *ChatSession) AddAssistantMessage(content string) {
	assistantMsg := Message{
		Role:    RoleAssistant,
		Content: content,
	}
	session.Messages = append(session.Messages, assistantMsg)
}

// ClearHistory clears the message history but keeps the system prompt.
func (session *ChatSession) ClearHistory() {
	var systemMessages []Message
	for _, msg := range session.Messages {
		if msg.Role == RoleSystem {
			systemMessages = append(systemMessages, msg)
		}
	}
	session.Messages = systemMessages
}

// GetLastMessage returns the last message in the session.
func (session *ChatSession) GetLastMessage() *Message {
	if len(session.Messages) == 0 {
		return nil
	}
	return &session.Messages[len(session.Messages)-1]
}

// DefaultPersonas provides some common chat personas.
var DefaultPersonas = struct {
	Assistant      *ChatPersona
	CodeHelper     *ChatPersona
	CreativeWriter *ChatPersona
	DataAnalyst    *ChatPersona
	Researcher     *ChatPersona
}{
	Assistant: &ChatPersona{
		Name:         "Assistant",
		SystemPrompt: "You are a helpful AI assistant. Be concise and informative.",
		Model:        "gpt-4",
		Temperature:  0.7,
	},
	CodeHelper: &ChatPersona{
		Name:         "CodeHelper",
		SystemPrompt: "You are a coding assistant. Help with programming questions, debugging, and best practices.",
		Model:        "gpt-4",
		Temperature:  0.3,
	},
	CreativeWriter: &ChatPersona{
		Name:         "CreativeWriter",
		SystemPrompt: "You are a creative writing assistant. Help with stories, poems, and creative content.",
		Model:        "gpt-4",
		Temperature:  0.8,
	},
	DataAnalyst: &ChatPersona{
		Name:         "DataAnalyst",
		SystemPrompt: "You are a data analysis expert. Help with data interpretation, statistics, and insights.",
		Model:        "gpt-4",
		Temperature:  0.4,
	},
	Researcher: &ChatPersona{
		Name:         "Researcher",
		SystemPrompt: "You are a research assistant. Help find information, summarize topics, and provide citations when possible.",
		Model:        "gpt-4",
		Temperature:  0.5,
	},
}
