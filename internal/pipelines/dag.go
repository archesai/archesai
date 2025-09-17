package pipelines

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/tools"
)

// DAGNode represents a node in the workflow DAG.
type DAGNode struct {
	ID           uuid.UUID
	Step         *PipelineStep
	Dependencies []*DAGNode
	Dependents   []*DAGNode
	Status       NodeStatus
	Result       interface{}
	Error        error
}

// NodeStatus represents the execution status of a DAG node.
type NodeStatus int

// Node status constants.
const (
	NodeStatusPending   NodeStatus = iota // Node is waiting for dependencies
	NodeStatusReady                       // Node is ready to execute
	NodeStatusRunning                     // Node is currently executing
	NodeStatusCompleted                   // Node completed successfully
	NodeStatusFailed                      // Node execution failed
	NodeStatusSkipped                     // Node was skipped due to upstream failure
)

// DAG represents a Directed Acyclic Graph for workflow execution.
type DAG struct {
	Nodes     map[uuid.UUID]*DAGNode
	RootNodes []*DAGNode // Nodes with no dependencies
	mu        sync.RWMutex
}

// NewDAG creates a new DAG from pipeline steps.
func NewDAG(steps []PipelineStep, dependencies map[uuid.UUID][]uuid.UUID) (*DAG, error) {
	dag := &DAG{
		Nodes:     make(map[uuid.UUID]*DAGNode),
		RootNodes: []*DAGNode{},
	}

	// Create nodes
	for i := range steps {
		step := &steps[i]
		node := &DAGNode{
			ID:           step.ID,
			Step:         step,
			Dependencies: []*DAGNode{},
			Dependents:   []*DAGNode{},
			Status:       NodeStatusPending,
		}
		dag.Nodes[step.ID] = node
	}

	// Build edges
	for nodeID, depIDs := range dependencies {
		node, exists := dag.Nodes[nodeID]
		if !exists {
			return nil, fmt.Errorf("node %s not found", nodeID)
		}

		for _, depID := range depIDs {
			depNode, exists := dag.Nodes[depID]
			if !exists {
				return nil, fmt.Errorf("dependency %s not found for node %s", depID, nodeID)
			}

			node.Dependencies = append(node.Dependencies, depNode)
			depNode.Dependents = append(depNode.Dependents, node)
		}
	}

	// Identify root nodes
	for _, node := range dag.Nodes {
		if len(node.Dependencies) == 0 {
			dag.RootNodes = append(dag.RootNodes, node)
			node.Status = NodeStatusReady
		}
	}

	// Validate DAG (check for cycles)
	if err := dag.ValidateCycles(); err != nil {
		return nil, err
	}

	return dag, nil
}

// ValidateCycles checks if the DAG contains any cycles.
func (d *DAG) ValidateCycles() error {
	visited := make(map[uuid.UUID]bool)
	recStack := make(map[uuid.UUID]bool)

	for id := range d.Nodes {
		if !visited[id] {
			if d.hasCycleDFS(id, visited, recStack) {
				return fmt.Errorf("cycle detected in workflow DAG")
			}
		}
	}
	return nil
}

// hasCycleDFS is a helper function for cycle detection using DFS.
func (d *DAG) hasCycleDFS(nodeID uuid.UUID, visited, recStack map[uuid.UUID]bool) bool {
	visited[nodeID] = true
	recStack[nodeID] = true

	node := d.Nodes[nodeID]
	for _, dependent := range node.Dependents {
		if !visited[dependent.ID] {
			if d.hasCycleDFS(dependent.ID, visited, recStack) {
				return true
			}
		} else if recStack[dependent.ID] {
			return true
		}
	}

	recStack[nodeID] = false
	return false
}

