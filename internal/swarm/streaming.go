package swarm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/archesai/archesai/internal/llm"
)

// StreamHandler represents a handler for streaming responses
type StreamHandler interface {
	OnStart()
	OnToken(token string)
	OnToolCall(toolCall llm.ToolCall)
	OnComplete(message llm.Message)
	OnError(err error)
}

// DefaultStreamHandler provides a basic implementation of StreamHandler
type DefaultStreamHandler struct{}

// OnStart is called when streaming starts
func (h *DefaultStreamHandler) OnStart() {}

// OnToken is called for each token received
func (h *DefaultStreamHandler) OnToken(_ string) {}

// OnToolCall is called when a tool call is received
func (h *DefaultStreamHandler) OnToolCall(_ llm.ToolCall) {}

// OnComplete is called when streaming completes
func (h *DefaultStreamHandler) OnComplete(_ llm.Message) {}

// OnError is called when an error occurs
func (h *DefaultStreamHandler) OnError(_ error) {}

// StreamingResponse handles streaming chat completions
//
//nolint:gocyclo // This function handles complex streaming logic
func (s *Swarm) StreamingResponse(
	ctx context.Context,
	agent *Agent,
	messages []llm.Message,
	contextVariables map[string]interface{},
	modelOverride string,
	handler StreamHandler,
	debug bool,
) error {
	if handler == nil {
		handler = &DefaultStreamHandler{}
	}

	if contextVariables == nil {
		contextVariables = make(map[string]interface{})
	}

	if debug {
		fmt.Printf("Debug: Using model: %s\n", agent.Model)
		fmt.Printf("Debug: Number of messages: %d\n", len(messages))
		fmt.Printf("Debug: Number of tools: %d\n", len(agent.Functions))
	}

	// Prepare the initial system message with agent instructions
	instructions := agent.Instructions
	if agent.InstructionsFunc != nil {
		instructions = agent.InstructionsFunc(contextVariables)
	}
	allMessages := append([]llm.Message{
		{
			Role:    llm.RoleSystem,
			Content: instructions,
		},
	}, messages...)

	// Build tool definitions
	var tools []llm.Tool
	for _, af := range agent.Functions {
		def := FunctionToDefinition(af)
		if debug {
			fmt.Printf("Debug: Adding tool: %s\n", def.Name)
		}
		tools = append(tools, llm.Tool{
			Type: "function",
			Function: &llm.Function{
				Name:        def.Name,
				Description: def.Description,
				Parameters:  def.Parameters,
			},
		})
	}

	// Prepare the streaming request
	model := agent.Model
	if modelOverride != "" {
		model = modelOverride
	}

	if debug {
		fmt.Printf("Debug: Final model: %s\n", model)
		fmt.Printf("Debug: Creating stream with %d messages\n", len(allMessages))
	}

	req := llm.ChatCompletionRequest{
		Model:    model,
		Messages: allMessages,
		Tools:    tools,
		Stream:   true,
	}

	stream, err := s.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		if debug {
			fmt.Printf("Debug: Stream creation error: %v\n", err)
		}
		handler.OnError(fmt.Errorf("failed to create chat completion stream: %v", err))
		return err
	}
	defer func() {
		if err := stream.Close(); err != nil {
			handler.OnError(fmt.Errorf("failed to close stream: %w", err))
		}
	}()

	handler.OnStart()

	var currentMessage llm.Message
	currentMessage.Role = llm.RoleAssistant
	currentMessage.Name = agent.Name

	// Track tool calls being built
	toolCallsInProgress := make(map[string]*llm.ToolCall)
	processedToolCalls := make(map[string]bool)

	// createNewStream creates a new stream and handles errors
	createNewStream := func() error {
		if err := stream.Close(); err != nil {
			handler.OnError(fmt.Errorf("failed to close stream: %v", err))
			return err
		}

		newStream, err := s.client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			if debug {
				fmt.Printf("Debug: Error creating new stream: %v\n", err)
			}
			handler.OnError(fmt.Errorf("failed to create new stream after tool call: %v", err))
			return err
		}
		stream = newStream
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			handler.OnError(ctx.Err())
			return ctx.Err()
		default:
			response, err := stream.Recv()
			if err != nil {
				if err.Error() == "EOF" {
					handler.OnComplete(currentMessage)
					return nil
				}
				if err.Error() == "stream closed" {
					// If stream is closed, try to create a new one
					if err := createNewStream(); err != nil {
						return err
					}
					continue
				}
				if debug {
					fmt.Printf("Debug: Error receiving from stream: %v\n", err)
				}
				handler.OnError(fmt.Errorf("error receiving from stream: %v", err))
				return err
			}

			if len(response.Choices) == 0 {
				continue
			}

			choice := response.Choices[0]

			// Handle content streaming
			if choice.Message.Content != "" {
				currentMessage.Content += choice.Message.Content
				handler.OnToken(choice.Message.Content)
			}

			// Handle tool calls
			if len(choice.Message.ToolCalls) > 0 {
				for _, toolCall := range choice.Message.ToolCalls {
					if debug {
						fmt.Printf("Debug: Processing tool call: ID=%s Name=%s\n",
							toolCall.ID, toolCall.Function.Name)
					}

					// Skip empty tool calls
					if toolCall.ID == "" {
						if debug {
							fmt.Printf("Debug: Skipping empty tool call ID\n")
						}
						continue
					}

					// Skip if we've already processed this tool call
					if processedToolCalls[toolCall.ID] {
						if debug {
							fmt.Printf("Debug: Skipping already processed tool call: %s\n", toolCall.ID)
						}
						continue
					}

					// Get or create the in-progress tool call
					inProgress, exists := toolCallsInProgress[toolCall.ID]
					if !exists {
						inProgress = &llm.ToolCall{
							ID:   toolCall.ID,
							Type: toolCall.Type,
							Function: llm.ToolCallFunction{
								Name:      toolCall.Function.Name,
								Arguments: "",
							},
						}
						toolCallsInProgress[toolCall.ID] = inProgress
						if debug {
							fmt.Printf("Debug: Created new tool call: %s, Name: %s\n",
								toolCall.ID, toolCall.Function.Name)
						}
					}

					// Update function name if provided
					if toolCall.Function.Name != "" && inProgress.Function.Name == "" {
						inProgress.Function.Name = toolCall.Function.Name
						if debug {
							fmt.Printf("Debug: Updated function name for tool call %s: %s\n",
								toolCall.ID, toolCall.Function.Name)
						}
					}

					// Accumulate function arguments
					if toolCall.Function.Arguments != "" {
						// Always append new arguments
						inProgress.Function.Arguments += toolCall.Function.Arguments
						if debug {
							fmt.Printf("Debug: Updated arguments for tool call %s: %s\n",
								toolCall.ID, inProgress.Function.Arguments)
						}

						// Try to parse the arguments to verify it's complete JSON
						var args map[string]interface{}
						if err := json.Unmarshal([]byte(inProgress.Function.Arguments), &args); err == nil {
							if debug {
								fmt.Printf("Debug: Valid JSON arguments for tool call %s: %v\n",
									toolCall.ID, args)
							}

							// Only execute if we haven't processed this tool call yet
							if !processedToolCalls[toolCall.ID] {
								// Find and execute the corresponding function
								var fn *AgentFunction
								for _, f := range agent.Functions {
									if f.Name == inProgress.Function.Name {
										fn = &f
										break
									}
								}

								if fn == nil {
									err := fmt.Errorf("unknown function: %s", inProgress.Function.Name)
									handler.OnError(err)
									continue
								}

								if debug {
									fmt.Printf("Debug: Executing function %s with args: %v\n",
										inProgress.Function.Name, args)
								}

								// Execute the function
								result := fn.Function(args, contextVariables)

								// Create function response message
								var resultContent string
								if result.Error != nil {
									resultContent = fmt.Sprintf("Error: %v", result.Error)
									if debug {
										fmt.Printf("Debug: Function execution error: %v\n", result.Error)
									}
								} else {
									resultContent = fmt.Sprintf("%v", result.Data)
									if debug {
										fmt.Printf("Debug: Function execution success: %v\n", result.Data)
									}
								}

								// Mark as processed and clean up
								processedToolCalls[toolCall.ID] = true
								delete(toolCallsInProgress, toolCall.ID)

								// Add to current message and notify handler
								currentMessage.ToolCalls = append(currentMessage.ToolCalls, *inProgress)
								handler.OnToolCall(*inProgress)

								// Add function response message
								functionMessage := llm.Message{
									Role:    llm.RoleFunction,
									Content: resultContent,
									Name:    inProgress.Function.Name,
								}

								// Add messages and create new stream
								allMessages = append(allMessages, currentMessage)
								allMessages = append(allMessages, functionMessage)
								req.Messages = allMessages

								if debug {
									fmt.Printf("Debug: Added function response message: %s = %s\n",
										functionMessage.Name, functionMessage.Content)
								}

								if err := createNewStream(); err != nil {
									handler.OnError(fmt.Errorf("failed to create new stream after tool call: %v", err))
									return err
								}

								if debug {
									fmt.Printf("Debug: Created new stream after tool call, messages count: %d\n", len(allMessages))
								}

								// Reset current message for new response
								currentMessage = llm.Message{
									Role: llm.RoleAssistant,
									Name: agent.Name,
								}
							}
						} else if debug {
							fmt.Printf("Debug: Incomplete JSON for tool call %s: %v\n", toolCall.ID, err)
						}
					}
				}
			}
		}
	}
}
