package swarm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/archesai/archesai/internal/llm"
	"github.com/google/uuid"
)

// LangGraph inspired workflow system

// NodeID represents a unique identifier for a node in the workflow graph
type NodeID string

// EdgeType represents the type of relationship between nodes
type EdgeType string

const (
	// StandardEdge represents regular flow between nodes
	StandardEdge EdgeType = "standard"
	// ConditionEdge represents flow based on condition
	ConditionEdge EdgeType = "condition"
	// FallbackEdge is used when other edges fail
	FallbackEdge EdgeType = "fallback"
	// CallbackEdge is used for callbacks
	CallbackEdge EdgeType = "callback"
)

// StateKey represents a key in the state map
type StateKey string

// GraphState represents the current state of the workflow
type GraphState map[StateKey]interface{}

// Clone creates a deep copy of the state
func (s GraphState) Clone() GraphState {
	newState := make(GraphState)
	for k, v := range s {
		newState[k] = v
	}
	return newState
}

// UpdateState updates the state with new values
func (s GraphState) UpdateState(updates GraphState) {
	for k, v := range updates {
		s[k] = v
	}
}

// Get retrieves a value from state, with type assertion
func (s GraphState) Get(key StateKey) interface{} {
	return s[key]
}

// GetString gets a string value from state
func (s GraphState) GetString(key StateKey) (string, bool) {
	val, exists := s[key]
	if !exists {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

// GetBool gets a boolean value from state
func (s GraphState) GetBool(key StateKey) (bool, bool) {
	val, exists := s[key]
	if !exists {
		return false, false
	}
	b, ok := val.(bool)
	return b, ok
}

// MessageKey is the default key for storing messages in state
const MessageKey StateKey = "messages"

// NodeFunc is a function that processes state and returns updates
type NodeFunc func(ctx context.Context, state GraphState) (GraphState, error)

// ConditionFunc determines which edge to follow from a node
type ConditionFunc func(state GraphState) (NodeID, error)

// Node represents a node in the workflow graph
type Node struct {
	ID          NodeID
	Name        string
	Description string
	Process     NodeFunc
	Agent       *Agent // Optional agent associated with this node
	Metadata    map[string]interface{}
}

// Edge represents a connection between nodes
type Edge struct {
	From      NodeID
	To        NodeID
	Type      EdgeType
	Condition ConditionFunc // For conditional edges
	Metadata  map[string]interface{}
}

// Graph represents the workflow graph
type Graph struct {
	ID          string
	Name        string
	Description string
	Nodes       map[NodeID]*Node
	Edges       map[NodeID][]Edge
	EntryPoint  NodeID
	ExitPoints  []NodeID // Optional exit points
	mutex       sync.RWMutex
	eventHooks  map[string][]func(state GraphState)
}

// NewGraph creates a new workflow graph
func NewGraph(name string, description string) *Graph {
	return &Graph{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Nodes:       make(map[NodeID]*Node),
		Edges:       make(map[NodeID][]Edge),
		eventHooks:  make(map[string][]func(state GraphState)),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(id NodeID, name string, process NodeFunc) *Node {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	node := &Node{
		ID:       id,
		Name:     name,
		Process:  process,
		Metadata: make(map[string]interface{}),
	}

	g.Nodes[id] = node
	return node
}

// AddAgentNode adds a node with an associated agent
func (g *Graph) AddAgentNode(id NodeID, name string, agent *Agent) *Node {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Create node function that runs the agent
	processFunc := func(ctx context.Context, state GraphState) (GraphState, error) {
		// Get messages from state
		messagesRaw, ok := state[MessageKey]
		if !ok {
			messagesRaw = []llm.Message{}
		}

		// Convert to proper message type
		var messages []llm.Message
		messagesData, err := json.Marshal(messagesRaw)
		if err != nil {
			return state, fmt.Errorf("error marshaling messages: %w", err)
		}

		err = json.Unmarshal(messagesData, &messages)
		if err != nil {
			return state, fmt.Errorf("error unmarshaling messages: %w", err)
		}

		// Create swarm client if needed
		apiKey, _ := state.GetString("api_key")
		providerStr, _ := state.GetString("provider")
		provider := llm.Provider(providerStr)
		if provider == "" {
			provider = llm.OpenAI
		}

		client := NewSwarm(apiKey, provider)

		// Extract context variables
		contextVars := make(map[string]interface{})
		for k, v := range state {
			if strings.HasPrefix(string(k), "var_") {
				key := strings.TrimPrefix(string(k), "var_")
				contextVars[key] = v
			}
		}

		// Run the agent
		response, err := client.Run(ctx, agent, messages, contextVars, "", false, false, 1, true)
		if err != nil {
			return state, fmt.Errorf("error running agent: %w", err)
		}

		// Update state with new messages
		newMessages := append(messages, response.Messages...)

		// Create new state
		newState := state.Clone()
		newState[MessageKey] = newMessages

		// Add tool results to state if any
		if len(response.ToolResults) > 0 {
			toolResultsMap := make(map[string]interface{})
			for _, result := range response.ToolResults {
				toolResultsMap[result.ToolName] = result.Result.Data
			}
			newState["tool_results"] = toolResultsMap
		}

		// Update context variables in state
		for k, v := range response.ContextVariables {
			newState[StateKey("var_"+k)] = v
		}

		return newState, nil
	}

	node := &Node{
		ID:       id,
		Name:     name,
		Process:  processFunc,
		Agent:    agent,
		Metadata: make(map[string]interface{}),
	}

	g.Nodes[id] = node
	return node
}

// AddDirectedEdge adds a simple directed edge between nodes
func (g *Graph) AddDirectedEdge(from NodeID, to NodeID) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.Nodes[from]; !exists {
		return fmt.Errorf("source node %s does not exist", from)
	}

	if _, exists := g.Nodes[to]; !exists {
		return fmt.Errorf("destination node %s does not exist", to)
	}

	edge := Edge{
		From: from,
		To:   to,
		Type: StandardEdge,
	}

	g.Edges[from] = append(g.Edges[from], edge)
	return nil
}

// AddConditionalEdge adds an edge with a condition
func (g *Graph) AddConditionalEdge(from NodeID, to NodeID, condition ConditionFunc) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.Nodes[from]; !exists {
		return fmt.Errorf("source node %s does not exist", from)
	}

	if _, exists := g.Nodes[to]; !exists {
		return fmt.Errorf("destination node %s does not exist", to)
	}

	edge := Edge{
		From:      from,
		To:        to,
		Type:      ConditionEdge,
		Condition: condition,
	}

	g.Edges[from] = append(g.Edges[from], edge)
	return nil
}

// SetEntryPoint sets the entry point for the graph
func (g *Graph) SetEntryPoint(nodeID NodeID) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.Nodes[nodeID]; !exists {
		return fmt.Errorf("node %s does not exist", nodeID)
	}

	g.EntryPoint = nodeID
	return nil
}

// AddExitPoint adds an exit point to the graph
func (g *Graph) AddExitPoint(nodeID NodeID) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.Nodes[nodeID]; !exists {
		return fmt.Errorf("node %s does not exist", nodeID)
	}

	g.ExitPoints = append(g.ExitPoints, nodeID)
	return nil
}

