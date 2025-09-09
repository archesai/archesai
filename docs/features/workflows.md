# Workflows

## Overview

ArchesAI Workflows provide a powerful DAG-based (Directed Acyclic Graph) automation system for
orchestrating complex data processing pipelines with AI tool integration.

## Core Concepts

### Workflow Architecture

```typescript
interface Workflow {
  id: string;
  name: string;
  description: string;
  dag: DirectedAcyclicGraph;
  tools: Tool[];
  triggers: Trigger[];
  status: "draft" | "active" | "paused";
  createdAt: Date;
  updatedAt: Date;
}
```

### DAG Execution

- **Parallel Processing**: Execute independent nodes simultaneously
- **Dependency Management**: Automatic ordering based on dependencies
- **Error Handling**: Retry logic and fallback paths
- **State Management**: Persistent workflow state across executions

## Pipeline Creation

### Visual Pipeline Builder

- Drag-and-drop interface
- Real-time validation
- Node connection visualization
- Pipeline testing mode

### YAML Definition

```yaml
name: Document Processing Pipeline
description: Extract and analyze document content
nodes:
  - id: extract_text
    type: tool
    tool: pdf_extractor
    inputs:
      file: "${workflow.input.document}"

  - id: analyze_sentiment
    type: tool
    tool: sentiment_analyzer
    dependencies: [extract_text]
    inputs:
      text: "${extract_text.output.content}"

  - id: generate_summary
    type: tool
    tool: ai_summarizer
    dependencies: [extract_text]
    inputs:
      text: "${extract_text.output.content}"
      max_length: 500

  - id: store_results
    type: tool
    tool: database_writer
    dependencies: [analyze_sentiment, generate_summary]
    inputs:
      sentiment: "${analyze_sentiment.output}"
      summary: "${generate_summary.output}"
```

### Programmatic Creation

```go
workflow := &Workflow{
    Name: "Data Processing Pipeline",
    DAG: &DAG{
        Nodes: []Node{
            {
                ID:   "fetch_data",
                Type: "tool",
                Tool: "http_fetcher",
                Config: map[string]interface{}{
                    "url": "https://api.example.com/data",
                },
            },
            {
                ID:           "process_data",
                Type:         "tool",
                Tool:         "data_processor",
                Dependencies: []string{"fetch_data"},
            },
        },
    },
}
```

## Tool Registry System

### Built-in Tools

#### Data Processing

- CSV Parser
- JSON Transformer
- XML Processor
- Excel Reader
- PDF Extractor

#### AI/ML Tools

- Text Summarizer
- Sentiment Analyzer
- Language Translator
- Image Classifier
- Entity Extractor

#### Integration Tools

- HTTP Client
- Database Connector
- S3 Uploader
- Email Sender
- Webhook Caller

#### Transformation Tools

- Data Mapper
- Format Converter
- Schema Validator
- Data Enricher
- Deduplicator

### Custom Tool Development

```go
type Tool interface {
    // Metadata
    GetID() string
    GetName() string
    GetDescription() string
    GetVersion() string

    // Schema
    GetInputSchema() *Schema
    GetOutputSchema() *Schema
    GetConfigSchema() *Schema

    // Execution
    Execute(ctx context.Context, input Input) (Output, error)
    Validate(input Input) error
}

// Example custom tool
type CustomProcessor struct {
    id   string
    name string
}

func (t *CustomProcessor) Execute(ctx context.Context, input Input) (Output, error) {
    // Process input data
    data := input.Get("data")

    // Perform transformation
    result := transform(data)

    // Return output
    return Output{
        "processed": result,
        "timestamp": time.Now(),
    }, nil
}
```

### Tool Registration

```go
// Register custom tool
registry.Register(&CustomProcessor{
    id:   "custom_processor",
    name: "Custom Data Processor",
})

// Use in workflow
workflow.AddNode(&Node{
    Type: "tool",
    Tool: "custom_processor",
})
```

## Run Management

### Execution Modes

#### Manual Execution

```bash
POST /api/v1/workflows/:id/execute
{
  "inputs": {
    "source": "manual",
    "data": {...}
  }
}
```

#### Scheduled Execution

```yaml
triggers:
  - type: schedule
    cron: "0 9 * * MON-FRI" # Every weekday at 9 AM
    timezone: "America/New_York"
```

#### Event-Driven Execution

```yaml
triggers:
  - type: webhook
    endpoint: /webhooks/workflow/:id
    secret: "${WEBHOOK_SECRET}"

  - type: file_upload
    bucket: "input-documents"
    pattern: "*.pdf"
```

### Execution Monitoring

```typescript
interface WorkflowRun {
  id: string;
  workflowId: string;
  status: "pending" | "running" | "completed" | "failed";
  startedAt: Date;
  completedAt?: Date;
  nodes: NodeExecution[];
  outputs: Record<string, any>;
  errors?: Error[];
}

interface NodeExecution {
  nodeId: string;
  status: "pending" | "running" | "completed" | "failed" | "skipped";
  startedAt?: Date;
  completedAt?: Date;
  inputs: Record<string, any>;
  outputs?: Record<string, any>;
  error?: Error;
  retries: number;
}
```

### Run History and Analytics

