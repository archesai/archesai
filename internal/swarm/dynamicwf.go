package swarm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/archesai/archesai/internal/llm"
)

// DynamicWorkflowCreator helps construct workflows dynamically based on user tasks
type DynamicWorkflowCreator struct {
	swarm        *Swarm
	baseAgents   map[string]*Agent // Pre-defined agent templates
	apiKey       string
	provider     llm.Provider
	plannerModel string
	taskAnalyzer *Agent
}

// NewDynamicWorkflowCreator creates a new workflow creator
func NewDynamicWorkflowCreator(apiKey string, provider llm.Provider) *DynamicWorkflowCreator {
	swarm := NewSwarm(apiKey, provider)

	// Create task analyzer agent with specialized prompt
	taskAnalyzer := &Agent{
		Name: "TaskAnalyzer",
		Instructions: `You are a specialized agent that analyzes user tasks and determines the optimal workflow structure.
Your job is to:
1. Identify the main goal of the user's task
2. Break down the task into logical sub-tasks
3. Determine which agent types would be needed (research, writing, analysis, coding, etc.)
4. Suggest a workflow structure (collaborative, hierarchical, or supervisor-based)
5. Define the agent relationships and data flow

Output your analysis in the following JSON format:
{
  "mainGoal": "Brief description of the overall goal",
  "workflowType": "collaborative|hierarchical|supervisor",
  "agents": [
    {
      "name": "AgentName",
      "role": "Brief description of agent role",
      "instructions": "Detailed instructions for this agent",
      "model": "Recommended model (e.g., gpt-4, gpt-3.5-turbo, etc.)",
      "connections": ["AgentName1", "AgentName2"]
    }
  ],
  "dataFlow": [
    {"from": "AgentName1", "to": "AgentName2", "description": "What data/results are passed"}
  ],
  "entryPoint": "Name of the agent that should start the workflow"
}`,
		Model: "gpt-4o",
	}

	return &DynamicWorkflowCreator{
		swarm:        swarm,
		baseAgents:   make(map[string]*Agent),
		apiKey:       apiKey,
		provider:     provider,
		plannerModel: "gpt-4o",
		taskAnalyzer: taskAnalyzer,
	}
}

// RegisterBaseAgent adds a pre-defined agent template
func (dwc *DynamicWorkflowCreator) RegisterBaseAgent(name string, agent *Agent) {
	dwc.baseAgents[name] = agent
}

// WorkflowSpec represents the specification for a dynamic workflow
type WorkflowSpec struct {
	MainGoal     string         `json:"mainGoal"`
	WorkflowType string         `json:"workflowType"`
	Agents       []AgentSpec    `json:"agents"`
	DataFlow     []DataFlowSpec `json:"dataFlow"`
	EntryPoint   string         `json:"entryPoint"`
}

// AgentSpec represents a specification for an agent
type AgentSpec struct {
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	Instructions string   `json:"instructions"`
	Model        string   `json:"model"`
	Connections  []string `json:"connections"`
}

// DataFlowSpec represents a data flow between agents
type DataFlowSpec struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Description string `json:"description"`
}

// CreateWorkflowFromTask generates a workflow specification based on a user task
func (dwc *DynamicWorkflowCreator) CreateWorkflowFromTask(ctx context.Context, userTask string) (*WorkflowSpec, error) {
	// Use the task analyzer to generate a workflow specification
	messages := []llm.Message{
		{Role: llm.RoleUser, Content: fmt.Sprintf("Analyze the following task and design an optimal workflow: %s", userTask)},
	}

	response, err := dwc.swarm.Run(ctx, dwc.taskAnalyzer, messages, nil, dwc.plannerModel, false, false, 1, false)
	if err != nil {
		return nil, fmt.Errorf("error analyzing task: %w", err)
	}

	// Extract and parse the JSON specification
	assistantMsg := response.Messages[len(response.Messages)-1].Content
	spec, err := extractJSONSpec(assistantMsg)
	if err != nil {
		return nil, fmt.Errorf("error extracting workflow specification: %w", err)
	}

	// Validate the workflow specification
	if err := validateWorkflowSpec(spec); err != nil {
		return nil, fmt.Errorf("invalid workflow specification: %w", err)
	}

	return spec, nil
}

