package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/ollama/ollama/api"
)

// OllamaLLM implements the LLM interface for Ollama
type OllamaLLM struct {
	client *api.Client
}

// NewOllamaLLM creates a new Ollama LLM client
func NewOllamaLLM() (*OllamaLLM, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}
	return &OllamaLLM{client: client}, nil
}

// NewOllamaLLMWithURL creates a new Ollama LLM client with a custom URL
func NewOllamaLLMWithURL(baseURL string) (*OllamaLLM, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	client := api.NewClient(parsedURL, nil)
	return &OllamaLLM{client: client}, nil
}

// convertToOllamaRole converts our Role type to Ollama's role string
func convertToOllamaRole(role Role) string {
	if role == RoleFunction {
		return "tool"
	}
	return string(role)
}

// convertFromOllamaRole converts Ollama's role string to our Role type
func convertFromOllamaRole(role string) Role {
	if role == "tool" {
		return RoleFunction
	}
	return Role(role)
}

// convertToOllamaMessages converts our generic Message type to Ollama's message format
func convertToOllamaMessages(messages []Message) []api.Message {
	ollamaMessages := make([]api.Message, len(messages))
	for i, msg := range messages {
		ollamaMessages[i] = api.Message{
			Role:      convertToOllamaRole(msg.Role),
			Content:   msg.Content,
			ToolCalls: convertToOllamaToolCalls(msg.ToolCalls),
		}
	}
	return ollamaMessages
}

// convertToOllamaTools converts our generic Tool type to Ollama's tool type
func convertToOllamaTools(tools []Tool) api.Tools {
	if len(tools) == 0 {
		return nil
	}

	ollamaTools := make([]api.Tool, len(tools))
	for i, tool := range tools {
		// Convert required array
		required := make([]string, len(tool.Function.Parameters["required"].([]interface{})))
		for i, v := range tool.Function.Parameters["required"].([]interface{}) {
			required[i] = v.(string)
		}

		// Convert properties map
		rawProps := tool.Function.Parameters["properties"].(map[string]interface{})
		properties := make(map[string]struct {
			Type        string   `json:"type"`
			Description string   `json:"description"`
			Enum        []string `json:"enum,omitempty"`
		})

		for propName, propValue := range rawProps {
			propMap := propValue.(map[string]interface{})
			prop := struct {
				Type        string   `json:"type"`
				Description string   `json:"description"`
				Enum        []string `json:"enum,omitempty"`
			}{
				Type:        propMap["type"].(string),
				Description: propMap["description"].(string),
			}

			// Handle optional enum field
			if enumVal, ok := propMap["enum"]; ok {
				enumInterface := enumVal.([]interface{})
				enumStrings := make([]string, len(enumInterface))
				for j, e := range enumInterface {
					enumStrings[j] = e.(string)
				}
				prop.Enum = enumStrings
			}

			properties[propName] = prop
		}

		ollamaTools[i] = api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Parameters: api.ToolFunctionParameters{
					Type:     tool.Function.Parameters["type"].(string),
					Required: required,
					// Properties: properties, // FIXME
				},
			},
		}
	}
	return ollamaTools
}

// convertToOllamaToolCalls converts our generic ToolCall type to Ollama's type
func convertToOllamaToolCalls(toolCalls []ToolCall) []api.ToolCall {
	if len(toolCalls) == 0 {
		return nil
	}

	calls := make([]api.ToolCall, len(toolCalls))
	for i, call := range toolCalls {
		// Convert map[string]interface{} to api.ToolCallFunctionArguments
		args := make(map[string]any)
		err := json.Unmarshal([]byte(call.Function.Arguments), &args)
		if err != nil {
			// Handle error (e.g., log it or use a default empty map)
			args = make(map[string]any)
		}

		calls[i] = api.ToolCall{
			Function: api.ToolCallFunction{
				Name:      call.Function.Name,
				Arguments: args,
			},
		}
	}
	return calls
}

// convertFromOllamaToolCalls converts Ollama's tool calls to our generic type
func convertFromOllamaToolCalls(toolCalls []api.ToolCall) []ToolCall {
	if len(toolCalls) == 0 {
		return nil
	}

	calls := make([]ToolCall, len(toolCalls))
	for i, call := range toolCalls {
		// Convert api.ToolCallFunctionArguments to map[string]interface{}

		calls[i] = ToolCall{
			ID:   call.Function.Name, // Using name as ID since Ollama doesn't have a separate ID field
			Type: "function",
			Function: ToolCallFunction{
				Name:      call.Function.Name,
				Arguments: call.Function.Arguments.String(),
			},
		}
	}
	return calls
}

// CreateChatCompletion implements the LLM interface for Ollama
func (o *OllamaLLM) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
	stream := false
	ollamaReq := &api.ChatRequest{
		Model:    req.Model,
		Messages: convertToOllamaMessages(req.Messages),
		Stream:   &stream,
		Tools:    convertToOllamaTools(req.Tools),
		Options:  make(map[string]interface{}),
	}

	var response ChatCompletionResponse
	var finalMessage Message

	err := o.client.Chat(ctx, ollamaReq, func(resp api.ChatResponse) error {
		if resp.Done {
			finalMessage = Message{
				Role:      convertFromOllamaRole(resp.Message.Role),
				Content:   resp.Message.Content,
				ToolCalls: convertFromOllamaToolCalls(resp.Message.ToolCalls),
			}
		}
		return nil
	})

	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("ollama chat completion failed: %w", err)
	}

	response.Choices = []Choice{
		{
			Index:        0,
			Message:      finalMessage,
			FinishReason: "stop",
		},
	}

	return response, nil
}

type ollamaStreamWrapper struct {
	ctx            context.Context
	client         *api.Client
	req            *api.ChatRequest
	done           bool
	content        string
	toolCallBuffer map[string]*ToolCall
}

func newOllamaStreamWrapper(ctx context.Context, client *api.Client, req *api.ChatRequest) *ollamaStreamWrapper {
	return &ollamaStreamWrapper{
		ctx:            ctx,
		client:         client,
		req:            req,
		toolCallBuffer: make(map[string]*ToolCall),
	}
}

func (s *ollamaStreamWrapper) Recv() (ChatCompletionResponse, error) {
	if s.done {
		return ChatCompletionResponse{}, io.EOF
	}

	var response ChatCompletionResponse
	err := s.client.Chat(s.ctx, s.req, func(resp api.ChatResponse) error {
		if resp.Done {
			s.done = true
			return io.EOF
		}

		s.content += resp.Message.Content
		response.Choices = []Choice{
			{
				Index: 0,
				Message: Message{
					Role:      convertFromOllamaRole(resp.Message.Role),
					Content:   resp.Message.Content,
					ToolCalls: convertFromOllamaToolCalls(resp.Message.ToolCalls),
				},
				FinishReason: "",
			},
		}
		return nil
	})

	if err == io.EOF {
		response.Choices[0].FinishReason = "stop"
		return response, nil
	}

	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("ollama stream failed: %w", err)
	}

	return response, nil
}

func (s *ollamaStreamWrapper) Close() error {
	return nil
}

// CreateChatCompletionStream implements the LLM interface for Ollama streaming
func (o *OllamaLLM) CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (ChatCompletionStream, error) {
	stream := true
	ollamaReq := &api.ChatRequest{
		Model:    req.Model,
		Messages: convertToOllamaMessages(req.Messages),
		Stream:   &stream,
		Tools:    convertToOllamaTools(req.Tools),
		Options:  make(map[string]interface{}),
	}

	return newOllamaStreamWrapper(ctx, o.client, ollamaReq), nil
}
