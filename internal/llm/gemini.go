package llm

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"strings"

// 	"google.golang.org/genai"
// )

// // GeminiLLM implements the LLM interface for Google's Gemini
// type GeminiLLM struct {
// 	client *genai.Client
// }

// // GeminiOptions contains configuration options for the Gemini model
// type GeminiOptions struct {
// 	Model          string
// 	HarmThreshold  genai.HarmBlockThreshold
// 	SafetySettings []*genai.SafetySetting
// }

// // NewGeminiLLM creates a new Gemini LLM client
// func NewGeminiLLM(apiKey string, opts ...GeminiOptions) (*GeminiLLM, error) {
// 	ctx := context.Background()
// 	client, err := genai.NewClient(ctx, &genai.ClientConfig{
// 		APIKey: apiKey,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
// 	}

// 	return &GeminiLLM{
// 		client: client,
// 	}, nil
// }

// // convertToGeminiMessages converts our generic Message type to Gemini's content type
// func convertToGeminiMessages(messages []Message) []genai.Part {
// 	var parts []genai.Part

// 	for _, msg := range messages {
// 		content := strings.TrimSpace(msg.Content)
// 		if content == "" {
// 			continue // Skip empty messages
// 		}

// 		switch msg.Role {
// 		case RoleSystem:
// 			parts = append(parts, genai.Text(fmt.Sprintf("[System]\n%s", content)))
// 		case "function":
// 			if msg.Name != "" {
// 				parts = append(parts, genai.Text(fmt.Sprintf("[Function: %s]\n%s", msg.Name, content)))
// 			} else {
// 				parts = append(parts, genai.Text(fmt.Sprintf("[Function Result]\n%s", content)))
// 			}
// 		case RoleAssistant:
// 			parts = append(parts, genai.Text(fmt.Sprintf("[Assistant]\n%s", content)))
// 		case RoleUser:
// 			parts = append(parts, genai.Text(fmt.Sprintf("[User]\n%s", content)))
// 		}
// 	}

// 	return parts
// }

// // convertFromGeminiResponse converts Gemini's response to our generic Message type
// func convertFromGeminiResponse(resp *genai.GenerateContentResponse) Message {
// 	if resp == nil || len(resp.Candidates) == 0 {
// 		return Message{Role: RoleAssistant, Content: ""}
// 	}

// 	content := resp.Candidates[0].Content
// 	if content == nil || len(content.Parts) == 0 {
// 		return Message{Role: RoleAssistant, Content: ""}
// 	}

// 	var textParts []string
// 	var toolCalls []ToolCall

// 	for _, part := range content.Parts {
// 		switch p := part.(type) {
// 		case genai.Text:
// 			text := string(p)
// 			// Remove role prefixes if present
// 			text = strings.TrimPrefix(text, "[Assistant]\n")
// 			text = strings.TrimSpace(text)
// 			if text != "" {
// 				textParts = append(textParts, text)
// 			}
// 		case genai.FunctionCall:
// 			args, err := json.Marshal(p.Args)
// 			if err == nil {
// 				toolCalls = append(toolCalls, ToolCall{
// 					Type: "function",
// 					Function: ToolCallFunction{
// 						Name:      p.Name,
// 						Arguments: string(args),
// 					},
// 				})
// 			}
// 		}
// 	}

// 	return Message{
// 		Role:      RoleAssistant,
// 		Content:   strings.Join(textParts, "\n"),
// 		ToolCalls: toolCalls,
// 	}
// }

// // convertToGeminiTools converts our generic Tool type to Gemini's tool type
// func convertToGeminiTools(tools []Tool) []*genai.Tool {
// 	if len(tools) == 0 {
// 		return nil
// 	}

// 	geminiTools := make([]*genai.Tool, len(tools))
// 	for i, tool := range tools {
// 		schema := &genai.Schema{
// 			Type: genai.TypeObject,
// 		}

// 		schema.Properties = make(map[string]*genai.Schema)

// 		if properties, ok := tool.Function.Parameters["properties"].(map[string]interface{}); ok {
// 			for name, prop := range properties {
// 				if propMap, ok := prop.(map[string]interface{}); ok {
// 					propSchema := &genai.Schema{}
// 					if typ, ok := propMap["type"].(string); ok {
// 						propSchema.Type = convertSchemaType(typ)
// 					}
// 					if desc, ok := propMap["description"].(string); ok {
// 						propSchema.Description = desc
// 					}
// 					schema.Properties[name] = propSchema
// 				}
// 			}
// 		}

