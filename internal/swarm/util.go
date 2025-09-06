package swarm

import (
	"fmt"
)

// ProcessAndPrintResponse processes and prints the response from the LLM.
// It uses different colors for different roles: blue for "assistant" and magenta for "function" or "tool".
func ProcessAndPrintResponse(response Response) {
	for _, message := range response.Messages {
		fmt.Printf("\033[90m%s\033[0m: %s\n", message.Role, message.Content)
		if message.Role == "assistant" {
			// Print assistant messages in blue, use agent name if available
			name := "Assistant"
			if response.Agent != nil && response.Agent.Name != "" {
				name = response.Agent.Name
			}

			// Print tool calls first
			if len(message.ToolCalls) > 0 {
				for _, toolCall := range message.ToolCalls {
					fmt.Printf("\033[94m%s\033[0m is calling function '%s' with arguments: %s\n",
						name, toolCall.Function.Name, toolCall.Function.Arguments)
				}
				continue // Skip printing empty content if we only have tool calls
			}

			// Print content if present
			if message.Content != "" {
				fmt.Printf("\033[94m%s\033[0m: %s\n", name, message.Content)
			}
		} else if message.Role == "function" || message.Role == "tool" {
			// Print function or tool results in magenta
			fmt.Printf("\033[95mFunction Result\033[0m: %s\n", message.Content)
		}
	}
}