// AddEventHook adds a hook for graph events
func (g *Graph) AddEventHook(event string, hook func(state GraphState)) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.eventHooks[event] = append(g.eventHooks[event], hook)
}

// fireEvent triggers event hooks
func (g *Graph) fireEvent(event string, state GraphState) {
	g.mutex.RLock()
	hooks, exists := g.eventHooks[event]
	g.mutex.RUnlock()

	if exists {
		for _, hook := range hooks {
			hook(state)
		}
	}
}

// ExecuteGraph runs the workflow graph from the entry point
func (g *Graph) ExecuteGraph(ctx context.Context, initialState GraphState) (GraphState, error) {
	if g.EntryPoint == "" {
		return initialState, errors.New("no entry point defined for graph")
	}

	currentNodeID := g.EntryPoint
	currentState := initialState
	visited := make(map[NodeID]int) // Track visited nodes to detect cycles

	// Start execution event
	g.fireEvent("graph_start", currentState)

	for {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return currentState, ctx.Err()
		default:
			// Continue execution
		}

		// Check for cycle
		visited[currentNodeID]++
		if visited[currentNodeID] > 10 { // Maximum cycle threshold
			return currentState, fmt.Errorf("potential infinite loop detected at node %s", currentNodeID)
		}

		// Get current node
		g.mutex.RLock()
		node, exists := g.Nodes[currentNodeID]
		g.mutex.RUnlock()

		if !exists {
			return currentState, fmt.Errorf("node %s not found", currentNodeID)
		}

		// Fire node entry event
		g.fireEvent(fmt.Sprintf("node_enter_%s", currentNodeID), currentState)

		// Execute node process
		newState, err := node.Process(ctx, currentState)
		if err != nil {
			g.fireEvent("node_error", currentState)
			return currentState, fmt.Errorf("error processing node %s: %w", currentNodeID, err)
		}

		currentState = newState

		// Fire node exit event
		g.fireEvent(fmt.Sprintf("node_exit_%s", currentNodeID), currentState)

		// Check if we've reached an exit point
		isExitPoint := false
		for _, exitPoint := range g.ExitPoints {
			if currentNodeID == exitPoint {
				isExitPoint = true
				break
			}
		}

		if isExitPoint {
			g.fireEvent("graph_complete", currentState)
			return currentState, nil
		}

		// Find next node
		g.mutex.RLock()
		edges, hasEdges := g.Edges[currentNodeID]
		g.mutex.RUnlock()

		if !hasEdges || len(edges) == 0 {
			return currentState, fmt.Errorf("node %s has no outgoing edges", currentNodeID)
		}

		// Determine next node based on edge types
		var nextNodeID NodeID
		found := false
		for _, edge := range edges {
			switch edge.Type {
			case StandardEdge:
				nextNodeID = edge.To
				found = true
			case ConditionEdge:
				if edge.Condition != nil {
					nodeID, err := edge.Condition(currentState)
					if err != nil {
						continue // Try next edge
					}
					nextNodeID = nodeID
					found = true
				}
			case FallbackEdge:
				// Only use fallback if no other edge matched
				if nextNodeID == "" {
					nextNodeID = edge.To
				}
			}
			if found {
				break
			}

			// If we found a next node, break out of the edge loop
			if nextNodeID != "" {
				break
			}
		}

		// If no valid edge was found, return error
		if nextNodeID == "" {
			return currentState, fmt.Errorf("no valid transition from node %s", currentNodeID)
		}

		// Update current node
		currentNodeID = nextNodeID
	}
}

