package swarm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/llm"
)

var (
	// ErrNilAgent is returned when a nil agent is provided
	ErrNilAgent = errors.New("agent cannot be nil")
	// ErrEmptyMessages is returned when message history is empty
	ErrEmptyMessages = errors.New("message history cannot be empty")
	// ErrLLMClientNotReady is returned when the LLM client is not initialized
	ErrLLMClientNotReady = errors.New("LLM client is not initialized")
	// ErrInvalidProvider is returned when an invalid LLM provider is specified
	ErrInvalidProvider = errors.New("invalid LLM provider specified")
	// ErrNoChoicesInResp is returned when there are no choices in the LLM response
	ErrNoChoicesInResp = errors.New("no choices in LLM response")
	// ErrMessageTooLong is returned when a message exceeds the maximum token limit
	ErrMessageTooLong = errors.New("message exceeds maximum token limit")
)

// Swarm represents the main structure
type Swarm struct {
	client       llm.LLM
	tokenCounter func(string) int // Optional token counter function
	initialized  bool             // Flag to check if Swarm is properly initialized
	config       *Config          // Configuration settings
}

// Config holds configuration options for Swarm
type Config struct {
	MaxRetries        int
	RetryBackoff      time.Duration
	RequestTimeout    time.Duration
	MaxTokens         int
	DefaultModel      string
	Debug             bool
	LogLevel          LogLevel
	TokenLimits       map[string]int // Model-specific token limits
	FailureHandlers   []FailureHandler
	RateLimitStrategy RateLimitStrategy
}

// LogLevel represents the level of logging
type LogLevel int

const (
	// LogSilent suppresses all log output
	LogSilent LogLevel = iota
	// LogError logs only errors
	LogError
	// LogWarning logs warnings and errors
	LogWarning
	// LogInfo logs info, warnings, and errors
	LogInfo
	// LogDebug logs debug messages and above
	LogDebug
	// LogTrace logs all messages including trace
	LogTrace
)

// FailureHandler defines a function to handle specific failures
type FailureHandler func(error) (bool, error)

// RateLimitStrategy defines how rate limits are handled
type RateLimitStrategy int

const (
	// RateLimitRetry retries requests when rate limited
	RateLimitRetry RateLimitStrategy = iota
	// RateLimitFail fails immediately when rate limited
	RateLimitFail
	// RateLimitQueue queues requests when rate limited
	RateLimitQueue
)

// DefaultConfig returns default configuration values
func DefaultConfig() *Config {
	return &Config{
		MaxRetries:     3,
		RetryBackoff:   time.Second,
		RequestTimeout: 60 * time.Second,
		MaxTokens:      4096,
		DefaultModel:   "gpt-3.5-turbo",
		Debug:          false,
		LogLevel:       LogError,
		TokenLimits: map[string]int{
			"gpt-3.5-turbo": 4096,
			"gpt-4":         8192,
			"gpt-4o":        128000,
			"claude-3-opus": 200000,
		},
		RateLimitStrategy: RateLimitRetry,
	}
}

// NewSwarm initializes a new Swarm instance with an LLM client
func NewSwarm(apiKey string, provider llm.Provider) *Swarm {
	return NewSwarmWithConfig(apiKey, provider, DefaultConfig())
}

