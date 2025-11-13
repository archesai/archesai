package llm

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/ssestream"

	"github.com/archesai/archesai/internal/core/models"
)

// OpenAILLM implements the LLM interface for OpenAI.
type OpenAILLM struct {
	client *openai.Client
}

// NewOpenAILLM creates a new OpenAI LLM client.
func NewOpenAILLM(apiKey string) *OpenAILLM {
	// If apiKey is empty, try to get from environment
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &OpenAILLM{client: &client}
}

// NewOpenAILLMWithHost creates an OpenAI client with custom host.
func NewOpenAILLMWithHost(apiKey string, host string) *OpenAILLM {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(host),
	)
	return &OpenAILLM{client: &client}
}

// convertToOpenAIMessages converts our generic Message type to OpenAI's message param type.
func convertToOpenAIMessages(
	messages []models.Message,
) []openai.ChatCompletionMessageParamUnion {
	var openAIMessages []openai.ChatCompletionMessageParamUnion

	for _, msg := range messages {
		switch msg.Role {
		case models.RoleSystem:
			openAIMessages = append(openAIMessages, openai.SystemMessage(msg.Content))
		case models.RoleUser:
			openAIMessages = append(openAIMessages, openai.UserMessage(msg.Content))
		case models.RoleAssistant:
			// AssistantMessage constructor doesn't support tool calls directly
			// Tool calls would need to be handled with a custom message type
			openAIMessages = append(openAIMessages, openai.AssistantMessage(msg.Content))
		case models.RoleFunction:
			// Tool/Function messages
			openAIMessages = append(openAIMessages, openai.ToolMessage(msg.Name, msg.Content))
		}
	}

	return openAIMessages
}

// convertToOpenAITools converts our Tool type to OpenAI's tool param.
func convertToOpenAITools(tools []models.LLMTool) []openai.ChatCompletionToolParam {
	if len(tools) == 0 {
		return nil
	}

	var openAITools []openai.ChatCompletionToolParam
	for _, tool := range tools {
		// Convert parameters to FunctionParameters (which is just map[string]any)
		var params openai.FunctionParameters
		if tool.Function.Parameters != nil {
			params = openai.FunctionParameters(tool.Function.Parameters)
		}

		openAITools = append(openAITools, openai.ChatCompletionToolParam{
			Type: openai.ChatCompletionToolParam{}.Type,
			Function: openai.FunctionDefinitionParam{
				Name:        tool.Function.Name,
				Description: openai.String(tool.Function.Description),
				Parameters:  params,
			},
		})
	}

	return openAITools
}

// CreateChatCompletion implements the LLM interface for OpenAI.
func (o *OpenAILLM) CreateChatCompletion(
	ctx context.Context,
	req models.ChatCompletionRequest,
) (models.ChatCompletionResponse, error) {
	params := openai.ChatCompletionNewParams{
		Model:    req.Model,
		Messages: convertToOpenAIMessages(req.Messages),
	}

	// Optional parameters
	if req.Temperature > 0 {
		params.Temperature = openai.Float(float64(req.Temperature))
	}
	if req.TopP > 0 {
		params.TopP = openai.Float(float64(req.TopP))
	}
	if req.MaxTokens > 0 {
		params.MaxTokens = openai.Int(int64(req.MaxTokens))
	}
	if len(req.Stop) > 0 {
		// Use the union type for stop sequences
		params.Stop = openai.ChatCompletionNewParamsStopUnion{
			OfStringArray: req.Stop,
		}
	}
	if req.PresencePenalty != 0 {
		params.PresencePenalty = openai.Float(float64(req.PresencePenalty))
	}

	// Add tools if present
	if len(req.Tools) > 0 {
		params.Tools = convertToOpenAITools(req.Tools)
	}

	// Make the API call
	completion, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return models.ChatCompletionResponse{}, fmt.Errorf("OpenAI API error: %w", err)
	}

	// Convert response
	var choices []models.Choice
	for _, c := range completion.Choices {
		msg := models.Message{
			Role:    models.RoleAssistant,
			Content: c.Message.Content,
		}

		// Convert tool calls
		if len(c.Message.ToolCalls) > 0 {
			for _, tc := range c.Message.ToolCalls {
				msg.ToolCalls = append(msg.ToolCalls, models.ToolCall{
					ID:   tc.ID,
					Type: string(tc.Type),
					Function: models.ToolCallFunction{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					},
				})
			}
		}

		choices = append(choices, models.Choice{
			Index:        int(c.Index),
			Message:      msg,
			FinishReason: c.FinishReason,
		})
	}

	return models.ChatCompletionResponse{
		ID:      completion.ID,
		Choices: choices,
		Usage: models.Usage{
			PromptTokens:     int(completion.Usage.PromptTokens),
			CompletionTokens: int(completion.Usage.CompletionTokens),
			TotalTokens:      int(completion.Usage.TotalTokens),
		},
	}, nil
}