// CreateAgentNode is a helper function to create common agent node types
func CreateAgentNode(g *Graph, id NodeID, name string, instructions string, model string, functions []AgentFunction, provider llm.Provider) *Node {
	agent := &Agent{
		Name:         name,
		Instructions: instructions,
		Model:        model,
		Functions:    functions,
		Provider:     provider,
	}

	return g.AddAgentNode(id, name, agent)
}

// CreateRouterNode creates a node that routes to different destinations based on content
func CreateRouterNode(g *Graph, id NodeID, destinations map[string]NodeID) *Node {
	routerFunc := func(_ context.Context, state GraphState) (GraphState, error) {
		// Router simply passes state through unchanged
		return state, nil
	}

	node := g.AddNode(id, fmt.Sprintf("Router-%s", id), routerFunc)

	// Create condition function for routing
	routeCondition := func(state GraphState) (NodeID, error) {
		// Get the last message
		messagesRaw, ok := state[MessageKey]
		if !ok {
			return "", errors.New("no messages in state")
		}

		// Convert to proper message type
		var messages []llm.Message
		messagesData, err := json.Marshal(messagesRaw)
		if err != nil {
			return "", fmt.Errorf("error marshaling messages: %w", err)
		}

		err = json.Unmarshal(messagesData, &messages)
		if err != nil {
			return "", fmt.Errorf("error unmarshaling messages: %w", err)
		}

		if len(messages) == 0 {
			return "", errors.New("no messages to route")
		}

		// Get the latest message content
		latestMsg := messages[len(messages)-1]
		content := strings.ToLower(latestMsg.Content)

		// Try to match with destinations
		for keyword, nodeID := range destinations {
			if strings.Contains(content, strings.ToLower(keyword)) {
				return nodeID, nil
			}
		}

		// Default destination (first one)
		for _, nodeID := range destinations {
			return nodeID, nil // Return the first one
		}

		return "", errors.New("no destination found")
	}

	// Add conditional edges to all destinations
	for _, destNodeID := range destinations {
		if err := g.AddConditionalEdge(id, destNodeID, routeCondition); err != nil {
			log.Printf("failed to add conditional edge from %s to %s: %v", id, destNodeID, err)
		}
	}

	return node
}

// CreateParallelNode creates a node that processes tasks in parallel
func CreateParallelNode(g *Graph, id NodeID, parallelProcesses []NodeFunc) *Node {
	parallelFunc := func(ctx context.Context, state GraphState) (GraphState, error) {
		var wg sync.WaitGroup
		results := make([]GraphState, len(parallelProcesses))
		errors := make([]error, len(parallelProcesses))

		// Create a copy of state for each parallel process
		for i, process := range parallelProcesses {
			wg.Add(1)
			go func(idx int, processFunc NodeFunc) {
				defer wg.Done()
				result, err := processFunc(ctx, state.Clone())
				results[idx] = result
				errors[idx] = err
			}(i, process)
		}

		// Wait for all processes to complete
		wg.Wait()

		// Check for errors
		for i, err := range errors {
			if err != nil {
				return state, fmt.Errorf("parallel process %d failed: %w", i, err)
			}
		}

		// Merge results
		mergedState := state.Clone()
		for _, result := range results {
			for k, v := range result {
				// Special handling for messages - combine them
				if k == MessageKey {
					// Combine message arrays
					existingMsgs, existingOk := mergedState[MessageKey].([]llm.Message)
					newMsgs, newOk := v.([]llm.Message)

					if existingOk && newOk {
						// Append new messages
						mergedState[MessageKey] = append(existingMsgs, newMsgs...)
					} else {
						// Just use the new messages
						mergedState[MessageKey] = v
					}
				} else {
					// For other keys, just overwrite
					mergedState[k] = v
				}
			}
		}

		return mergedState, nil
	}

	return g.AddNode(id, fmt.Sprintf("Parallel-%s", id), parallelFunc)
}