// 		if required, ok := tool.Function.Parameters["required"].([]interface{}); ok {
// 			reqFields := make([]string, len(required))
// 			for i, r := range required {
// 				if str, ok := r.(string); ok {
// 					reqFields[i] = str
// 				}
// 			}
// 			schema.Required = reqFields
// 		}

// 		geminiTools[i] = &genai.Tool{
// 			FunctionDeclarations: []*genai.FunctionDeclaration{
// 				{
// 					Name:        tool.Function.Name,
// 					Description: tool.Function.Description,
// 					Parameters:  schema,
// 				},
// 			},
// 		}
// 	}
// 	return geminiTools
// }

// // convertSchemaType converts a JSON Schema type to Gemini schema type
// func convertSchemaType(typ string) genai.Type {
// 	switch typ {
// 	case "object":
// 		return genai.TypeObject
// 	case "string":
// 		return genai.TypeString
// 	case "number":
// 		return genai.TypeNumber
// 	case "integer":
// 		return genai.TypeInteger
// 	case "boolean":
// 		return genai.TypeBoolean
// 	case "array":
// 		return genai.TypeArray
// 	default:
// 		return genai.TypeUnspecified
// 	}
// }

// // convertFromGeminiToolCalls converts Gemini's tool calls to our generic type
// func convertFromGeminiToolCalls(parts []genai.Part) []ToolCall {
// 	var calls []ToolCall

// 	for _, part := range parts {
// 		if fc, ok := part.(genai.FunctionCall); ok {
// 			args, _ := json.Marshal(fc.Args)
// 			calls = append(calls, ToolCall{
// 				Type: "function",
// 				Function: ToolCallFunction{
// 					Name:      fc.Name,
// 					Arguments: string(args),
// 				},
// 			})
// 		}
// 	}

// 	return calls
// }

// // CreateChatCompletion implements the LLM interface for Gemini
// func (g *GeminiLLM) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
// 	// Create model and configure settings
// 	model := g.client.GenerativeModel(req.Model)

// 	if req.Temperature > 0 {
// 		model.SetTemperature(float32(req.Temperature))
// 	}
// 	if req.TopP > 0 {
// 		model.SetTopP(float32(req.TopP))
// 	}
// 	if req.MaxTokens > 0 {
// 		model.SetMaxOutputTokens(int32(req.MaxTokens))
// 	}

// 	// Check if we're in a function calling cycle
// 	inFunctionCall := false
// 	if len(req.Messages) > 0 {
// 		lastMsg := req.Messages[len(req.Messages)-1]
// 		inFunctionCall = lastMsg.Role == "function"
// 	}

// 	// Only set tools if we're not in a function calling cycle
// 	if len(req.Tools) > 0 && !inFunctionCall {
// 		model.Tools = convertToGeminiTools(req.Tools)
// 	}

// 	// Convert messages to Gemini format
// 	parts := convertToGeminiMessages(req.Messages)

// 	// Generate response
// 	resp, err := model.GenerateContent(ctx, parts...)
// 	if err != nil {
// 		return ChatCompletionResponse{}, fmt.Errorf("failed to generate content: %v", err)
// 	}

// 	// Convert response to our format
// 	choices := make([]Choice, len(resp.Candidates))
// 	for i, c := range resp.Candidates {
// 		msg := Message{
// 			Role:    RoleAssistant,
// 			Content: "",
// 		}

// 		// Handle function calls and text separately
// 		var textParts []string
// 		for _, part := range c.Content.Parts {
// 			switch p := part.(type) {
// 			case genai.Text:
// 				textParts = append(textParts, string(p))
// 			case genai.FunctionCall:
// 				if !inFunctionCall {
// 					args, err := json.Marshal(p.Args)
// 					if err != nil {
// 						continue
// 					}
// 					msg.ToolCalls = append(msg.ToolCalls, ToolCall{
// 						Type: "function",
// 						Function: ToolCallFunction{
// 							Name:      p.Name,
// 							Arguments: string(args),
// 						},
// 					})
// 				}
// 			}
// 		}
// 		msg.Content = strings.Join(textParts, "")

// 		// If we have function calls and we're not in a function calling cycle
// 		if len(msg.ToolCalls) > 0 && !inFunctionCall {
// 			// Create new messages array with original messages plus function results
// 			var newMessages []Message
// 			newMessages = append(newMessages, req.Messages...)