// NewSwarmWithConfig initializes a new Swarm with custom configuration
func NewSwarmWithConfig(apiKey string, provider llm.Provider, config *Config) *Swarm {
	if apiKey == "" {
		log.Println("Warning: Empty API key provided")
		return &Swarm{
			initialized: false,
			config:      config,
		}
	}

	var client llm.LLM
	var err error

	switch provider {
	case llm.OpenAI:
		client = llm.NewOpenAILLM(apiKey)
	case llm.Gemini:
		client, err = llm.NewGeminiLLM(apiKey)
		if err != nil {
			log.Printf("Failed to create Gemini client: %v", err)
			return &Swarm{
				initialized: false,
				config:      config,
			}
		}
	case llm.Claude:
		client = llm.NewClaudeLLM(apiKey)
	case llm.Ollama:
		client, err = llm.NewOllamaLLM()
		if err != nil {
			log.Printf("Failed to create Ollama client: %v", err)
			return &Swarm{
				initialized: false,
				config:      config,
			}
		}
	case llm.DeepSeek:
		client = llm.NewDeepSeekLLM(apiKey)
	default:
		log.Printf("Unsupported LLM provider: %v", provider)
		return &Swarm{
			initialized: false,
			config:      config,
		}
	}

	// Verify the client was created properly
	if client == nil {
		log.Println("Warning: Failed to initialize LLM client")
		return &Swarm{
			initialized: false,
			config:      config,
		}
	}

	return &Swarm{
		client:      client,
		initialized: true,
		config:      config,
	}
}

// NewSwarmWithHost creates a Swarm with a custom host
func NewSwarmWithHost(apiKey, host string, provider llm.Provider) *Swarm {
	if provider == llm.OpenAI {
		client := llm.NewOpenAILLMWithHost(apiKey, host)
		if client == nil {
			log.Println("Warning: Failed to initialize OpenAI client with custom host")
			return &Swarm{
				initialized: false,
				config:      DefaultConfig(),
			}
		}
		return &Swarm{
			client:      client,
			initialized: true,
			config:      DefaultConfig(),
		}
	}
	log.Printf("Custom host %s not supported for provider: %v", host, provider)
	return NewSwarm(apiKey, provider)
}

// NewSwarmWithCustomProvider creates a Swarm with a custom LLM provider implementation
func NewSwarmWithCustomProvider(providerImpl llm.LLM, config *Config) *Swarm {
	return &Swarm{
		client:      providerImpl,
		initialized: true,
		config:      config,
	}
}

// SetTokenCounter sets a function to count tokens in messages
func (s *Swarm) SetTokenCounter(counter func(string) int) {
	s.tokenCounter = counter
}

// IsInitialized returns whether the Swarm is properly initialized
func (s *Swarm) IsInitialized() bool {
	return s.initialized && s.client != nil
}

// ValidateConnection tests the LLM connection with a simple request
func (s *Swarm) ValidateConnection(ctx context.Context) error {
	if !s.IsInitialized() {
		return ErrLLMClientNotReady
	}

	// Create a simple test request
	testRequest := llm.ChatCompletionRequest{
		Model: s.config.DefaultModel,
		Messages: []llm.Message{
			{Role: llm.RoleUser, Content: "Test connection"},
		},
		MaxTokens: 5,
	}

	// Attempt to send the request
	_, err := s.client.CreateChatCompletion(ctx, testRequest)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	return nil
}

