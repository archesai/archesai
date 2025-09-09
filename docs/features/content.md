# Content Management

## Overview

ArchesAI's Content Management system provides comprehensive artifact storage, vector embeddings for
semantic search, and intelligent content processing pipelines.

## Core Features

### Artifact Storage

- **Multi-format Support**: Documents, images, videos, structured data
- **Version Control**: Track content changes over time
- **Metadata Management**: Rich metadata and tagging
- **Cloud Storage**: S3-compatible object storage integration

### Vector Embeddings

- **Semantic Search**: Find content by meaning, not just keywords
- **Multiple Models**: Support for OpenAI, Cohere, and custom embeddings
- **pgvector Integration**: PostgreSQL vector similarity search
- **Hybrid Search**: Combine vector and traditional search

### Content Processing

- **Automatic Extraction**: Text, metadata, and structure extraction
- **Format Conversion**: Convert between document formats
- **Thumbnail Generation**: Automatic preview generation
- **Content Enrichment**: AI-powered tagging and categorization

## Architecture

### Content Model

```go
type Content struct {
    ID           uuid.UUID
    Type         ContentType
    Title        string
    Description  string
    MimeType     string
    Size         int64
    StorageURL   string
    Embedding    pgvector.Vector
    Metadata     JSONB
    Tags         []string
    Version      int
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type ContentType string

const (
    ContentTypeDocument ContentType = "document"
    ContentTypeImage    ContentType = "image"
    ContentTypeVideo    ContentType = "video"
    ContentTypeAudio    ContentType = "audio"
    ContentTypeData     ContentType = "data"
)
```

### Storage Architecture

```yaml
storage:
  primary:
    type: s3
    bucket: archesai-content
    region: us-east-1

  cache:
    type: redis
    ttl: 3600

  database:
    type: postgresql
    extensions:
      - pgvector
      - pg_trgm
```

## Content Operations

### Upload and Storage

```typescript
// Upload content
POST /api/v1/content
Content-Type: multipart/form-data

{
  "file": File,
  "title": "Q4 Report",
  "tags": ["finance", "quarterly"],
  "metadata": {
    "department": "finance",
    "year": 2024
  }
}

// Response
{
  "id": "content_123",
  "url": "https://storage.archesai.com/content_123",
  "embedding_status": "processing",
  "processing_status": "queued"
}
```

### Search Functionality

#### Semantic Search

```typescript
// Vector similarity search
POST /api/v1/content/search
{
  "query": "financial performance metrics",
  "type": "semantic",
  "limit": 10,
  "threshold": 0.8
}
```

#### Hybrid Search

```typescript
// Combined vector and keyword search
POST /api/v1/content/search
{
  "query": "Q4 2024 revenue",
  "type": "hybrid",
  "filters": {
    "tags": ["finance"],
    "date_range": {
      "start": "2024-10-01",
      "end": "2024-12-31"
    }
  },
  "weights": {
    "semantic": 0.7,
    "keyword": 0.3
  }
}
```

### Processing Pipeline

```mermaid
graph LR
    A[Upload] --> B[Validation]
    B --> C[Storage]
    C --> D[Extraction]
    D --> E[Embedding]
    E --> F[Indexing]
    F --> G[Available]
```

## Vector Embeddings

### Embedding Generation

```go
func GenerateEmbedding(content Content) ([]float32, error) {
    // Extract text from content
    text := ExtractText(content)

    // Generate embedding using configured model
    embedding := embeddingModel.Embed(text)

    // Store in pgvector
    err := db.StoreEmbedding(content.ID, embedding)

    return embedding, err
}
```

### Similarity Search

```sql
-- Find similar content using pgvector
SELECT
    id,
    title,
    1 - (embedding <=> $1) as similarity
FROM content
WHERE 1 - (embedding <=> $1) > $2
ORDER BY embedding <=> $1
LIMIT $3;
```

### Embedding Models

#### OpenAI

- Model: text-embedding-ada-002
- Dimensions: 1536
- Best for: General purpose

#### Cohere

- Model: embed-english-v3.0
- Dimensions: 1024
- Best for: Domain-specific

#### Custom Models

- Sentence Transformers
- Domain-trained models
- Fine-tuned embeddings

## Processing Capabilities

### Document Processing

#### Text Extraction

- PDF text extraction with OCR
- Word document parsing
- HTML content extraction
- Markdown processing

#### Metadata Extraction

- Author information
- Creation/modification dates
- Document properties
- Custom metadata fields

### Image Processing

#### Analysis

- Object detection
- Face detection
- OCR for text in images
- EXIF data extraction

#### Transformation

- Thumbnail generation
- Format conversion
- Compression
- Watermarking

### Video Processing

#### Extraction

- Frame extraction
- Audio transcription
- Scene detection
- Metadata parsing

#### Generation

- Preview clips
- Thumbnails
- Transcripts
- Subtitles