```sql
-- Get run statistics
SELECT
    workflow_id,
    COUNT(*) as total_runs,
    AVG(EXTRACT(EPOCH FROM (completed_at - started_at))) as avg_duration,
    SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as successful_runs,
    SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed_runs
FROM workflow_runs
WHERE started_at >= NOW() - INTERVAL '30 days'
GROUP BY workflow_id;
```

## Advanced Features

### Conditional Logic

```yaml
nodes:
  - id: check_condition
    type: condition
    expression: "${previous.output.value > 100}"

  - id: high_value_path
    type: tool
    tool: premium_processor
    dependencies: [check_condition]
    condition: "${check_condition.result == true}"

  - id: normal_path
    type: tool
    tool: standard_processor
    dependencies: [check_condition]
    condition: "${check_condition.result == false}"
```

### Loops and Iteration

```yaml
nodes:
  - id: process_items
    type: foreach
    items: "${input.items}"
    iterator: item
    nodes:
      - id: process_single
        type: tool
        tool: item_processor
        inputs:
          data: "${item}"
```

### Error Handling

```yaml
nodes:
  - id: risky_operation
    type: tool
    tool: external_api
    retry:
      attempts: 3
      backoff: exponential
      initial_delay: 1s
    on_error:
      - type: fallback
        node: backup_processor
      - type: notify
        channel: slack
        message: "Operation failed after retries"
```

### Parallel Processing

```yaml
nodes:
  - id: split_data
    type: tool
    tool: data_splitter
    outputs:
      chunks: array

  - id: process_chunks
    type: parallel
    dependencies: [split_data]
    max_concurrency: 10
    items: "${split_data.outputs.chunks}"
    node:
      type: tool
      tool: chunk_processor

  - id: merge_results
    type: tool
    tool: data_merger
    dependencies: [process_chunks]
```

## API Endpoints

### Workflow Management

- `GET /api/v1/workflows` - List workflows
- `POST /api/v1/workflows` - Create workflow
- `GET /api/v1/workflows/:id` - Get workflow details
- `PUT /api/v1/workflows/:id` - Update workflow
- `DELETE /api/v1/workflows/:id` - Delete workflow
- `POST /api/v1/workflows/:id/clone` - Clone workflow

### Execution

- `POST /api/v1/workflows/:id/execute` - Execute workflow
- `GET /api/v1/workflows/:id/runs` - List workflow runs
- `GET /api/v1/runs/:id` - Get run details
- `POST /api/v1/runs/:id/cancel` - Cancel running workflow
- `GET /api/v1/runs/:id/logs` - Get execution logs

### Tools

- `GET /api/v1/tools` - List available tools
- `GET /api/v1/tools/:id` - Get tool details
- `POST /api/v1/tools` - Register custom tool
- `PUT /api/v1/tools/:id` - Update tool
- `DELETE /api/v1/tools/:id` - Unregister tool

## Performance Optimization

### Caching Strategies

- Node output caching
- Tool result memoization
- Dependency graph caching
- Connection pooling

### Resource Management

```yaml
resources:
  cpu: 2
  memory: 4Gi
  timeout: 30m
  max_retries: 3
  queue_priority: high
```

### Scaling

- Horizontal workflow execution
- Distributed node processing
- Queue-based load balancing
- Auto-scaling based on metrics

## Monitoring and Observability

### Metrics

- Workflow execution time
- Node processing duration
- Success/failure rates
- Resource utilization
- Queue depths

### Logging

```json
{
  "workflow_id": "wf_123",
  "run_id": "run_456",
  "node_id": "process_data",
  "level": "info",
  "message": "Processing 1000 records",
  "timestamp": "2024-01-15T10:30:00Z",
  "metadata": {
    "record_count": 1000,
    "processing_time_ms": 250
  }
}
```

### Alerting

- Workflow failure notifications
- SLA breach alerts
- Resource exhaustion warnings
- Long-running workflow alerts

## Testing

### Unit Testing

```go
func TestWorkflowExecution(t *testing.T) {
    workflow := CreateTestWorkflow()
    executor := NewExecutor()

    result, err := executor.Execute(workflow, TestInputs())

    assert.NoError(t, err)
    assert.Equal(t, "completed", result.Status)
    assert.Contains(t, result.Outputs, "processed_data")
}
```

### Integration Testing

- End-to-end workflow testing
- Tool integration verification
- Trigger testing
- Performance benchmarking

## Best Practices

### Design Patterns

- Keep workflows focused and single-purpose
- Use sub-workflows for reusable components
- Implement proper error handling
- Version control workflow definitions

### Security

- Validate all inputs
- Use secrets management for credentials
- Implement rate limiting
- Audit workflow executions

### Performance

- Minimize node dependencies
- Use caching appropriately
- Batch process where possible
- Monitor resource usage

## Troubleshooting

### Common Issues

#### Workflow Stuck in Running

- Check for deadlocks in dependencies
- Verify external service availability
- Review timeout settings
- Check resource limits

#### High Failure Rate

- Review error logs
- Check input validation
- Verify tool configurations
- Monitor external dependencies

#### Performance Degradation

- Analyze execution metrics
- Check database query performance
- Review caching effectiveness
- Consider workflow optimization

## Migration Guide

### Importing Workflows

1. Export workflow as YAML/JSON
2. Validate schema compatibility
3. Update tool references
4. Test in staging environment
5. Deploy to production

### Version Migration

- Backup existing workflows
- Run migration scripts
- Update tool versions
- Test backward compatibility
- Update documentation