// getChatCompletion requests a chat completion from the LLM with retries and error handling
// nolint:unused
func (s *Swarm) getChatCompletion(
	ctx context.Context,
	agent *Agent,
	history []llm.Message,
	contextVariables map[string]interface{},
	modelOverride string,
	_ bool,
	debug bool,
) (llm.ChatCompletionResponse, error) {
	// Validate inputs
	if !s.IsInitialized() {
		return llm.ChatCompletionResponse{}, ErrLLMClientNotReady
	}

	if agent == nil {
		return llm.ChatCompletionResponse{}, ErrNilAgent
	}

	if len(history) == 0 {
		// Instead of failing, create an empty initial message
		history = []llm.Message{}
	}

	// Prepare the initial system message with agent instructions
	instructions := agent.Instructions
	if agent.InstructionsFunc != nil {
		if contextVariables == nil {
			contextVariables = make(map[string]interface{})
		}
		instructions = agent.InstructionsFunc(contextVariables)
	}

	if instructions != "" {
		// Add system message at the beginning
		messages := []llm.Message{
			{
				Role:    llm.RoleSystem,
				Content: instructions,
			},
		}
		messages = append(messages, history...)
		history = messages
	}

	// Build tool definitions from agent's functions
	var tools []llm.Tool
	if agent.Functions != nil {
		for _, af := range agent.Functions {
			def := FunctionToDefinition(af)
			tools = append(tools, llm.Tool{
				Type: "function",
				Function: &llm.Function{
					Name:        def.Name,
					Description: def.Description,
					Parameters:  def.Parameters,
				},
			})
		}
	}

	// Prepare the chat completion request
	model := agent.Model
	if modelOverride != "" {
		model = modelOverride
	}

	// Use default model if none specified
	if model == "" {
		model = s.config.DefaultModel
	}

	req := llm.ChatCompletionRequest{
		Model:    model,
		Messages: history,
		Tools:    tools,
	}

	if debug {
		log.Printf("Debug - Model: %s, Messages: %d, Tools: %d\n",
			model, len(history), len(tools))
	}

	// Implement retry logic
	var lastErr error
	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		if attempt > 0 && s.config.Debug {
			log.Printf("Retry attempt %d after error: %v", attempt, lastErr)
		}

		// Create a timeout context for this request if none was provided
		requestCtx := ctx
		if _, hasDeadline := ctx.Deadline(); !hasDeadline {
			var cancel context.CancelFunc
			requestCtx, cancel = context.WithTimeout(ctx, s.config.RequestTimeout)
			defer cancel()
		}

		// Call the LLM to get a chat completion
		resp, err := s.client.CreateChatCompletion(requestCtx, req)
		if err == nil {
			// Success
			return resp, nil
		}

		// Handle the error
		lastErr = err

		// Check for rate limit errors and apply the rate limit strategy
		if isRateLimitError(err) {
			switch s.config.RateLimitStrategy {
			case RateLimitFail:
				return llm.ChatCompletionResponse{}, fmt.Errorf("rate limit exceeded: %w", err)
			case RateLimitQueue:
				// Implement exponential backoff
				backoff := s.config.RetryBackoff * time.Duration(1<<uint(attempt))
				if s.config.Debug {
					log.Printf("Rate limit hit, backing off for %v", backoff)
				}
				select {
				case <-ctx.Done():
					return llm.ChatCompletionResponse{}, ctx.Err()
				case <-time.After(backoff):
					// Continue to next retry
				}
			default: // RateLimitRetry
				// Simple retry with backoff
				backoff := s.config.RetryBackoff * time.Duration(attempt+1)
				if s.config.Debug {
					log.Printf("Backing off for %v before retry", backoff)
				}
				time.Sleep(backoff)
			}
		} else if isFatalError(err) {
			// Don't retry fatal errors
			return llm.ChatCompletionResponse{}, err
		} else {
			// For other errors, apply backoff
			backoff := s.config.RetryBackoff * time.Duration(attempt+1)
			time.Sleep(backoff)
		}
	}

	// All retries failed
	return llm.ChatCompletionResponse{}, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// isRateLimitError checks if an error is related to rate limiting
// nolint:unused
func isRateLimitError(err error) bool {
	return err != nil && (strings.Contains(strings.ToLower(err.Error()), "rate limit") ||
		strings.Contains(strings.ToLower(err.Error()), "too many requests") ||
		strings.Contains(strings.ToLower(err.Error()), "429"))
}

// isFatalError checks if an error is fatal and should not be retried
// nolint:unused
func isFatalError(err error) bool {
	return err != nil && (strings.Contains(strings.ToLower(err.Error()), "invalid auth") ||
		strings.Contains(strings.ToLower(err.Error()), "authentication") ||
		strings.Contains(strings.ToLower(err.Error()), "not found") ||
		strings.Contains(strings.ToLower(err.Error()), "invalid model"))
}

// Helper function to clone a slice of messages
// nolint:unused
func cloneMessages(msgs []llm.Message) []llm.Message {
	if msgs == nil {
		return []llm.Message{}
	}

	cloned := make([]llm.Message, len(msgs))
	copy(cloned, msgs)
	return cloned
}