## API Endpoints

### Content Management

- `POST /api/v1/content` - Upload content
- `GET /api/v1/content/:id` - Get content details
- `PUT /api/v1/content/:id` - Update content metadata
- `DELETE /api/v1/content/:id` - Delete content
- `GET /api/v1/content/:id/download` - Download content

### Search

- `POST /api/v1/content/search` - Search content
- `GET /api/v1/content/similar/:id` - Find similar content
- `POST /api/v1/content/query` - Advanced query

### Processing

- `POST /api/v1/content/:id/process` - Trigger processing
- `GET /api/v1/content/:id/status` - Processing status
- `POST /api/v1/content/:id/extract` - Extract specific data

### Embeddings

- `POST /api/v1/content/:id/embed` - Generate embedding
- `GET /api/v1/content/:id/embedding` - Get embedding
- `PUT /api/v1/content/:id/embedding` - Update embedding

## Storage Configuration

### S3 Configuration

```yaml
storage:
  s3:
    endpoint: https://s3.amazonaws.com
    bucket: archesai-content
    region: us-east-1
    access_key: ${AWS_ACCESS_KEY_ID}
    secret_key: ${AWS_SECRET_ACCESS_KEY}

    # Lifecycle policies
    lifecycle:
      - rule: archive_old_content
        transition:
          days: 90
          storage_class: GLACIER

      - rule: delete_temp_files
        expiration:
          days: 7
          prefix: temp/
```

### Local Storage

```yaml
storage:
  local:
    base_path: /var/archesai/content
    max_size: 100GB

    # Cleanup policies
    cleanup:
      temp_files:
        max_age: 24h
      orphaned_files:
        check_interval: 1h
```

## Performance Optimization

### Caching Strategy

#### Redis Cache

```go
// Cache frequently accessed content
func CacheContent(id string, content []byte) {
    key := fmt.Sprintf("content:%s", id)
    redis.Set(key, content, 1*time.Hour)
}

// Cache search results
func CacheSearchResults(query string, results []Content) {
    key := fmt.Sprintf("search:%s", hash(query))
    redis.Set(key, results, 15*time.Minute)
}
```

#### CDN Integration

- CloudFront for static content
- Edge caching for thumbnails
- Geo-distributed content delivery

### Database Optimization

#### Indexes

```sql
-- Vector similarity index
CREATE INDEX content_embedding_idx ON content
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- Full-text search index
CREATE INDEX content_search_idx ON content
USING gin(to_tsvector('english', title || ' ' || description));

-- Metadata JSONB index
CREATE INDEX content_metadata_idx ON content
USING gin(metadata);
```

### Processing Optimization

- Async processing queues
- Parallel extraction
- Batch embedding generation
- Incremental indexing

## Security

### Access Control

- Row-level security
- Content encryption at rest
- Signed URLs for downloads
- IP-based restrictions

### Data Protection

- Automatic PII detection
- Content sanitization
- Virus scanning
- DLP integration

### Audit Trail

```json
{
  "event": "content_accessed",
  "content_id": "content_123",
  "user_id": "user_456",
  "action": "download",
  "timestamp": "2024-01-15T10:30:00Z",
  "ip_address": "192.168.1.1",
  "user_agent": "Mozilla/5.0..."
}
```

## Integration

### Workflow Integration

- Content as workflow triggers
- Processing pipelines
- Automated tagging
- Content routing

### External Systems

- SharePoint connector
- Google Drive sync
- Dropbox integration
- Box.com support

### APIs and Webhooks

```javascript
// Webhook notification
{
  "event": "content.processed",
  "content": {
    "id": "content_123",
    "status": "completed",
    "results": {
      "text_extracted": true,
      "embedding_generated": true,
      "thumbnails_created": 3
    }
  }
}
```

## Monitoring

### Metrics

- Upload/download rates
- Processing queue depth
- Search response times
- Storage utilization
- Embedding generation time

### Health Checks

- Storage connectivity
- Database performance
- Processing worker status
- Cache hit rates

## Best Practices

### Content Organization

- Use consistent naming conventions
- Apply comprehensive tagging
- Maintain metadata standards
- Regular cleanup procedures

### Search Optimization

- Pre-generate embeddings
- Use appropriate embedding models
- Optimize vector dimensions
- Cache frequent queries

### Storage Management

- Implement lifecycle policies
- Use appropriate storage tiers
- Regular backup procedures
- Monitor storage costs

## Troubleshooting

### Common Issues

#### Slow Search Performance

- Check pgvector indexes
- Verify embedding dimensions
- Review query complexity
- Monitor database load

#### Processing Failures

- Check file format support
- Verify processing worker status
- Review error logs
- Check resource limits

#### Storage Issues

- Verify S3 permissions
- Check storage quotas
- Review lifecycle policies
- Monitor bandwidth usage

## Related Documentation

- [Workflows](workflows.md) - Content processing workflows