// TopologicalSort returns nodes in topological order.
func (d *DAG) TopologicalSort() ([]*DAGNode, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Create a copy of in-degree counts
	inDegree := make(map[uuid.UUID]int)
	for id, node := range d.Nodes {
		inDegree[id] = len(node.Dependencies)
	}

	// Queue for nodes with no dependencies
	queue := make([]*DAGNode, 0, len(d.RootNodes))
	queue = append(queue, d.RootNodes...)

	result := make([]*DAGNode, 0, len(d.Nodes))

	for len(queue) > 0 {
		// Dequeue
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Decrease in-degree for dependent nodes
		for _, dependent := range current.Dependents {
			inDegree[dependent.ID]--
			if inDegree[dependent.ID] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	if len(result) != len(d.Nodes) {
		return nil, fmt.Errorf("cycle detected: topological sort incomplete")
	}

	return result, nil
}

// GetReadyNodes returns all nodes that are ready to execute.
func (d *DAG) GetReadyNodes() []*DAGNode {
	d.mu.RLock()
	defer d.mu.RUnlock()

	ready := []*DAGNode{}
	for _, node := range d.Nodes {
		if node.Status == NodeStatusReady {
			ready = append(ready, node)
		}
	}
	return ready
}

// MarkNodeCompleted marks a node as completed and updates dependent nodes.
func (d *DAG) MarkNodeCompleted(nodeID uuid.UUID, result interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()

	node, exists := d.Nodes[nodeID]
	if !exists {
		return
	}

	node.Status = NodeStatusCompleted
	node.Result = result

	// Check if dependent nodes are ready
	for _, dependent := range node.Dependents {
		if d.areAllDependenciesCompleted(dependent) {
			dependent.Status = NodeStatusReady
		}
	}
}

// MarkNodeFailed marks a node as failed and skips dependent nodes.
func (d *DAG) MarkNodeFailed(nodeID uuid.UUID, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	node, exists := d.Nodes[nodeID]
	if !exists {
		return
	}

	node.Status = NodeStatusFailed
	node.Error = err

	// Skip all dependent nodes
	d.skipDependents(node)
}

// skipDependents recursively skips all dependent nodes.
func (d *DAG) skipDependents(node *DAGNode) {
	for _, dependent := range node.Dependents {
		if dependent.Status != NodeStatusSkipped {
			dependent.Status = NodeStatusSkipped
			d.skipDependents(dependent)
		}
	}
}

// areAllDependenciesCompleted checks if all dependencies of a node are completed.
func (d *DAG) areAllDependenciesCompleted(node *DAGNode) bool {
	for _, dep := range node.Dependencies {
		if dep.Status != NodeStatusCompleted {
			return false
		}
	}
	return true
}

// GetExecutionPlan returns the execution plan as levels of parallel tasks.
func (d *DAG) GetExecutionPlan() ([][]uuid.UUID, error) {
	sorted, err := d.TopologicalSort()
	if err != nil {
		return nil, err
	}

	levels := [][]uuid.UUID{}
	processed := make(map[uuid.UUID]int)

	for _, node := range sorted {
		level := 0
		for _, dep := range node.Dependencies {
			if depLevel, ok := processed[dep.ID]; ok && depLevel >= level {
				level = depLevel + 1
			}
		}

		processed[node.ID] = level

		// Extend levels if needed
		for len(levels) <= level {
			levels = append(levels, []uuid.UUID{})
		}

		levels[level] = append(levels[level], node.ID)
	}

	return levels, nil
}

// DAGExecutor handles the execution of a DAG.
type DAGExecutor struct {
	dag         *DAG
	executor    ToolExecutor
	maxParallel int
}

// ToolExecutor defines the interface for executing tools.
type ToolExecutor interface {
	Execute(ctx context.Context, tool *tools.Tool, input interface{}) (interface{}, error)
}

// NewDAGExecutor creates a new DAG executor.
func NewDAGExecutor(dag *DAG, executor ToolExecutor, maxParallel int) *DAGExecutor {
	if maxParallel <= 0 {
		maxParallel = 1
	}
	return &DAGExecutor{
		dag:         dag,
		executor:    executor,
		maxParallel: maxParallel,
	}
}

// Execute runs the DAG execution.
func (e *DAGExecutor) Execute(ctx context.Context) error {
	// Get execution plan
	plan, err := e.dag.GetExecutionPlan()
	if err != nil {
		return fmt.Errorf("failed to get execution plan: %w", err)
	}

	// Execute level by level
	for levelIdx, level := range plan {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute nodes in parallel within each level
		var wg sync.WaitGroup
		semaphore := make(chan struct{}, e.maxParallel)
		errors := make(chan error, len(level))

		for _, nodeID := range level {
			node := e.dag.Nodes[nodeID]

			// Skip if node is not ready (due to previous failures)
			if node.Status == NodeStatusSkipped {
				continue
			}

			wg.Add(1)
			go func(n *DAGNode) {
				defer wg.Done()

				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				// Execute node
				if err := e.executeNode(ctx, n); err != nil {
					errors <- fmt.Errorf("node %s failed: %w", n.ID, err)
					e.dag.MarkNodeFailed(n.ID, err)
				} else {
					e.dag.MarkNodeCompleted(n.ID, n.Result)
				}
			}(node)
		}

		wg.Wait()
		close(errors)

		// Check for errors
		for err := range errors {
			if err != nil {
				return fmt.Errorf("level %d execution failed: %w", levelIdx, err)
			}
		}
	}

	return nil
}

// executeNode executes a single node.
func (e *DAGExecutor) executeNode(ctx context.Context, node *DAGNode) error {
	node.Status = NodeStatusRunning

	// Get input from dependencies
	input := e.collectInputs(node)

	// Execute the tool
	// Note: You'll need to implement the actual tool execution logic
	// This is a placeholder for the actual implementation
	result, err := e.executeStep(ctx, node.Step, input)
	if err != nil {
		return err
	}

	node.Result = result
	return nil
}

// collectInputs collects outputs from dependencies.
func (e *DAGExecutor) collectInputs(node *DAGNode) map[uuid.UUID]interface{} {
	inputs := make(map[uuid.UUID]interface{})
	for _, dep := range node.Dependencies {
		if dep.Status == NodeStatusCompleted {
			inputs[dep.ID] = dep.Result
		}
	}
	return inputs
}

// executeStep executes a pipeline step
// This is a placeholder - you'll need to implement actual tool execution.
func (e *DAGExecutor) executeStep(
	_ context.Context,
	step *PipelineStep,
	inputs map[uuid.UUID]interface{},
) (interface{}, error) {
	// TODO: Implement actual tool execution
	// This would typically:
	// 1. Load the tool configuration
	// 2. Prepare the input based on the tool's requirements
	// 3. Execute the tool (could be a container, function, API call, etc.)
	// 4. Return the output

	// For now, just return a placeholder result
	_ = inputs // will be used in actual implementation

	// Placeholder: in real implementation, this could return an error
	if step.ToolID == uuid.Nil {
		return nil, fmt.Errorf("invalid tool ID")
	}

	return map[string]interface{}{
		"stepID": step.ID,
		"toolID": step.ToolID,
		"status": "completed",
	}, nil
}
