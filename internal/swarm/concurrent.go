package swarm

import (
	"context"
	"sync"

	"github.com/archesai/archesai/internal/llm"
)

// ConcurrentResult represents the result from a single agent's execution
type ConcurrentResult struct {
	AgentName string
	Response  Response
	Error     error
}

// ConcurrentSwarm manages concurrent execution of multiple agents
type ConcurrentSwarm struct {
	*Swarm
}

// NewConcurrentSwarm creates a new ConcurrentSwarm instance
func NewConcurrentSwarm(apiKey string, provider llm.Provider) *ConcurrentSwarm {
	return &ConcurrentSwarm{
		Swarm: NewSwarm(apiKey, provider),
	}
}

// AgentConfig holds the configuration for a single agent execution
type AgentConfig struct {
	Agent            *Agent
	Messages         []llm.Message
	ContextVariables map[string]interface{}
	ModelOverride    string
	Stream           bool
	Debug            bool
	MaxTurns         int
	ExecuteTools     bool
}

// RunConcurrent executes multiple agents concurrently and returns their results
func (cs *ConcurrentSwarm) RunConcurrent(ctx context.Context, configs map[string]AgentConfig) []ConcurrentResult {
	var (
		wg      sync.WaitGroup
		results = make([]ConcurrentResult, 0, len(configs))
		mu      sync.Mutex
	)

	// Create a channel for results to handle potential context cancellation
	resultChan := make(chan ConcurrentResult, len(configs))

	// Start each agent in its own goroutine
	for name, config := range configs {
		wg.Add(1)
		go func(name string, cfg AgentConfig) {
			defer wg.Done()

			resp, err := cs.Run(
				ctx,
				cfg.Agent,
				cfg.Messages,
				cfg.ContextVariables,
				cfg.ModelOverride,
				cfg.Stream,
				cfg.Debug,
				cfg.MaxTurns,
				cfg.ExecuteTools,
			)
			result := ConcurrentResult{
				AgentName: name,
				Response:  resp,
				Error:     err,
			}

			select {
			case <-ctx.Done():
				// Context was cancelled, but we still want to record the partial result
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			case resultChan <- result:
			}
		}(name, config)
	}

	// Wait for all goroutines to complete in a separate goroutine
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results, respecting context cancellation
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				// Channel closed, all agents have completed
				return results
			}
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		case <-ctx.Done():
			// Context cancelled, return partial results
			return results
		}
	}
}

// RunConcurrentOrdered executes multiple agents concurrently and returns their results in the order specified
func (cs *ConcurrentSwarm) RunConcurrentOrdered(ctx context.Context, orderedConfigs []struct {
	Name   string
	Config AgentConfig
}) []ConcurrentResult {
	configs := make(map[string]AgentConfig)
	for _, cfg := range orderedConfigs {
		configs[cfg.Name] = cfg.Config
	}

	results := cs.RunConcurrent(ctx, configs)

	// Reorder results to match input order
	ordered := make([]ConcurrentResult, len(results))
	resultMap := make(map[string]ConcurrentResult)
	for _, result := range results {
		resultMap[result.AgentName] = result
	}

	for i, cfg := range orderedConfigs {
		ordered[i] = resultMap[cfg.Name]
	}

	return ordered
}
