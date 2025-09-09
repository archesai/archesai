# API Reference

Welcome to the ArchesAI API reference documentation. Our REST API is built using OpenAPI 3.0
specifications and provides comprehensive access to all platform features.

## Overview

The ArchesAI API is organized around REST principles with predictable URLs, standard HTTP response
codes, and JSON request/response bodies. All API endpoints are prefixed with the API version.

**Base URL**: `https://api.archesai.com/v1`

## Authentication

All API requests require authentication using Bearer tokens. See
[Authentication](../architecture/authentication.md) for detailed information about obtaining and
using API keys.

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
     https://api.archesai.com/v1/organizations
```

## API Domains

### [Authentication](../architecture/authentication.md)

User authentication, session management, and OAuth integration.

### Organizations

Multi-tenant organization management, member invitations, and role-based access control.

### Workflows

Pipeline creation, execution, and monitoring with DAG-based workflow automation.

### Content

Artifact storage, vector embeddings, and content processing operations.

## Response Format

All API responses follow a consistent format:

```json
{
  "data": {
    // Response data
  },
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

### Error Responses

Errors follow RFC 7807 Problem Details format:

```json
{
  "type": "/errors/validation-error",
  "title": "Validation Error",
  "status": 400,
  "detail": "The request body is invalid",
  "instance": "/v1/organizations"
}
```

## Rate Limiting

API requests are rate limited to prevent abuse:

- **Free tier**: 100 requests per minute
- **Pro tier**: 1000 requests per minute
- **Enterprise**: Custom limits

Rate limit headers are included in all responses:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## SDKs and Tools

- **Go SDK**: Generated client in `web/client/`
- **TypeScript SDK**: NPM package `@archesai/client`
- **OpenAPI Spec**: Available in the `api/` directory
- **Postman Collection**: [Import](https://api.archesai.com/postman)

## Interactive API Explorer

Try our interactive API explorer at [api.archesai.com/docs](https://api.archesai.com/docs) to test
endpoints directly from your browser.