// CreateHumanInputNode creates a node that collects input from a human
func CreateHumanInputNode(g *Graph, id NodeID, prompt string) *Node {
	inputFunc := func(_ context.Context, state GraphState) (GraphState, error) {
		// Get current messages
		messagesRaw, ok := state[MessageKey]
		if !ok {
			messagesRaw = []llm.Message{}
		}

		// Convert to proper message type
		var messages []llm.Message
		messagesData, err := json.Marshal(messagesRaw)
		if err != nil {
			return state, fmt.Errorf("error marshaling messages: %w", err)
		}

		err = json.Unmarshal(messagesData, &messages)
		if err != nil {
			return state, fmt.Errorf("error unmarshaling messages: %w", err)
		}

		// Add system prompt for input
		messages = append(messages, llm.Message{
			Role:    llm.RoleAssistant,
			Content: prompt,
		})

		// Return updated state with the prompt
		newState := state.Clone()
		newState[MessageKey] = messages
		newState["waiting_for_input"] = true

		return newState, nil
	}

	return g.AddNode(id, fmt.Sprintf("HumanInput-%s", id), inputFunc)
}

// GraphBuilder provides a fluent interface for building graphs
type GraphBuilder struct {
	graph *Graph
}

// NewGraphBuilder creates a new graph builder
func NewGraphBuilder(name, description string) *GraphBuilder {
	return &GraphBuilder{
		graph: NewGraph(name, description),
	}
}

// Build returns the constructed graph
func (b *GraphBuilder) Build() *Graph {
	return b.graph
}

// WithAgent adds an agent node to the graph
func (b *GraphBuilder) WithAgent(id NodeID, name string, agent *Agent) *GraphBuilder {
	b.graph.AddAgentNode(id, name, agent)
	return b
}

// WithNode adds a generic node to the graph
func (b *GraphBuilder) WithNode(id NodeID, name string, process NodeFunc) *GraphBuilder {
	b.graph.AddNode(id, name, process)
	return b
}

// WithEdge adds a directed edge between nodes
func (b *GraphBuilder) WithEdge(from NodeID, to NodeID) *GraphBuilder {
	if err := b.graph.AddDirectedEdge(from, to); err != nil {
		log.Printf("failed to add directed edge from %s to %s: %v", from, to, err)
	}
	return b
}

// WithConditionalEdge adds a conditional edge between nodes
func (b *GraphBuilder) WithConditionalEdge(from NodeID, to NodeID, condition ConditionFunc) *GraphBuilder {
	if err := b.graph.AddConditionalEdge(from, to, condition); err != nil {
		log.Printf("failed to add conditional edge from %s to %s: %v", from, to, err)
	}
	return b
}

// WithEntryPoint sets the entry point for the graph
func (b *GraphBuilder) WithEntryPoint(nodeID NodeID) *GraphBuilder {
	if err := b.graph.SetEntryPoint(nodeID); err != nil {
		log.Printf("failed to set entry point %s: %v", nodeID, err)
	}
	return b
}

// WithExitPoint adds an exit point to the graph
func (b *GraphBuilder) WithExitPoint(nodeID NodeID) *GraphBuilder {
	if err := b.graph.AddExitPoint(nodeID); err != nil {
		log.Printf("failed to add exit point %s: %v", nodeID, err)
	}
	return b
}

// GraphRunner handles the execution of workflow graphs
type GraphRunner struct {
	graphs map[string]*Graph
	mu     sync.RWMutex
}

// NewGraphRunner creates a new graph runner
func NewGraphRunner() *GraphRunner {
	return &GraphRunner{
		graphs: make(map[string]*Graph),
	}
}

// RegisterGraph adds a graph to the runner
func (r *GraphRunner) RegisterGraph(graph *Graph) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.graphs[graph.ID] = graph
}

// ExecuteGraph runs a graph with the given initial state
func (r *GraphRunner) ExecuteGraph(ctx context.Context, graphID string, initialState GraphState) (GraphState, error) {
	r.mu.RLock()
	graph, exists := r.graphs[graphID]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("graph %s not found", graphID)
	}

	return graph.ExecuteGraph(ctx, initialState)
}