// Helper function to return the last message in a slice
// nolint:unused
func lastMessage(msgs []llm.Message) *llm.Message {
	if len(msgs) == 0 {
		return nil
	}
	return &msgs[len(msgs)-1]
}

// handleToolCall processes a tool call and ensures proper context
func (s *Swarm) handleToolCall(
	_ context.Context,
	toolCall *llm.ToolCall,
	agent *Agent,
	contextVariables map[string]interface{},
	debug bool,
) Response {
	toolName := toolCall.Function.Name
	argsJSON := toolCall.Function.Arguments

	// Parse the tool call arguments
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		errorMsg := fmt.Sprintf("Error parsing tool call arguments: %v", err)
		if debug {
			log.Println(errorMsg)
		}
		return Response{
			Messages: []llm.Message{
				{
					Role:    llm.RoleFunction,
					Content: errorMsg,
					Name:    toolName,
				},
			},
		}
	}

	if debug {
		log.Printf("Processing tool call: %s with arguments %v\n", toolName, args)
	}

	// Find the corresponding function in the agent's functions
	var functionFound *AgentFunction
	for _, af := range agent.Functions {
		if af.Name == toolName {
			functionFound = &af
			break
		}
	}

	// Handle case where function is not found
	if functionFound == nil {
		errorMsg := fmt.Sprintf("Error: Tool %s not found", toolName)
		if debug {
			log.Println(errorMsg)
		}
		return Response{
			Messages: []llm.Message{
				{
					Role:    llm.RoleFunction,
					Content: errorMsg,
					Name:    toolName,
				},
			},
		}
	}

	// Execute the function
	result := functionFound.Function(args, contextVariables)

	// Create a message with the tool result
	var resultContent string
	if result.Error != nil {
		resultContent = fmt.Sprintf("Error: %v", result.Error)
	} else {
		resultContent = fmt.Sprintf("%v", result.Data)
	}

	// Create function response message properly formatted for tool call
	toolResultMessage := llm.Message{
		Role:    llm.RoleFunction,
		Content: resultContent,
		Name:    toolName,
	}

	// Return the response with the tool result
	return Response{
		Messages:         []llm.Message{toolResultMessage},
		Agent:            result.Agent,
		ContextVariables: contextVariables,
	}
}

