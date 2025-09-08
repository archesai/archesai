# Contributing to ArchesAI

Thank you for your interest in contributing to ArchesAI! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- **Be respectful**: Treat all contributors with respect
- **Be constructive**: Provide helpful feedback and suggestions
- **Be inclusive**: Welcome contributors of all backgrounds and experience levels
- **Be professional**: Keep discussions focused on the project

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- Go 1.21+ installed
- PostgreSQL 15+ with pgvector extension
- Node.js 20+ and pnpm 8+
- Git configured with your GitHub account
- Make installed for running build commands

### Setting Up Your Development Environment

1. **Fork the repository** on GitHub

2. **Clone your fork**:

   ```bash
   git clone https://github.com/YOUR_USERNAME/archesai.git
   cd archesai
   ```

3. **Add upstream remote**:

   ```bash
   git remote add upstream https://github.com/archesai/archesai.git
   ```

4. **Set up the project**:

   ```bash
   # Copy environment variables
   cp .env.example .env

   # Install dependencies
   make tools
   go mod download
   pnpm install

   # Set up database
   createdb archesai
   psql archesai -c "CREATE EXTENSION IF NOT EXISTS vector;"
   make migrate-up

   # Generate code
   make generate
   ```

5. **Verify setup**:
   ```bash
   make test
   make lint
   ```

## Development Workflow

### Branch Naming

Use descriptive branch names:

- `feature/add-webhook-support`
- `fix/auth-token-expiry`
- `docs/update-api-guide`
- `refactor/optimize-queries`
- `test/add-workflow-tests`

### Making Changes

1. **Create a new branch**:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards

3. **Generate code if needed**:

   ```bash
   make generate  # After API or database changes
   ```

4. **Run tests**:

   ```bash
   make test
   ```

5. **Lint your code**:

   ```bash
   make lint
   ```

6. **Format your code**:

   ```bash
   make format
   ```

7. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

### Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements

Examples:

```
feat(auth): add refresh token rotation
fix(workflows): resolve DAG cycle detection bug
docs(api): update authentication examples
refactor(database): optimize artifact queries
```

### Submitting a Pull Request

1. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request** on GitHub

3. **Fill out the PR template** with:
   - Description of changes
   - Related issue numbers
   - Testing performed
   - Screenshots (if UI changes)

4. **Wait for review** - maintainers will review your PR

5. **Address feedback** - make requested changes and push updates

## Coding Standards

### Go Code

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Write descriptive variable and function names
- Add comments for exported functions
- Handle errors explicitly
- Use interfaces for dependency injection
- Write table-driven tests

Example:

```go
// CreateUser creates a new user account with the provided details.
// It returns an error if the email is already registered.
func (s *AuthService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    // Check if user exists
    existing, _ := s.repo.GetUserByEmail(ctx, req.Email)
    if existing != nil {
        return nil, ErrUserExists
    }

    // Create user
    user := &User{
        ID:    uuid.New(),
        Email: req.Email,
        Name:  req.Name,
    }

    if err := s.repo.CreateUser(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}
```

### TypeScript/React Code

- Use TypeScript for all new code
- Follow React best practices and hooks
- Use functional components
- Implement proper error boundaries
- Write meaningful component and prop names
- Add JSDoc comments for complex functions

Example:

```typescript
/**
 * WorkflowEditor component for creating and editing workflow DAGs
 */
export const WorkflowEditor: React.FC<WorkflowEditorProps> = ({
  workflow,
  onSave,
  readonly = false,
}) => {
  const [nodes, setNodes] = useState<Node[]>(workflow?.nodes || []);
  const [edges, setEdges] = useState<Edge[]>(workflow?.edges || []);

  const handleNodeAdd = useCallback((node: Node) => {
    setNodes((prev) => [...prev, node]);
  }, []);

  // ... rest of component
};
```

### SQL Queries