// openAIStreamWrapper wraps the OpenAI stream.
type openAIStreamWrapper struct {
	stream *ssestream.Stream[openai.ChatCompletionChunk]
}

func newOpenAIStreamWrapper(
	stream *ssestream.Stream[openai.ChatCompletionChunk],
) *openAIStreamWrapper {
	return &openAIStreamWrapper{
		stream: stream,
	}
}

func (w *openAIStreamWrapper) Recv() (models.ChatCompletionResponse, error) {
	if !w.stream.Next() {
		err := w.stream.Err()
		if err == nil {
			return models.ChatCompletionResponse{}, io.EOF
		}
		return models.ChatCompletionResponse{}, err
	}

	chunk := w.stream.Current()

	var choices []models.Choice
	for _, c := range chunk.Choices {
		msg := models.Message{}

		// Handle delta content
		if c.Delta.Content != "" {
			msg.Content = c.Delta.Content
			msg.Role = models.RoleAssistant
		}

		// Handle delta tool calls
		if len(c.Delta.ToolCalls) > 0 {
			for _, tc := range c.Delta.ToolCalls {
				msg.ToolCalls = append(msg.ToolCalls, models.ToolCall{
					ID:   tc.ID,
					Type: tc.Type,
					Function: models.ToolCallFunction{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					},
				})
			}
		}

		choices = append(choices, models.Choice{
			Index:        int(c.Index),
			Message:      msg,
			FinishReason: c.FinishReason,
		})
	}

	return models.ChatCompletionResponse{
		ID:      chunk.ID,
		Choices: choices,
	}, nil
}

func (w *openAIStreamWrapper) Close() error {
	// The stream closes automatically when iteration is done
	return nil
}

// CreateChatCompletionStream implements the LLM interface for OpenAI streaming.
func (o *OpenAILLM) CreateChatCompletionStream(
	ctx context.Context,
	req models.ChatCompletionRequest,
) (ChatCompletionStream, error) {
	params := openai.ChatCompletionNewParams{
		Model:    req.Model,
		Messages: convertToOpenAIMessages(req.Messages),
	}

	// Optional parameters
	if req.Temperature > 0 {
		params.Temperature = openai.Float(float64(req.Temperature))
	}
	if req.TopP > 0 {
		params.TopP = openai.Float(float64(req.TopP))
	}
	if req.MaxTokens > 0 {
		params.MaxTokens = openai.Int(int64(req.MaxTokens))
	}
	if len(req.Stop) > 0 {
		// Use the union type for stop sequences
		params.Stop = openai.ChatCompletionNewParamsStopUnion{
			OfStringArray: req.Stop,
		}
	}
	if req.PresencePenalty != 0 {
		params.PresencePenalty = openai.Float(float64(req.PresencePenalty))
	}

	// Add tools if present
	if len(req.Tools) > 0 {
		params.Tools = convertToOpenAITools(req.Tools)
	}

	// Create streaming response
	stream := o.client.Chat.Completions.NewStreaming(ctx, params)

	return newOpenAIStreamWrapper(stream), nil
}