// handleToolCalls handles multiple tool calls with correct context forwarding
func (s *Swarm) handleToolCalls(
	ctx context.Context,
	toolCalls []llm.ToolCall,
	history []llm.Message,
	agent *Agent,
	contextVariables map[string]interface{},
	modelOverride string,
	_ bool,
	debug bool,
	_ bool,
) ([]ToolResult, []llm.Message, *Agent, error) {
	var toolResults []ToolResult
	updatedAgent := agent
	updatedHistory := make([]llm.Message, len(history))
	copy(updatedHistory, history)

	// Execute tools sequentially for now for simplicity
	for _, toolCall := range toolCalls {
		// Execute the tool call
		toolResp := s.handleToolCall(ctx, &toolCall, updatedAgent, contextVariables, debug)

		// Parse arguments for the result
		var args interface{}
		_ = json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

		// Record the tool result
		toolResults = append(toolResults, ToolResult{
			ToolName: toolCall.Function.Name,
			Args:     args,
			Result: Result{
				Success: true,
				Data:    toolResp.Messages[0].Content,
				Error:   nil,
				Agent:   toolResp.Agent,
			},
		})

		// Add the function result to history with proper role and name
		updatedHistory = append(updatedHistory, llm.Message{
			Role:    llm.RoleFunction,
			Content: toolResp.Messages[0].Content,
			Name:    toolCall.Function.Name,
		})

		// Update agent if needed
		if toolResp.Agent != nil {
			updatedAgent = toolResp.Agent
		}
	}

	// Only try to get follow-up if we executed at least one tool successfully
	if len(toolResults) > 0 {
		// CRITICAL FIX: Create a new request with the updated history INCLUDING function results
		followUpReq := llm.ChatCompletionRequest{
			Model:    updatedAgent.Model,
			Messages: updatedHistory,
		}

		if modelOverride != "" {
			followUpReq.Model = modelOverride
		}

		// Create a timeout for the follow-up request
		followUpCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if debug {
			log.Printf("Getting follow-up with %d messages in history", len(updatedHistory))

			// Print last few messages for debugging
			lastN := 3
			if len(updatedHistory) < lastN {
				lastN = len(updatedHistory)
			}

			log.Printf("Last %d messages:", lastN)
			for i := len(updatedHistory) - lastN; i < len(updatedHistory); i++ {
				msg := updatedHistory[i]
				log.Printf("[%s] %s: %s", msg.Role, msg.Name, truncateString(msg.Content, 50))
			}
		}

		// Get the follow-up response with proper context
		followUpResp, err := s.client.CreateChatCompletion(followUpCtx, followUpReq)

		if err != nil {
			if debug {
				log.Printf("Error getting follow-up: %v", err)
			}
			// Continue without follow-up rather than failing
		} else if len(followUpResp.Choices) > 0 {
			// Add follow-up response to history
			followUpMessage := followUpResp.Choices[0].Message

			// Only add if it has content and only append content (no tools)
			if followUpMessage.Content != "" {
				cleanedFollowUp := llm.Message{
					Role:    followUpMessage.Role,
					Content: followUpMessage.Content,
				}
				updatedHistory = append(updatedHistory, cleanedFollowUp)

				if debug {
					log.Printf("Added follow-up: %s", truncateString(followUpMessage.Content, 50))
				}
			} else if debug {
				log.Println("Follow-up was empty")
			}
		}
	}

	return toolResults, updatedHistory, updatedAgent, nil
}

// Helper function to truncate strings for debugging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Run is the main entry point for agent execution
func (s *Swarm) Run(
	ctx context.Context,
	agent *Agent,
	messages []llm.Message,
	contextVariables map[string]interface{},
	modelOverride string,
	stream bool,
	debug bool,
	_ int,
	executeTools bool,
) (Response, error) {
	// Validate inputs
	if agent == nil {
		return Response{}, fmt.Errorf("agent cannot be nil")
	}

	// Use a cloned copy of messages for history
	history := make([]llm.Message, len(messages))
	copy(history, messages)

	if contextVariables == nil {
		contextVariables = make(map[string]interface{})
	}

	// Set up context variables
	if agent.InstructionsFunc != nil {
		agent.Instructions = agent.InstructionsFunc(contextVariables)
	}

	// Add system instruction as first message if not already present
	hasSystemMessage := false
	for _, msg := range history {
		if msg.Role == llm.RoleSystem {
			hasSystemMessage = true
			break
		}
	}

	if !hasSystemMessage && agent.Instructions != "" {
		newHistory := make([]llm.Message, 0, len(history)+1)
		newHistory = append(newHistory, llm.Message{
			Role:    llm.RoleSystem,
			Content: agent.Instructions,
		})
		newHistory = append(newHistory, history...)
		history = newHistory
	}

	// Get chat completion from LLM
	if debug {
		log.Printf("Getting initial response with %d messages", len(history))
	}

	model := agent.Model
	if modelOverride != "" {
		model = modelOverride
	}

	// Prepare tools for the request
	var tools []llm.Tool
	if agent.Functions != nil {
		for _, fn := range agent.Functions {
			tools = append(tools, llm.Tool{
				Type: "function",
				Function: &llm.Function{
					Name:        fn.Name,
					Description: fn.Description,
					Parameters:  fn.Parameters,
				},
			})
		}
	}

	// Create the request
	req := llm.ChatCompletionRequest{
		Model:    model,
		Messages: history,
		Tools:    tools,
	}

	// Get initial response
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return Response{}, fmt.Errorf("chat completion error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return Response{}, fmt.Errorf("no choices in response")
	}

	// Extract the response
	choice := resp.Choices[0]
	history = append(history, choice.Message)

	// Handle tool calls if present and execution is enabled
	if len(choice.Message.ToolCalls) > 0 && executeTools {
		if debug {
			log.Printf("Handling %d tool calls", len(choice.Message.ToolCalls))
		}

		// Execute tools and get the updated history including follow-up
		toolResults, updatedHistory, updatedAgent, err := s.handleToolCalls(
			ctx, choice.Message.ToolCalls, history, agent,
			contextVariables, modelOverride, stream, debug,
			agent.ParallelToolCalls)

		if err != nil {
			return Response{}, fmt.Errorf("tool execution error: %v", err)
		}

		// Calculate which messages to return (only the new ones)
		newMessages := updatedHistory[len(messages):]

		return Response{
			Messages:         newMessages,
			Agent:            updatedAgent,
			ContextVariables: contextVariables,
			ToolResults:      toolResults,
		}, nil
	}

	// No tool calls - just return the normal response
	return Response{
		Messages:         history[len(messages):],
		Agent:            agent,
		ContextVariables: contextVariables,
	}, nil
}

