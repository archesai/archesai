# Arches Roadmap

This document outlines planned features and enhancements for Arches. These features are
not yet implemented but represent the future direction of the project.

## Status Key

- 🔵 **Planned** - On the roadmap, timeline TBD
- 🟡 **In Progress** - Active development
- 🟢 **Partially Complete** - Some functionality exists
- ⏸️ **On Hold** - Paused pending other work

## Q4 2025 Goals

### CLI Enhancements

- 🔵 **`archesai new` command** - Scaffold new Arches projects from templates
- 🔵 **`archesai init` command** - Initialize Arches in existing projects
- 🟡 **`archesai deploy` command** - One-command deployment to cloud platforms

### Core Features

#### Visual Studio UI

- 🔵 **Web-based Schema Designer** - Drag-and-drop OpenAPI schema builder
  - Visual entity relationship diagram
  - Real-time validation
  - AI-powered suggestions
- 🔵 **Live Preview** - See generated app as you build
- 🔵 **Git Integration** - Version control within the studio
- 🔵 **Collaborative Editing** - Real-time multi-user editing

#### Code Generation Improvements

- 🟢 **Multi-language Support**
  - ✅ Go (complete)
  - 🟡 Python (runner exists, templates in progress)
  - 🔵 Node.js/TypeScript (planned)
  - 🔵 Rust (planned)
- 🔵 **Custom Template System** - User-defined generation templates
- 🔵 **Plugin Architecture** - Extend generation with custom plugins

#### AI Features

- 🟡 **Natural Language to OpenAPI** - Describe your API in plain English
- 🔵 **Smart Handler Implementation** - AI generates handler logic from descriptions
- 🔵 **Test Generation** - Automatic test case creation
- 🔵 **Documentation Generation** - AI-written API documentation

## Q2 2025 Goals

### Enterprise Features

- 🔵 **Team Management** - Organizations, roles, permissions
- 🔵 **Private Schema Registry** - Share and reuse API components
- 🔵 **CI/CD Integration** - GitHub Actions, GitLab CI, Jenkins plugins
- 🔵 **Audit Logging** - Track all changes and generations

### Platform Enhancements

- 🔵 **Cloud Deployments**
  - AWS (ECS, Lambda, RDS)
  - Google Cloud (Cloud Run, Cloud SQL)
  - Azure (Container Instances, SQL Database)
  - Vercel/Netlify for frontend
- 🔵 **Database Flexibility**
  - MySQL support
  - MongoDB support
  - Supabase integration
- 🔵 **Monitoring & Observability**
  - Built-in APM
  - Distributed tracing
  - Custom metrics dashboards

### Developer Experience

- 🔵 **Hot Module Replacement** - Frontend HMR during development
- 🔵 **Database Migrations UI** - Visual migration management
- 🔵 **API Testing Suite** - Built-in API testing tools
- 🔵 **Performance Profiler** - Identify bottlenecks

## Q3 2025 Goals

### Advanced Features

- 🔵 **GraphQL Support** - Generate GraphQL APIs from OpenAPI
- 🔵 **WebSocket Subscriptions** - Real-time event streaming
- 🔵 **Message Queue Integration** - RabbitMQ, Kafka support
- 🔵 **Microservices Mode** - Generate service mesh architectures

### Ecosystem

- 🔵 **Marketplace** - Share and sell templates, plugins
- 🔵 **Community Templates** - Pre-built app templates
- 🔵 **Integration Hub** - Pre-built integrations with popular services

## Long-term Vision

### Platform Evolution

- 🔵 **No-Code Mode** - Full visual development without code
- 🔵 **Mobile App Generation** - React Native/Flutter from OpenAPI
- 🔵 **Desktop App Generation** - Electron apps
- 🔵 **Edge Deployment** - Cloudflare Workers, Deno Deploy

### AI Evolution

- 🔵 **Autonomous Development** - AI handles entire features end-to-end
- 🔵 **Code Review AI** - Automated PR reviews and suggestions
- 🔵 **Performance Optimization AI** - Automatic performance improvements
- 🔵 **Security Scanning** - AI-powered vulnerability detection

## Recently Completed

### ✅ Core Platform (Completed)

- OpenAPI to Go code generation
- Database schema generation
- Basic authentication/authorization
- Docker containerization
- Kubernetes manifests
- CLI tooling foundation

### ✅ Development Tools (Completed)

- Hot reload development mode
- Code generation pipeline
- Testing infrastructure
- CI/CD workflows

## Contributing

Want to help accelerate our roadmap? Check out our [Contributing Guide](contributing.md) to get started.

## Feature Requests

Have ideas for features not on this roadmap? Please
[open an issue](https://github.com/archesai/archesai/issues) with the "feature-request" label.

---

_Last updated: November 2024_
_This roadmap is subject to change based on community feedback and project priorities._
