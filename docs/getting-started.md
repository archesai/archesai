# ArchesAI Documentation

Welcome to the ArchesAI documentation! This guide will help you understand, develop, and deploy the ArchesAI platform.

## What is ArchesAI?

ArchesAI is a high-performance data processing platform with AI-powered chat interface and workflow automation. It provides:

- **Multi-Provider AI**: Support for OpenAI, Claude, Gemini, Ollama
- **Chat Interface**: Simple persona-based chat system with session management
- **Beautiful TUI**: Terminal interface for configuration and chat
- **Workflow Automation**: DAG-based data processing pipelines
- **Code Generation**: OpenAPI and SQL-driven development
- **Modern Stack**: Go, PostgreSQL/SQLite, Redis

## Quick Start

### 1. Configuration Viewer (No API Key Needed)

```bash
archesai tui
```

### 2. API Server

```bash
archesai api
```

### 3. AI Chat Interface

```bash
export OPENAI_API_KEY=your-key
archesai tui --chat

# Or with Ollama (local)
archesai tui --chat --provider=ollama
```

## Documentation Structure

This documentation is organized into several sections:

### 🏗️ [Architecture](architecture/system-design.md)

Learn about the system design, patterns, and overall architecture of ArchesAI.

### 🚀 [Development](guides/development.md)

Everything you need to know about setting up your development environment and contributing to the project.

### 📚 [API Reference](api-reference/overview.md)

Complete API documentation with endpoints, schemas, and examples.

### 🎯 [Features](features/overview.md)

Detailed guides for each feature domain: authentication, organizations, workflows, and content management.

### 🐳 [Deployment](deployment/overview.md)

Production deployment guides including Docker, Kubernetes, and infrastructure setup.

### 🔧 [Troubleshooting](troubleshooting/common-issues.md)

Common issues, debugging guides, and solutions.

### 🔒 [Security](security/overview.md)

Security best practices and guidelines.

### ⚡ [Performance](performance/overview.md)

Performance optimization and monitoring guides.

## Need Help?

- **Email**: <support@archesai.com>
- **Issues**: [GitHub Issues](https://github.com/archesai/archesai/issues)
- **Contributing**: See our [Contributing Guide](contributing.md)

## License

See [LICENSE](https://github.com/archesai/archesai/blob/main/LICENSE) file for details.
