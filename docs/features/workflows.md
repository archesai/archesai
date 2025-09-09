# Workflows

## Overview

ArchesAI Workflows provide a powerful DAG-based (Directed Acyclic Graph) automation system for
orchestrating complex data processing pipelines with AI tool integration, enabling automated,
scalable, and intelligent data processing.

## Core Concepts

### Workflow Architecture

Workflows in ArchesAI are built on a flexible, extensible architecture that supports complex processing patterns:

```typescript
interface Workflow {
  id: string;
  name: string;
  description: string;
  organizationId: string;
  dag: DirectedAcyclicGraph;
  tools: Tool[];
  triggers: Trigger[];
  status: "draft" | "active" | "paused" | "archived";
  version: number;
  metadata: {
    tags: string[];
    category: string;
    estimatedDuration: number;
  };
  createdAt: Date;
  updatedAt: Date;
}
```

### Directed Acyclic Graph (DAG)

The DAG structure ensures proper execution order and prevents circular dependencies:

```typescript
interface DirectedAcyclicGraph {
  nodes: Node[];
  edges: Edge[];
  entryPoints: string[];
  exitPoints: string[];
}

interface Node {
  id: string;
  type: "tool" | "condition" | "parallel" | "loop" | "subworkflow";
  config: NodeConfig;
  inputs: Input[];
  outputs: Output[];
  retryPolicy?: RetryPolicy;
}
```

## Workflow Components

### 1. Tools

Tools are the building blocks of workflows, providing specific functionality:

```go
// Tool registration example
type Tool struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Category    string                 `json:"category"`
    Description string                 `json:"description"`
    InputSchema json.RawMessage        `json:"input_schema"`
    OutputSchema json.RawMessage       `json:"output_schema"`
    Execute     func(context.Context, map[string]interface{}) (map[string]interface{}, error)
}
```

**Built-in Tools:**

- **Data Processing**: JSON transform, CSV parser, XML processor
- **AI Integration**: LLM completion, embedding generation, classification
- **File Operations**: Read/write, format conversion, compression
- **HTTP Operations**: API calls, webhooks, data fetching
- **Database Operations**: Query, insert, update, delete
- **Notification**: Email, Slack, webhooks

### 2. Triggers

Triggers initiate workflow execution:

```yaml
triggers:
  - type: schedule
    config:
      cron: "0 9 * * MON-FRI"  # Every weekday at 9 AM
      timezone: "America/New_York"
  
  - type: webhook
    config:
      path: "/webhooks/document-upload"
      method: "POST"
      authentication: "bearer"
  
  - type: event
    config:
      source: "s3"
      event: "object.created"
      bucket: "documents"
  
  - type: manual
    config:
      requireApproval: true
      approvers: ["admin", "manager"]
```

### 3. Conditions

Conditional logic for dynamic workflow paths:

```yaml
nodes:
  - id: check_document_type
    type: condition
    config:
      expression: "input.fileType == 'pdf'"
      trueBranch: "process_pdf"
      falseBranch: "process_other"
```

### 4. Parallel Processing

Execute multiple branches simultaneously:

```yaml
nodes:
  - id: parallel_analysis
    type: parallel
    config:
      branches:
        - id: sentiment_analysis
          nodes: [analyze_sentiment, store_sentiment]
        - id: entity_extraction
          nodes: [extract_entities, validate_entities]
      joinStrategy: "wait_all"  # or "wait_any", "wait_n"
```

## Pipeline Creation

### Visual Pipeline Builder

The web interface provides an intuitive pipeline creation experience:

**Features:**

- **Drag-and-drop interface**: Visual node placement
- **Real-time validation**: Instant feedback on configuration errors
- **Node connection visualization**: See data flow clearly
- **Pipeline testing mode**: Test with sample data
- **Version control**: Track changes and rollback
- **Template library**: Start from pre-built templates

### YAML Definition

Define workflows programmatically:

```yaml
name: Document Processing Pipeline
description: Extract, analyze, and store document content
version: 1.0.0

metadata:
  tags: ["document", "nlp", "extraction"]
  category: "data-processing"
  estimatedDuration: 300  # seconds

# Input parameters
parameters:
  - name: document_url
    type: string
    required: true
    description: URL of the document to process
  
  - name: output_format
    type: string
    default: "json"
    enum: ["json", "xml", "csv"]

# Workflow nodes
nodes:
  - id: fetch_document
    type: tool
    tool: http_fetch
    config:
      url: "${parameters.document_url}"
      method: GET
    outputs:
      - name: document_content
        type: binary

  - id: extract_text
    type: tool
    tool: pdf_extractor
    inputs:
      - from: fetch_document.document_content
    config:
      ocr_enabled: true
      language: "en"
    outputs:
      - name: extracted_text
        type: string

  - id: analyze_content
    type: parallel
    config:
      branches:
        - id: sentiment
          nodes:
            - id: analyze_sentiment
              tool: sentiment_analyzer
              inputs:
                - from: extract_text.extracted_text
        
        - id: entities
          nodes:
            - id: extract_entities
              tool: ner_extractor
              inputs:
                - from: extract_text.extracted_text

  - id: generate_summary
    type: tool
    tool: llm_summarizer
    inputs:
      - from: extract_text.extracted_text
    config:
      model: "gpt-4"
      max_tokens: 500
      temperature: 0.3

  - id: store_results
    type: tool
    tool: database_insert
    inputs:
      - from: extract_text.extracted_text
      - from: analyze_sentiment.sentiment_score
      - from: extract_entities.entities
      - from: generate_summary.summary
    config:
      table: "processed_documents"
      mapping:
        text: "${extracted_text}"
        sentiment: "${sentiment_score}"
        entities: "${entities}"
        summary: "${summary}"

# Define execution flow
edges:
  - from: fetch_document
    to: extract_text
  - from: extract_text
    to: analyze_content
  - from: extract_text
    to: generate_summary
  - from: analyze_content
    to: store_results
  - from: generate_summary
    to: store_results

# Error handling
error_handling:
  default_strategy: "retry"
  retry_config:
    max_attempts: 3
    backoff: "exponential"
    initial_delay: 1000  # ms
  
  node_overrides:
    fetch_document:
      strategy: "fail_fast"
    store_results:
      strategy: "dead_letter_queue"
```

## Execution Engine

### Run Management

```go
// Workflow execution
type WorkflowRun struct {
    ID           uuid.UUID              `json:"id"`
    WorkflowID   uuid.UUID              `json:"workflow_id"`
    Status       RunStatus              `json:"status"`
    StartedAt    time.Time              `json:"started_at"`
    CompletedAt  *time.Time             `json:"completed_at"`
    Duration     time.Duration          `json:"duration"`
    NodeStates   map[string]NodeState   `json:"node_states"`
    Results      map[string]interface{} `json:"results"`
    Errors       []RunError             `json:"errors"`
}

type RunStatus string

const (
    RunStatusPending   RunStatus = "pending"
    RunStatusRunning   RunStatus = "running"
    RunStatusCompleted RunStatus = "completed"
    RunStatusFailed    RunStatus = "failed"
    RunStatusCancelled RunStatus = "cancelled"
)
```

### Execution Features

#### **Parallel Execution**

- Automatic detection of parallelizable nodes
- Worker pool management
- Resource allocation and limits

#### **State Management**

- Persistent state across node executions
- Checkpoint and resume capability
- Distributed state for scaled deployments

#### **Error Handling**

```yaml
error_handling:
  strategies:
    retry:
      max_attempts: 3
      backoff: exponential
      initial_delay: 1000ms
    
    circuit_breaker:
      threshold: 5
      timeout: 30s
      half_open_requests: 3
    
    fallback:
      handler: alternate_processing
      
    dead_letter_queue:
      queue: failed_workflows
      retention: 7d
```

## Monitoring & Observability

### Execution Metrics

```prometheus
# Workflow execution metrics
workflow_runs_total{workflow="document_processing", status="completed"} 1234
workflow_run_duration_seconds{workflow="document_processing", quantile="0.99"} 45.2
workflow_node_execution_time_seconds{node="extract_text", quantile="0.95"} 2.3
workflow_errors_total{workflow="document_processing", error_type="timeout"} 12
```

### Run History

```sql
-- Query run history
SELECT 
    wr.id,
    w.name as workflow_name,
    wr.status,
    wr.started_at,
    wr.duration,
    COUNT(CASE WHEN ns.status = 'failed' THEN 1 END) as failed_nodes
FROM workflow_runs wr
JOIN workflows w ON wr.workflow_id = w.id
JOIN node_states ns ON ns.run_id = wr.id
WHERE wr.started_at > NOW() - INTERVAL '24 hours'
GROUP BY wr.id, w.name
ORDER BY wr.started_at DESC;
```

