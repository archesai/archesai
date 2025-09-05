# API Documentation

This document provides comprehensive documentation for the ArchesAI REST API.

## Base URL

```
http://localhost:8080/api
```

## Authentication

The API uses JWT-based authentication with refresh tokens. Most endpoints require authentication.

### Token Types

- **Access Token**: Short-lived token (15 minutes) used for API requests
- **Refresh Token**: Long-lived token (7 days) used to obtain new access tokens

### Headers

Authenticated requests must include the access token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## Error Responses

All endpoints follow a consistent error response format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {} // Optional additional details
  }
}
```

Common HTTP status codes:

- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation errors
- `500 Internal Server Error` - Server error

## Endpoints

### Authentication

#### Register User

Create a new user account.

**Endpoint:** `POST /api/auth/register`

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "name": "John Doe"
}
```

**Response:** `201 Created`

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "emailVerified": false,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "tokens": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 900
  }
}
```

#### Login

Authenticate user and receive tokens.

**Endpoint:** `POST /api/auth/login`

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response:** `200 OK`

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "emailVerified": true,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "tokens": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 900
  }
}
```

#### Refresh Token

Get a new access token using a refresh token.

**Endpoint:** `POST /api/auth/refresh`

**Request Body:**

```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:** `200 OK`

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "expiresIn": 900
}
```

#### Logout

Invalidate the current session.

**Endpoint:** `POST /api/auth/logout`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`

```json
{
  "message": "Successfully logged out"
}
```

#### Get Current User

Get information about the authenticated user.

**Endpoint:** `GET /api/auth/me`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "emailVerified": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Organizations

#### List Organizations

Get a paginated list of organizations the user has access to.

**Endpoint:** `GET /api/organizations`

**Query Parameters:**

- `limit` (integer, optional): Number of items per page (default: 50, max: 100)
- `offset` (integer, optional): Number of items to skip (default: 0)
- `sort` (string, optional): Sort field (name, createdAt, updatedAt)
- `order` (string, optional): Sort order (asc, desc)

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Acme Corporation",
      "slug": "acme-corp",
      "description": "Leading provider of innovative solutions",
      "ownerId": "user-uuid",
      "settings": {},
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 100,
    "limit": 50,
    "offset": 0,
    "hasMore": true
  }
}
```

#### Create Organization

Create a new organization.

**Endpoint:** `POST /api/organizations`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "name": "Acme Corporation",
  "slug": "acme-corp",
  "description": "Leading provider of innovative solutions",
  "settings": {
    "allowPublicAccess": false
  }
}
```

**Response:** `201 Created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corporation",
  "slug": "acme-corp",
  "description": "Leading provider of innovative solutions",
  "ownerId": "user-uuid",
  "settings": {
    "allowPublicAccess": false
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Get Organization

Get details of a specific organization.

**Endpoint:** `GET /api/organizations/{id}`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corporation",
  "slug": "acme-corp",
  "description": "Leading provider of innovative solutions",
  "ownerId": "user-uuid",
  "settings": {},
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Update Organization

Update an organization's details.

**Endpoint:** `PUT /api/organizations/{id}`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "name": "Acme Corporation Updated",
  "description": "Updated description",
  "settings": {
    "allowPublicAccess": true
  }
}
```

**Response:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corporation Updated",
  "slug": "acme-corp",
  "description": "Updated description",
  "ownerId": "user-uuid",
  "settings": {
    "allowPublicAccess": true
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-02T00:00:00Z"
}
```

#### Delete Organization

Delete an organization.

**Endpoint:** `DELETE /api/organizations/{id}`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response:** `204 No Content`

### Workflows

#### List Workflows

Get a paginated list of workflows.

**Endpoint:** `GET /api/workflows`

**Query Parameters:**

- `limit` (integer, optional): Number of items per page (default: 50)
- `offset` (integer, optional): Number of items to skip (default: 0)
- `organizationId` (string, optional): Filter by organization
- `status` (string, optional): Filter by status (draft, published, archived)

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Document Processing Pipeline",
      "description": "Processes and analyzes documents",
      "organizationId": "org-uuid",
      "status": "published",
      "config": {
        "nodes": [],
        "edges": []
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 25,
    "limit": 50,
    "offset": 0,
    "hasMore": false
  }
}
```

#### Create Workflow

Create a new workflow.

**Endpoint:** `POST /api/workflows`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "name": "Document Processing Pipeline",
  "description": "Processes and analyzes documents",
  "organizationId": "org-uuid",
  "config": {
    "nodes": [
      {
        "id": "node1",
        "type": "input",
        "data": {}
      },
      {
        "id": "node2",
        "type": "text-extraction",
        "data": {}
      }
    ],
    "edges": [
      {
        "source": "node1",
        "target": "node2"
      }
    ]
  }
}
```

**Response:** `201 Created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Document Processing Pipeline",
  "description": "Processes and analyzes documents",
  "organizationId": "org-uuid",
  "status": "draft",
  "config": {
    "nodes": [...],
    "edges": [...]
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Execute Workflow

Execute a workflow to create a pipeline run.

**Endpoint:** `POST /api/workflows/{id}/runs`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "input": {
    "fileUrl": "https://example.com/document.pdf",
    "parameters": {
      "language": "en"
    }
  }
}
```

**Response:** `201 Created`

```json
{
  "id": "run-uuid",
  "workflowId": "workflow-uuid",
  "status": "running",
  "input": {...},
  "output": null,
  "startedAt": "2024-01-01T00:00:00Z",
  "completedAt": null
}
```

### Content/Artifacts

#### List Artifacts

Get a paginated list of content artifacts.

**Endpoint:** `GET /api/content/artifacts`