// 			// Add assistant's message with function calls
// 			newMessages = append(newMessages, Message{
// 				Role:      RoleAssistant,
// 				Content:   msg.Content,
// 				ToolCalls: msg.ToolCalls,
// 			})

// 			// Add function results
// 			for _, toolCall := range msg.ToolCalls {
// 				newMessages = append(newMessages, Message{
// 					Role:    "function",
// 					Content: toolCall.Function.Arguments,
// 					Name:    toolCall.Function.Name,
// 				})
// 			}

// 			// Create a new request for the final response
// 			finalReq := ChatCompletionRequest{
// 				Model:       req.Model,
// 				Messages:    newMessages,
// 				Temperature: req.Temperature,
// 				TopP:        req.TopP,
// 				MaxTokens:   req.MaxTokens,
// 			}

// 			// Create a new model instance for the final response
// 			finalModel := g.client.GenerativeModel(finalReq.Model)
// 			if finalReq.Temperature > 0 {
// 				finalModel.SetTemperature(float32(finalReq.Temperature))
// 			}
// 			if finalReq.TopP > 0 {
// 				finalModel.SetTopP(float32(finalReq.TopP))
// 			}
// 			if finalReq.MaxTokens > 0 {
// 				finalModel.SetMaxOutputTokens(int32(finalReq.MaxTokens))
// 			}

// 			// Convert messages to Gemini format
// 			finalParts := convertToGeminiMessages(newMessages)

// 			// Get the final response
// 			finalResp, err := finalModel.GenerateContent(ctx, finalParts...)
// 			if err != nil {
// 				return ChatCompletionResponse{}, fmt.Errorf("failed to generate final response: %v", err)
// 			}

// 			if len(finalResp.Candidates) > 0 {
// 				var finalTextParts []string
// 				for _, part := range finalResp.Candidates[0].Content.Parts {
// 					if t, ok := part.(genai.Text); ok {
// 						finalTextParts = append(finalTextParts, string(t))
// 					}
// 				}
// 				msg.Content = strings.Join(finalTextParts, "")
// 			}
// 		}

// 		choices[i] = Choice{
// 			Index:        i,
// 			Message:      msg,
// 			FinishReason: string(c.FinishReason),
// 		}
// 	}

// 	// Build response with usage metrics if available
// 	response := ChatCompletionResponse{
// 		Choices: choices,
// 	}
// 	if resp.UsageMetadata != nil {
// 		response.Usage = Usage{
// 			PromptTokens:     int(resp.UsageMetadata.PromptTokenCount),
// 			CompletionTokens: int(resp.UsageMetadata.CandidatesTokenCount),
// 			TotalTokens:      int(resp.UsageMetadata.TotalTokenCount),
// 		}
// 	}

// 	return response, nil
// }

// // geminiStreamWrapper wraps Gemini's stream to implement our ChatCompletionStream interface
// type geminiStreamWrapper struct {
// 	iter            *genai.GenerateContentResponseIterator
// 	client          *GeminiLLM
// 	req             ChatCompletionRequest
// 	ctx             context.Context
// 	inFunctionCall  bool
// 	currentToolCall *ToolCall
// 	toolCallBuffer  map[string]*ToolCall
// }

// func (w *geminiStreamWrapper) Recv() (ChatCompletionResponse, error) {
// 	// Get next response from iterator
// 	resp, err := w.iter.Next()
// 	if err != nil {
// 		if w.inFunctionCall {
// 			// If we're in a function call and hit the end, start a new stream for the response
// 			w.inFunctionCall = false
// 			w.iter = nil
// 			return w.handleFunctionResponse()
// 		}
// 		return ChatCompletionResponse{}, err
// 	}

// 	choices := make([]Choice, len(resp.Candidates))
// 	for i, c := range resp.Candidates {
// 		msg := Message{
// 			Role:    RoleAssistant,
// 			Content: "",
// 		}

// 		// Handle function calls and text separately
// 		var textParts []string
// 		for _, part := range c.Content.Parts {
// 			switch p := part.(type) {
// 			case genai.Text:
// 				textParts = append(textParts, string(p))
// 			case genai.FunctionCall:
// 				if !w.inFunctionCall {
// 					args, err := json.Marshal(p.Args)
// 					if err != nil {
// 						continue
// 					}

