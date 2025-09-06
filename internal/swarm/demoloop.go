package swarm

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/archesai/archesai/internal/llm"
)

var (
	// ErrDemoLoopInterrupted is returned when the demo loop is interrupted by user input
	ErrDemoLoopInterrupted = errors.New("demo loop interrupted by user")
	// ErrAgentExecutionFailed is returned when agent execution fails
	ErrAgentExecutionFailed = errors.New("agent execution failed")
)

// DemoLoopConfig contains configuration options for the demo loop
type DemoLoopConfig struct {
	Timeout             time.Duration // Timeout for each agent execution
	MaxHistoryMessages  int           // Maximum number of messages to keep in history
	MaxInputLength      int           // Maximum length of user input
	ShowFunctionResults bool          // Whether to display function results
	ColorOutput         bool          // Whether to use color in output
	Debug               bool          // Whether to show debug information
	SaveHistory         bool          // Whether to save conversation history to file
	HistoryFile         string        // Path to file for saving history
}

// DefaultDemoLoopConfig returns default configuration for demo loop
func DefaultDemoLoopConfig() *DemoLoopConfig {
	return &DemoLoopConfig{
		Timeout:             60 * time.Second,
		MaxHistoryMessages:  50,
		MaxInputLength:      1000,
		ShowFunctionResults: true,
		ColorOutput:         true,
		Debug:               false,
		SaveHistory:         false,
		HistoryFile:         "swarm_conversation.json",
	}
}

// RunDemoLoop starts an interactive CLI loop for conversing with an agent
func RunDemoLoop(client *Swarm, agent *Agent) {
	RunDemoLoopWithConfig(client, agent, DefaultDemoLoopConfig())
}