// BuildWorkflow creates a concrete Workflow instance from a WorkflowSpec
func (dwc *DynamicWorkflowCreator) BuildWorkflow(spec *WorkflowSpec) (*Workflow, error) {
	// Determine workflow type
	var workflowType WorkflowType
	switch strings.ToLower(spec.WorkflowType) {
	case "collaborative":
		workflowType = CollaborativeWorkflow
	case "hierarchical":
		workflowType = HierarchicalWorkflow
	case string(SupervisorTeam):
		workflowType = SupervisorWorkflow
	default:
		return nil, fmt.Errorf("unknown workflow type: %s", spec.WorkflowType)
	}

	// Ensure we have a valid API key and provider
	if dwc.apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Create the workflow with properly initialized fields
	workflow := NewWorkflow(dwc.apiKey, dwc.provider, workflowType)

	// Create and add agents
	for _, agentSpec := range spec.Agents {
		// Check if we have a base template for this agent type
		baseAgent, hasTemplate := dwc.baseAgents[agentSpec.Name]

		agent := &Agent{
			Name:         agentSpec.Name,
			Instructions: agentSpec.Instructions,
			Model:        agentSpec.Model,
		}

		// If we have a template, inherit its functions and other properties
		if hasTemplate {
			agent.Functions = baseAgent.Functions
			// If model is not specified, use the one from the template
			if agent.Model == "" {
				agent.Model = baseAgent.Model
			}
		}

		workflow.AddAgent(agent)
	}

	// Create connections between agents
	for _, agentSpec := range spec.Agents {
		for _, connection := range agentSpec.Connections {
			if err := workflow.ConnectAgents(agentSpec.Name, connection); err != nil {
				return nil, fmt.Errorf("error creating connection from %s to %s: %w",
					agentSpec.Name, connection, err)
			}
		}
	}

	return workflow, nil
}

// CreateAndExecuteWorkflow is a convenience method to create and execute a workflow in one step
func (dwc *DynamicWorkflowCreator) CreateAndExecuteWorkflow(ctx context.Context, userTask string) (*WorkflowResult, error) {
	// Analyze task and create workflow spec
	spec, err := dwc.CreateWorkflowFromTask(ctx, userTask)
	if err != nil {
		return nil, err
	}

	// Build workflow from spec
	workflow, err := dwc.BuildWorkflow(spec)
	if err != nil {
		return nil, err
	}

	// Execute the workflow
	result, err := workflow.Execute(spec.EntryPoint, userTask)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Helper function to extract JSON from text
func extractJSONSpec(text string) (*WorkflowSpec, error) {
	var spec WorkflowSpec

	// Find JSON content between braces
	startIdx := strings.Index(text, "{")
	endIdx := strings.LastIndex(text, "}")

	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return nil, errors.New("could not find valid JSON in the response")
	}

	jsonContent := text[startIdx : endIdx+1]

	// Unmarshal JSON
	if err := json.Unmarshal([]byte(jsonContent), &spec); err != nil {
		return nil, fmt.Errorf("error parsing JSON specification: %w", err)
	}

	return &spec, nil
}

// Helper function to validate a workflow specification
func validateWorkflowSpec(spec *WorkflowSpec) error {
	if spec.MainGoal == "" {
		return errors.New("main goal is required")
	}

	if spec.WorkflowType == "" {
		return errors.New("workflow type is required")
	}

	if len(spec.Agents) == 0 {
		return errors.New("at least one agent is required")
	}

	// Check entry point
	if spec.EntryPoint == "" {
		return errors.New("entry point agent is required")
	}

	// Check that entry point agent exists
	entryPointExists := false
	agentMap := make(map[string]bool)

	for _, agent := range spec.Agents {
		if agent.Name == "" {
			return errors.New("agent name cannot be empty")
		}
		agentMap[agent.Name] = true

		if agent.Name == spec.EntryPoint {
			entryPointExists = true
		}
	}

	if !entryPointExists {
		return errors.New("entry point agent does not exist in the agent list")
	}

	// Validate connections
	for _, agent := range spec.Agents {
		for _, conn := range agent.Connections {
			if !agentMap[conn] {
				return fmt.Errorf("agent %s has connection to non-existent agent %s", agent.Name, conn)
			}
		}
	}

	return nil
}
