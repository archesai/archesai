package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const deepseekAPIEndpoint = "https://api.deepseek.com/chat/completions"

// DeepSeekLLM implements the LLM interface for DeepSeek
type DeepSeekLLM struct {
	apiKey string
	client *http.Client
}

// NewDeepSeekLLM creates a new DeepSeek LLM client
func NewDeepSeekLLM(apiKey string) *DeepSeekLLM {
	return &DeepSeekLLM{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type deepseekMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	Name       string     `json:"name,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// Convert Message to deepseekMessage
func convertToDeepSeekMessage(msg Message) deepseekMessage {
	dsMsg := deepseekMessage{
		Role:      convertToDeepSeekRole(msg.Role),
		Content:   msg.Content,
		Name:      msg.Name,
		ToolCalls: msg.ToolCalls,
	}
	return dsMsg
}

// convertFromDeepSeekMessage converts deepseekMessage to Message (currently unused but kept for future use)

type deepseekRequest struct {
	Model            string            `json:"model"`
	Messages         []deepseekMessage `json:"messages"`
	FrequencyPenalty float32           `json:"frequency_penalty,omitempty"`
	MaxTokens        int               `json:"max_tokens,omitempty"`
	PresencePenalty  float32           `json:"presence_penalty,omitempty"`
	ResponseFormat   *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
	Stream      bool     `json:"stream,omitempty"`
	Temperature float32  `json:"temperature,omitempty"`
	TopP        float32  `json:"top_p,omitempty"`
	Tools       []Tool   `json:"tools,omitempty"`
	Stop        []string `json:"stop,omitempty"`
}

type deepseekResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type deepseekStreamResponse struct {
	ID      string         `json:"id"`
	Choices []StreamChoice `json:"choices"`
	Usage   Usage          `json:"usage"`
}

func convertToDeepSeekRole(role Role) string {
	if role == RoleFunction {
		return string(RoleTool)
	}
	return string(role)
}

func convertFromDeepSeekRole(role string) Role {
	if role == string(RoleTool) {
		return RoleFunction
	}
	return Role(role)
}

// CreateChatCompletion implements the LLM interface for DeepSeek
func (l *DeepSeekLLM) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
	// Convert messages to DeepSeek format
	var deepseekMessages []deepseekMessage
	var lastToolCalls []ToolCall

	for i, msg := range req.Messages {
		if msg.Role == RoleFunction {
			// For function responses, we need to find the corresponding tool call
			var toolCallID string
			for j := i - 1; j >= 0; j-- {
				if req.Messages[j].Role == RoleAssistant && len(req.Messages[j].ToolCalls) > 0 {
					for _, toolCall := range req.Messages[j].ToolCalls {
						if toolCall.Function.Name == msg.Name {
							toolCallID = toolCall.ID
							break
						}
					}
					if toolCallID != "" {
						break
					}
				}
			}
			if toolCallID == "" {
				// If we can't find a tool call ID, skip this message
				continue
			}
			dsMsg := convertToDeepSeekMessage(msg)
			dsMsg.ToolCallID = toolCallID
			deepseekMessages = append(deepseekMessages, dsMsg)
		} else {
			dsMsg := convertToDeepSeekMessage(msg)
			deepseekMessages = append(deepseekMessages, dsMsg)
			if msg.Role == RoleAssistant && len(msg.ToolCalls) > 0 {
				lastToolCalls = msg.ToolCalls
			}
		}
	}

	// If the last message had tool calls but no responses, skip the follow-up
	if len(lastToolCalls) > 0 {
		hasAllResponses := true
		for _, toolCall := range lastToolCalls {
			found := false
			for _, msg := range deepseekMessages {
				if msg.Role == string(RoleTool) && msg.ToolCallID == toolCall.ID {
					found = true
					break
				}
			}
			if !found {
				hasAllResponses = false
				break
			}
		}
		if !hasAllResponses {
			return ChatCompletionResponse{}, fmt.Errorf("missing tool responses")
		}
	}

	deepseekReq := deepseekRequest{
		Model:            req.Model,
		Messages:         deepseekMessages,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxTokens:        req.MaxTokens,
		PresencePenalty:  req.PresencePenalty,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		Tools:            req.Tools,
		Stop:             req.Stop,
	}

	// For follow-up responses after tool calls, disable tools to prevent loops
	if len(req.Messages) > 0 && req.Messages[len(req.Messages)-1].Role == RoleFunction {
		deepseekReq.Tools = nil
	}

	// Set default values if not provided
	if deepseekReq.Temperature == 0 {
		deepseekReq.Temperature = 0.7
	}
	if deepseekReq.TopP == 0 {
		deepseekReq.TopP = 0.95
	}
	if deepseekReq.MaxTokens == 0 {
		deepseekReq.MaxTokens = 2000
	}

	body, err := json.Marshal(deepseekReq)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", deepseekAPIEndpoint, bytes.NewReader(body))
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+l.apiKey)

	resp, err := l.client.Do(httpReq)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return ChatCompletionResponse{}, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var deepseekResp deepseekResponse
	if err := json.NewDecoder(resp.Body).Decode(&deepseekResp); err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return ChatCompletionResponse(deepseekResp), nil
}

type deepseekStreamWrapper struct {
	ctx            context.Context
	reader         *bufio.Reader
	response       *http.Response
	toolCallBuffer map[string]*ToolCall
}

func newDeepseekStreamWrapper(ctx context.Context, response *http.Response) *deepseekStreamWrapper {
	return &deepseekStreamWrapper{
		ctx:            ctx,
		reader:         bufio.NewReader(response.Body),
		response:       response,
		toolCallBuffer: make(map[string]*ToolCall),
	}
}

func (s *deepseekStreamWrapper) Close() error {
	return s.response.Body.Close()
}

func (s *deepseekStreamWrapper) Recv() (ChatCompletionResponse, error) {
	select {
	case <-s.ctx.Done():
		return ChatCompletionResponse{}, s.ctx.Err()
	default:
	}

	line, err := s.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return ChatCompletionResponse{}, io.EOF
		}
		return ChatCompletionResponse{}, fmt.Errorf("failed to read stream: %w", err)
	}

	// Remove "data: " prefix
	line = bytes.TrimPrefix(line, []byte("data: "))
	line = bytes.TrimSpace(line)

	// Check for stream end
	if bytes.Equal(line, []byte("[DONE]")) {
		return ChatCompletionResponse{}, io.EOF
	}

	var streamResp deepseekStreamResponse
	if err := json.Unmarshal(line, &streamResp); err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to unmarshal stream response: %w", err)
	}

	return ChatCompletionResponse{
		ID:      streamResp.ID,
		Choices: convertStreamChoicesToChoices(streamResp.Choices),
		Usage:   streamResp.Usage,
	}, nil
}

// CreateChatCompletionStream implements the LLM interface for DeepSeek streaming
func (l *DeepSeekLLM) CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (ChatCompletionStream, error) {
	// Convert messages to DeepSeek format
	var deepseekMessages []deepseekMessage
	var lastToolCalls []ToolCall

	for i, msg := range req.Messages {
		if msg.Role == RoleFunction {
			// For function responses, we need to find the corresponding tool call
			var toolCallID string
			for j := i - 1; j >= 0; j-- {
				if req.Messages[j].Role == RoleAssistant && len(req.Messages[j].ToolCalls) > 0 {
					for _, toolCall := range req.Messages[j].ToolCalls {
						if toolCall.Function.Name == msg.Name {
							toolCallID = toolCall.ID
							break
						}
					}
					if toolCallID != "" {
						break
					}
				}
			}
			if toolCallID == "" {
				// If we can't find a tool call ID, skip this message
				continue
			}
			dsMsg := convertToDeepSeekMessage(msg)
			dsMsg.ToolCallID = toolCallID
			deepseekMessages = append(deepseekMessages, dsMsg)
		} else {
			dsMsg := convertToDeepSeekMessage(msg)
			deepseekMessages = append(deepseekMessages, dsMsg)
			if msg.Role == RoleAssistant && len(msg.ToolCalls) > 0 {
				lastToolCalls = msg.ToolCalls
			}
		}
	}

	// If the last message had tool calls but no responses, skip the follow-up
	if len(lastToolCalls) > 0 {
		hasAllResponses := true
		for _, toolCall := range lastToolCalls {
			found := false
			for _, msg := range deepseekMessages {
				if msg.Role == string(RoleTool) && msg.ToolCallID == toolCall.ID {
					found = true
					break
				}
			}
			if !found {
				hasAllResponses = false
				break
			}
		}
		if !hasAllResponses {
			return nil, fmt.Errorf("missing tool responses")
		}
	}

	req.Stream = true
	deepseekReq := deepseekRequest{
		Model:            req.Model,
		Messages:         deepseekMessages,
		FrequencyPenalty: req.FrequencyPenalty,
		MaxTokens:        req.MaxTokens,
		PresencePenalty:  req.PresencePenalty,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		Tools:            req.Tools,
		Stop:             req.Stop,
		Stream:           true,
	}

	// For follow-up responses after tool calls, disable tools to prevent loops
	if len(req.Messages) > 0 && req.Messages[len(req.Messages)-1].Role == RoleFunction {
		deepseekReq.Tools = nil
	}

	// Set default values if not provided
	if deepseekReq.Temperature == 0 {
		deepseekReq.Temperature = 0.7
	}
	if deepseekReq.TopP == 0 {
		deepseekReq.TopP = 0.95
	}
	if deepseekReq.MaxTokens == 0 {
		deepseekReq.MaxTokens = 2000
	}

	body, err := json.Marshal(deepseekReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", deepseekAPIEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+l.apiKey)

	resp, err := l.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("failed to close response body: %v", closeErr)
		}
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return newDeepseekStreamWrapper(ctx, resp), nil
}

func convertStreamChoicesToChoices(streamChoices []StreamChoice) []Choice {
	choices := make([]Choice, len(streamChoices))
	for i, sc := range streamChoices {
		choices[i] = Choice{
			Index: sc.Index,
			Message: Message{
				Role:      convertFromDeepSeekRole(string(sc.Delta.Role)),
				Content:   sc.Delta.Content,
				ToolCalls: sc.Delta.ToolCalls,
			},
			FinishReason: sc.FinishReason,
		}
	}
	return choices
}