// handleToolCallsParallel executes multiple tool calls concurrently
// nolint:unused
func (s *Swarm) handleToolCallsParallel(
	ctx context.Context,
	toolCalls []llm.ToolCall,
	history []llm.Message,
	agent *Agent,
	contextVariables map[string]interface{},
	_ string,
	_ bool,
	debug bool,
) ([]ToolResult, []llm.Message, *Agent, error) {
	type toolCallResult struct {
		index  int
		result Response
		err    error
	}

	resultChan := make(chan toolCallResult, len(toolCalls))
	updatedHistory := make([]llm.Message, len(history))
	copy(updatedHistory, history)
	updatedAgent := agent

	// Create a cancellable context for all tool calls
	execCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Launch goroutines for each tool call
	for i, toolCall := range toolCalls {
		go func(idx int, tc llm.ToolCall) {
			toolResp := s.handleToolCall(execCtx, &tc, agent, contextVariables, debug)
			resultChan <- toolCallResult{index: idx, result: toolResp, err: nil}
		}(i, toolCall)
	}

	// Collect results
	var toolResults []ToolResult
	agentTransferred := false

	for i := 0; i < len(toolCalls); i++ {
		select {
		case <-ctx.Done():
			return nil, history, agent, ctx.Err()
		case result := <-resultChan:
			if result.err != nil {
				if debug {
					log.Printf("Error in tool call %d: %v", result.index, result.err)
				}
				continue
			}

			// Get the original tool call
			toolCall := toolCalls[result.index]
			var args interface{}
			_ = json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

			// Add to tool results
			toolResults = append(toolResults, ToolResult{
				ToolName: toolCall.Function.Name,
				Args:     args,
				Result: Result{
					Success: true,
					Data:    result.result.Messages[0].Content,
					Error:   nil,
					Agent:   result.result.Agent,
				},
			})

			// Add to history
			updatedHistory = append(updatedHistory, llm.Message{
				Role:    llm.RoleFunction,
				Content: result.result.Messages[0].Content,
				Name:    toolCall.Function.Name,
			})

			// Only update agent if not already transferred
			if result.result.Agent != nil && !agentTransferred {
				updatedAgent = result.result.Agent
				agentTransferred = true
			}

			// Store in memory
			if agent.Memory != nil {
				agent.Memory.AddMemory(Memory{
					Content: fmt.Sprintf("Tool %s call with args: %v, result: %s",
						toolCall.Function.Name, args, result.result.Messages[0].Content),
					Type:       "tool_call",
					Context:    map[string]interface{}{"tool": toolCall.Function.Name},
					Timestamp:  time.Now(),
					Importance: 0.7,
				})
			}
		}
	}

	return toolResults, updatedHistory, updatedAgent, nil
}