// RunDemoLoopWithConfig starts an interactive CLI loop with custom configuration
//
//nolint:gocyclo // This function handles many interactive cases
func RunDemoLoopWithConfig(client *Swarm, agent *Agent, config *DemoLoopConfig) {
	// Validate inputs
	if client == nil {
		log.Fatal("Swarm client cannot be nil")
		return
	}

	if agent == nil {
		log.Fatal("Agent cannot be nil")
		return
	}

	if config == nil {
		config = DefaultDemoLoopConfig()
	}

	// Check if client is properly initialized
	if !client.IsInitialized() {
		log.Fatal("Swarm client is not properly initialized")
		return
	}

	// Create a new context for the operation
	baseCtx := context.Background()

	// Create a context that can be cancelled on SIGINT
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("\nReceived interrupt signal, shutting down...")
		cancel()
	}()

	// Print a starting message to the console
	printColoredText(config.ColorOutput, "\n==== Starting swarm CLI Demo ====\n", "cyan")
	fmt.Printf("Agent: %s\nModel: %s\n\n", agent.Name, agent.Model)
	if config.Debug {
		fmt.Printf("Debug mode: ON\n")
		if agent.Functions != nil {
			fmt.Printf("Available functions: %d\n", len(agent.Functions))
			for _, fn := range agent.Functions {
				fmt.Printf("- %s: %s\n", fn.Name, fn.Description)
			}
		}
		fmt.Println()
	}

	// Initialize a slice to store chat messages
	messages := []llm.Message{}

	// Create a new reader to read user input from the standard input
	reader := bufio.NewReader(os.Stdin)

	activeAgent := agent

	// Main interaction loop
	for {
		select {
		case <-ctx.Done():
			printColoredText(config.ColorOutput, "\nExiting demo loop...\n", "yellow")
			return
		default:
			// Prompt the user for input
			printColoredText(config.ColorOutput, "You: ", "green")

			// Read the user's input from the console
			userInput, err := reader.ReadString('\n')
			if err != nil {
				printColoredText(config.ColorOutput, fmt.Sprintf("Error reading input: %v\n", err), "red")
				continue
			}

			// Trim any leading or trailing whitespace from the input
			userInput = strings.TrimSpace(userInput)

			// Check for exit commands
			if userInput == "exit" || userInput == "quit" || userInput == "q" {
				printColoredText(config.ColorOutput, "Exiting demo loop...\n", "yellow")
				return
			}

			// Check if input is empty
			if userInput == "" {
				continue
			}

			// Check input length
			if config.MaxInputLength > 0 && len(userInput) > config.MaxInputLength {
				printColoredText(config.ColorOutput,
					fmt.Sprintf("Input too long (%d characters). Maximum is %d characters.\n",
						len(userInput), config.MaxInputLength), "red")
				continue
			}

			// Append the user's input as a new message to the messages slice
			messages = append(messages, llm.Message{
				Role:    llm.RoleUser,
				Content: userInput,
			})

			// Trim history if it exceeds the maximum
			if config.MaxHistoryMessages > 0 && len(messages) > config.MaxHistoryMessages {
				// Keep system message if it exists, plus recent messages
				var systemMsg *llm.Message
				for _, msg := range messages {
					if msg.Role == llm.RoleSystem {
						systemMsg = &msg
						break
					}
				}

				newMessages := make([]llm.Message, 0, config.MaxHistoryMessages+1)
				if systemMsg != nil {
					newMessages = append(newMessages, *systemMsg)
				}

				// Add recent messages
				startIdx := len(messages) - config.MaxHistoryMessages
				if startIdx < 0 {
					startIdx = 0
				}

				newMessages = append(newMessages, messages[startIdx:]...)
				messages = newMessages

				if config.Debug {
					fmt.Printf("History trimmed to %d messages\n", len(messages))
				}
			}

			// Create execution context with timeout
			execCtx, execCancel := context.WithTimeout(ctx, config.Timeout)

			// Show thinking indicator
			fmt.Print("Thinking...")

			// Execute agent
			startTime := time.Now()
			response, err := client.Run(execCtx, activeAgent, messages, nil, "", false, config.Debug, 5, true)
			execCancel() // Always cancel context

			// Clear thinking indicator
			fmt.Print("\r          \r")

			if err != nil {
				printColoredText(config.ColorOutput, fmt.Sprintf("Error: %v\n", err), "red")

				// Check for specific errors
				if errors.Is(err, context.DeadlineExceeded) {
					printColoredText(config.ColorOutput,
						fmt.Sprintf("Request timed out after %v\n", config.Timeout), "yellow")
				}

				continue
			}

			// Process the response and print it to the console
			var lastAssistantMessage llm.Message
			var functionMessages []llm.Message

			// Display response messages
			for _, msg := range response.Messages {
				switch msg.Role {
				case llm.RoleAssistant:
					if msg.Content != "" {
						name := response.Agent.Name
						if name == "" {
							name = "Assistant"
						}

						printColoredText(config.ColorOutput, fmt.Sprintf("%s: ", name), "blue")
						fmt.Println(msg.Content)
						lastAssistantMessage = msg
					}
				case llm.RoleFunction:
					if config.ShowFunctionResults {
						printColoredText(config.ColorOutput,
							fmt.Sprintf("%s function result: ", msg.Name), "magenta")
						fmt.Println(msg.Content)
					}
					functionMessages = append(functionMessages, msg)
				}
			}

			// Display timing information in debug mode
			if config.Debug {
				duration := time.Since(startTime)
				fmt.Printf("\n[Debug] Response time: %v\n", duration)

				if len(response.ToolResults) > 0 {
					fmt.Printf("[Debug] Tool calls: %d\n", len(response.ToolResults))
					for i, tool := range response.ToolResults {
						fmt.Printf("  %d. %s\n", i+1, tool.ToolName)
					}
				}
			}

			// Add function results to history first, then the assistant message
			messages = append(messages, functionMessages...)
			if lastAssistantMessage.Content != "" {
				messages = append(messages, lastAssistantMessage)
			}

			// Handle agent transfer
			if response.Agent != nil && response.Agent.Name != activeAgent.Name {
				printColoredText(config.ColorOutput,
					fmt.Sprintf("\nTransferring conversation to %s.\n\n", response.Agent.Name), "yellow")
				activeAgent = response.Agent
			}

			// Save conversation history if enabled
			if config.SaveHistory && config.HistoryFile != "" {
				saveConversationHistory(messages, config.HistoryFile)
			}
		}
	}
}

// Helper function to print colored text
func printColoredText(useColor bool, text, color string) {
	if !useColor {
		fmt.Print(text)
		return
	}

	var colorCode string
	switch strings.ToLower(color) {
	case "red":
		colorCode = "\033[91m"
	case "green":
		colorCode = "\033[92m"
	case "yellow":
		colorCode = "\033[93m"
	case "blue":
		colorCode = "\033[94m"
	case "magenta":
		colorCode = "\033[95m"
	case "cyan":
		colorCode = "\033[96m"
	case "gray", "grey":
		colorCode = "\033[90m"
	default:
		colorCode = "\033[0m"
	}

	fmt.Print(colorCode + text + "\033[0m")
}

// Helper function to save conversation history to a file
func saveConversationHistory(messages []llm.Message, filePath string) {
	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		log.Printf("Error serializing conversation history: %v", err)
		return
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Printf("Error saving conversation history: %v", err)
	}
}