- Use lowercase for SQL keywords
- Use snake_case for column names
- Add comments for complex queries
- Use parameterized queries (never concatenate)
- Follow SQLC conventions

Example:

```sql
-- name: GetUserByEmail :one
-- GetUserByEmail retrieves a user by their email address
select
  id,
  email,
  name,
  email_verified,
  created_at,
  updated_at
from users
where email = $1
limit 1;
```

### API Design

- Follow RESTful principles
- Use proper HTTP status codes
- Implement pagination for list endpoints
- Use consistent error responses
- Document all endpoints in OpenAPI

## Testing Guidelines

### Go Tests

Write table-driven tests:

```go
func TestCreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserRequest
        want    *User
        wantErr error
    }{
        {
            name: "valid user",
            input: CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            want: &User{
                Email: "test@example.com",
                Name:  "Test User",
            },
            wantErr: nil,
        },
        {
            name: "duplicate email",
            input: CreateUserRequest{
                Email: "existing@example.com",
                Name:  "Test User",
            },
            want:    nil,
            wantErr: ErrUserExists,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### TypeScript Tests

Use Vitest for unit tests:

```typescript
describe('WorkflowEditor', () => {
  it('should render workflow nodes', () => {
    const workflow = {
      nodes: [{ id: '1', type: 'input', data: {} }],
      edges: [],
    };

    const { getByTestId } = render(
      <WorkflowEditor workflow={workflow} onSave={jest.fn()} />
    );

    expect(getByTestId('node-1')).toBeInTheDocument();
  });
});
```

## Documentation

### Code Documentation

- Add package-level documentation
- Document all exported functions and types
- Include examples for complex functionality
- Keep documentation up to date with code changes

### README Updates

Update the README when:

- Adding new features
- Changing setup procedures
- Modifying configuration options
- Adding new dependencies

### API Documentation

- Update OpenAPI spec for API changes
- Include request/response examples
- Document error responses
- Add descriptions for all parameters

## Architecture Guidelines

### Hexagonal Architecture

Follow the hexagonal architecture pattern:

1. **Domain Layer** (Core):
   - Business logic
   - Domain entities
   - Use cases
   - Port interfaces

2. **Infrastructure Layer** (Adapters):
   - Database implementations
   - External service clients
   - File system operations

3. **Application Layer** (Handlers):
   - HTTP handlers
   - Middleware
   - Request/response mapping

### Adding New Features

When adding a new feature:

1. **Define the API** in `api/openapi.yaml`
2. **Create database migration** if needed
3. **Write SQL queries** in `queries/`
4. **Generate code**: `make generate`
5. **Implement domain logic** in use cases
6. **Implement repository** if needed
7. **Create HTTP handler**
8. **Wire dependencies** in `deps.go`
9. **Add routes** in `routes.go`
10. **Write tests**
11. **Update documentation**

## Review Process

### What We Look For

- **Code quality**: Clean, maintainable, efficient code
- **Tests**: Adequate test coverage for new functionality
- **Documentation**: Clear comments and updated docs
- **Architecture**: Follows project patterns and conventions
- **Performance**: No significant performance regressions
- **Security**: No security vulnerabilities introduced

### Review Timeline

- Initial review: Within 2-3 business days
- Follow-up reviews: Within 1-2 business days
- Small fixes: Usually same day

## Getting Help

### Resources

- [Development Guide](DEVELOPMENT.md)
- [Architecture Documentation](ARCHITECTURE.md)
- [API Documentation](API.md)
- [Project Issues](https://github.com/archesai/archesai/issues)

### Communication

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and discussions
- **Pull Requests**: Code contributions and reviews

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):

- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)

### Release Cycle

- **Patch releases**: As needed for critical fixes
- **Minor releases**: Monthly
- **Major releases**: Quarterly or as needed

## Recognition

Contributors are recognized in:

- Release notes
- Contributors file
- Project documentation

Thank you for contributing to ArchesAI! Your efforts help make the project better for everyone.