### Real-time Monitoring

- **Live execution view**: See nodes executing in real-time
- **Performance metrics**: CPU, memory, execution time per node
- **Error tracking**: Immediate error notifications
- **Resource usage**: Track API calls, database queries

## Advanced Features

### 1. Subworkflows

Compose complex workflows from simpler ones:

```yaml
nodes:
  - id: process_batch
    type: subworkflow
    config:
      workflow_id: "batch_processor_v2"
      inputs:
        items: "${batch_items}"
      inherit_context: true
```

### 2. Dynamic Node Generation

Create nodes dynamically based on input:

```yaml
nodes:
  - id: dynamic_processor
    type: loop
    config:
      items: "${input.file_list}"
      iterator: "file"
      template:
        type: tool
        tool: file_processor
        config:
          path: "${file.path}"
          format: "${file.format}"
```

### 3. Custom Tools

Register custom tools for specialized processing:

```go
// Register custom tool
func RegisterCustomTool() error {
    tool := &Tool{
        ID:   "custom_analyzer",
        Name: "Custom Data Analyzer",
        Execute: func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
            // Custom processing logic
            data := input["data"].(string)
            result := analyzeData(data)
            return map[string]interface{}{
                "analysis": result,
                "timestamp": time.Now(),
            }, nil
        },
    }
    return toolRegistry.Register(tool)
}
```

### 4. Workflow Templates

Pre-built templates for common use cases:

- **ETL Pipeline**: Extract, transform, load data
- **Document Processing**: OCR, extraction, analysis
- **Data Enrichment**: Enhance data with external sources
- **ML Pipeline**: Training, evaluation, deployment
- **Alert System**: Monitor, evaluate, notify

## API Integration

### REST API

```bash
# Create workflow
POST /api/v1/workflows
Content-Type: application/json
{
  "name": "My Workflow",
  "dag": {...},
  "triggers": [...]
}

# Execute workflow
POST /api/v1/workflows/{id}/runs
{
  "parameters": {
    "document_url": "https://example.com/doc.pdf"
  }
}

# Get run status
GET /api/v1/workflows/{id}/runs/{run_id}

# List workflow runs
GET /api/v1/workflows/{id}/runs?status=completed&limit=10
```

### SDK Usage

```go
// Go SDK example
client := archesai.NewClient(apiKey)

// Create workflow
workflow, err := client.Workflows.Create(ctx, &archesai.WorkflowCreateRequest{
    Name: "Document Processor",
    DAG:  dag,
})

// Execute workflow
run, err := client.Workflows.Execute(ctx, workflow.ID, map[string]interface{}{
    "document_url": "https://example.com/document.pdf",
})

// Monitor execution
for run.Status == archesai.RunStatusRunning {
    run, err = client.Workflows.GetRun(ctx, workflow.ID, run.ID)
    time.Sleep(2 * time.Second)
}
```

## Best Practices

### Design Principles

1. **Keep nodes focused**: Single responsibility per node
2. **Handle errors gracefully**: Use appropriate error strategies
3. **Monitor performance**: Set up alerts for slow nodes
4. **Version workflows**: Track changes and enable rollback
5. **Test thoroughly**: Use test mode with sample data
6. **Document workflows**: Clear descriptions and parameter docs

### Performance Optimization

- **Batch processing**: Group similar operations
- **Caching**: Cache expensive computations
- **Async operations**: Use async tools when possible
- **Resource limits**: Set appropriate timeouts and memory limits
- **Parallel execution**: Maximize parallelism where possible

### Security Considerations

- **Input validation**: Validate all workflow inputs
- **Secret management**: Use secure credential storage
- **Access control**: Implement proper RBAC
- **Audit logging**: Track all workflow executions
- **Data encryption**: Encrypt sensitive data in transit and at rest

## Getting Started

1. **Define your workflow**: Start with YAML or visual builder
2. **Configure tools**: Set up required tool integrations
3. **Set triggers**: Define how workflows start
4. **Test execution**: Run with sample data
5. **Deploy**: Activate workflow for production
6. **Monitor**: Track execution and performance

For detailed examples and tutorials, see the
[Workflow Examples](https://github.com/archesai/archesai/tree/main/examples/workflows) repository.