**Query Parameters:**

- `limit` (integer, optional): Number of items per page (default: 50)
- `offset` (integer, optional): Number of items to skip (default: 0)
- `organizationId` (string, optional): Filter by organization
- `type` (string, optional): Filter by type (document, image, audio, video)
- `labels` (string, optional): Comma-separated list of labels

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "artifact-uuid",
      "name": "report.pdf",
      "type": "document",
      "mimeType": "application/pdf",
      "size": 1048576,
      "organizationId": "org-uuid",
      "metadata": {
        "pages": 10,
        "author": "John Doe"
      },
      "labels": ["report", "q4-2024"],
      "embedding": null,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 150,
    "limit": 50,
    "offset": 0,
    "hasMore": true
  }
}
```

#### Upload Artifact

Upload a new content artifact.

**Endpoint:** `POST /api/content/artifacts`

**Headers:**

```
Authorization: Bearer <access_token>
Content-Type: multipart/form-data
```

**Request Body (multipart/form-data):**

- `file` (file, required): The file to upload
- `organizationId` (string, required): Organization ID
- `labels` (string, optional): Comma-separated labels
- `metadata` (JSON string, optional): Additional metadata

**Response:** `201 Created`

```json
{
  "id": "artifact-uuid",
  "name": "report.pdf",
  "type": "document",
  "mimeType": "application/pdf",
  "size": 1048576,
  "organizationId": "org-uuid",
  "url": "https://storage.archesai.com/artifacts/artifact-uuid",
  "metadata": {},
  "labels": ["report"],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Get Artifact

Get details of a specific artifact.

**Endpoint:** `GET /api/content/artifacts/{id}`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`

```json
{
  "id": "artifact-uuid",
  "name": "report.pdf",
  "type": "document",
  "mimeType": "application/pdf",
  "size": 1048576,
  "organizationId": "org-uuid",
  "url": "https://storage.archesai.com/artifacts/artifact-uuid",
  "metadata": {
    "pages": 10,
    "author": "John Doe"
  },
  "labels": ["report", "q4-2024"],
  "embedding": [0.1, 0.2, ...],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Process Artifact

Process an artifact using a specific transformation.

**Endpoint:** `POST /api/content/artifacts/{id}/process`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "operation": "extract-text",
  "parameters": {
    "format": "markdown",
    "preserveFormatting": true
  }
}
```

**Response:** `200 OK`

```json
{
  "id": "process-job-uuid",
  "artifactId": "artifact-uuid",
  "operation": "extract-text",
  "status": "processing",
  "result": null,
  "startedAt": "2024-01-01T00:00:00Z",
  "completedAt": null
}
```

#### Delete Artifact

Delete an artifact.

**Endpoint:** `DELETE /api/content/artifacts/{id}`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response:** `204 No Content`

### Search

#### Semantic Search

Search for content using semantic similarity.

**Endpoint:** `POST /api/search/semantic`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "query": "financial reports from Q4",
  "organizationId": "org-uuid",
  "limit": 10,
  "threshold": 0.7,
  "filters": {
    "type": "document",
    "labels": ["report"]
  }
}
```

**Response:** `200 OK`

```json
{
  "results": [
    {
      "id": "artifact-uuid",
      "name": "Q4-financial-report.pdf",
      "type": "document",
      "score": 0.92,
      "snippet": "...quarterly financial results...",
      "metadata": {}
    }
  ],
  "total": 5,
  "query": "financial reports from Q4"
}
```

## Rate Limiting

API endpoints are rate-limited to prevent abuse:

- **Authentication endpoints**: 5 requests per minute per IP
- **Regular endpoints**: 100 requests per minute per user
- **File uploads**: 10 requests per minute per user

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704067200
```

## Pagination

List endpoints support pagination using `limit` and `offset` parameters:

- `limit`: Number of items to return (max: 100)
- `offset`: Number of items to skip

Paginated responses include a `pagination` object:

```json
{
  "data": [...],
  "pagination": {
    "total": 250,
    "limit": 50,
    "offset": 0,
    "hasMore": true
  }
}
```

## Filtering and Sorting

List endpoints support filtering and sorting:

### Filtering

Use query parameters to filter results:

```
GET /api/workflows?status=published&organizationId=org-uuid
```

### Sorting

Use `sort` and `order` parameters:

```
GET /api/artifacts?sort=createdAt&order=desc
```

## Webhooks (Coming Soon)

Webhooks will allow you to receive real-time notifications about events in your ArchesAI account.

### Planned Events

- `workflow.completed` - Workflow execution completed
- `artifact.processed` - Artifact processing completed
- `organization.updated` - Organization settings updated

## SDK and Client Libraries

### TypeScript/JavaScript

An auto-generated TypeScript client is available:

```typescript
import { ArchesAIClient } from "@archesai/client";

const client = new ArchesAIClient({
  baseURL: "https://api.archesai.com",
  accessToken: "your-access-token",
});

// List workflows
const workflows = await client.workflows.list({
  limit: 10,
  organizationId: "org-uuid",
});

// Upload artifact
const artifact = await client.artifacts.upload({
  file: fileBlob,
  organizationId: "org-uuid",
  labels: ["document", "report"],
});
```

### Other Languages

SDKs for other languages are planned:

- Python
- Go
- Java
- Ruby

## API Versioning

The API uses URL versioning. The current version is v1:

```
https://api.archesai.com/v1/...
```

Breaking changes will result in a new API version. Deprecated versions will be supported for at least 6 months after a new version is released.

## Support

For API support:

- Documentation: https://docs.archesai.com/api
- Issues: https://github.com/archesai/archesai/issues
- Email: api-support@archesai.com