// 					// Create or update tool call
// 					toolCall, exists := w.toolCallBuffer[p.Name]
// 					if !exists {
// 						toolCall = &ToolCall{
// 							ID:   p.Name, // Use function name as ID since Gemini doesn't provide one
// 							Type: "function",
// 							Function: ToolCallFunction{
// 								Name:      p.Name,
// 								Arguments: "",
// 							},
// 						}
// 						w.toolCallBuffer[p.Name] = toolCall
// 						w.currentToolCall = toolCall
// 					}

// 					// Update arguments
// 					toolCall.Function.Arguments += string(args)

// 					// Try to parse the arguments to verify if it's complete JSON
// 					var parsedArgs map[string]interface{}
// 					if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &parsedArgs); err == nil {
// 						// Add to message's tool calls and remove from buffer
// 						msg.ToolCalls = append(msg.ToolCalls, *toolCall)
// 						delete(w.toolCallBuffer, p.Name)
// 						w.currentToolCall = nil
// 						w.inFunctionCall = true
// 					}
// 				}
// 			}
// 		}
// 		msg.Content = strings.Join(textParts, "")

// 		choices[i] = Choice{
// 			Index:        i,
// 			Message:      msg,
// 			FinishReason: string(c.FinishReason),
// 		}
// 	}

// 	return ChatCompletionResponse{
// 		Choices: choices,
// 	}, nil
// }

// func (w *geminiStreamWrapper) handleFunctionResponse() (ChatCompletionResponse, error) {
// 	// Create new messages array with original messages plus function results
// 	var newMessages []Message
// 	newMessages = append(newMessages, w.req.Messages...)

// 	// Add the last assistant message with function calls
// 	lastAssistantMsg := w.req.Messages[len(w.req.Messages)-1]
// 	newMessages = append(newMessages, lastAssistantMsg)

// 	// Add function results from tool calls
// 	for _, toolCall := range lastAssistantMsg.ToolCalls {
// 		newMessages = append(newMessages, Message{
// 			Role:    "function",
// 			Content: toolCall.Function.Arguments,
// 			Name:    toolCall.Function.Name,
// 		})
// 	}

// 	// Create a new model instance for the final response
// 	model := w.client.client.GenerativeModel(w.req.Model)
// 	if w.req.Temperature > 0 {
// 		model.SetTemperature(float32(w.req.Temperature))
// 	}
// 	if w.req.TopP > 0 {
// 		model.SetTopP(float32(w.req.TopP))
// 	}
// 	if w.req.MaxTokens > 0 {
// 		model.SetMaxOutputTokens(int32(w.req.MaxTokens))
// 	}

// 	// Convert messages to Gemini format
// 	parts := convertToGeminiMessages(newMessages)

// 	// Start new stream for the response
// 	w.iter = model.GenerateContentStream(w.ctx, parts...)

// 	// Get first response from new stream
// 	return w.Recv()
// }

// func (w *geminiStreamWrapper) Close() error {
// 	w.iter = nil
// 	return nil
// }

// // CreateChatCompletionStream implements the LLM interface for Gemini streaming
// func (g *GeminiLLM) CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (ChatCompletionStream, error) {
// 	// Create model and configure settings
// 	model := g.client.GenerativeModel(req.Model)

// 	if req.Temperature > 0 {
// 		model.SetTemperature(float32(req.Temperature))
// 	}
// 	if req.TopP > 0 {
// 		model.SetTopP(float32(req.TopP))
// 	}
// 	if req.MaxTokens > 0 {
// 		model.SetMaxOutputTokens(int32(req.MaxTokens))
// 	}

// 	// Check if we're in a function calling cycle
// 	inFunctionCall := false
// 	if len(req.Messages) > 0 {
// 		lastMsg := req.Messages[len(req.Messages)-1]
// 		inFunctionCall = lastMsg.Role == "function" || (lastMsg.Role == RoleAssistant && len(lastMsg.ToolCalls) > 0)
// 	}

// 	// Only set tools if we're not in a function calling cycle
// 	if len(req.Tools) > 0 && !inFunctionCall {
// 		model.Tools = convertToGeminiTools(req.Tools)
// 	}

// 	// Convert messages to Gemini format
// 	parts := convertToGeminiMessages(req.Messages)

// 	// Generate streaming response
// 	iter := model.GenerateContentStream(ctx, parts...)

// 	return &geminiStreamWrapper{
// 		iter:           iter,
// 		client:         g,
// 		req:            req,
// 		ctx:            ctx,
// 		inFunctionCall: inFunctionCall,
// 		toolCallBuffer: make(map[string]*ToolCall),
// 	}, nil
// }
