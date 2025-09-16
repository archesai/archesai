# Testing Documentation

## Overview

This document outlines the testing strategy, patterns, and technical decisions for the Arches
backend.

## Testing Philosophy

- **Minimal Dependencies**: Tests use Go's standard library with minimal external dependencies
- **Mock Repository Pattern**: All database interactions are mocked using in-memory implementations
- **Table-Driven Tests**: Comprehensive test cases using Go's table-driven test pattern
- **Domain Isolation**: Each domain is tested independently without cross-domain dependencies

## Technical Decisions

### 1. Mock Repository Pattern

We use mockery v3 to generate mocks from interfaces. The unified code generator creates repository interfaces
that mockery then uses to generate test mocks.

```bash
# Generate repository interfaces from x-codegen
make generate-codegen

# Generate mocks from interfaces
make generate-mocks

# Or generate everything
make generate
```

**Benefits:**

- Automatic generation from generated interfaces
- Consistent mock patterns across all domains
- Type-safe mocks matching repository contracts
- Integration with testify assertions

### 2. Test Organization

Each domain follows this structure:

```text
internal/{domain}/
├── service_test.go      # Business logic tests
├── handler_test.go      # HTTP handler tests (if applicable)
├── middleware_test.go   # Middleware tests (if applicable)
├── mocks_test.gen.go    # Generated mocks from interfaces
└── postgres_test.go     # Integration tests with real database
```

### 3. Integration Testing

For database integration tests, we use testcontainers:

- PostgreSQL containers for testing real database operations
- Redis containers for testing cache operations
- Automatic cleanup after tests complete
- Located in `internal/testutil/containers.go`

### 4. Coverage Strategy

Current coverage by package:

- `internal/auth` - 20.2%
- `internal/sessions` - 27.2%
- `internal/config` - 47.2%
- `internal/artifacts` - 17.2%
- `internal/labels` - 16.3%
- `internal/pipelines` - 12.8%
- `internal/users` - 12.2%
- `internal/runs` - 6.4%
- `internal/health` - 3.9%
- `internal/organizations` - 3.3%
- `internal/tools` - 0.0%
- `internal/accounts` - 0.0%
- `internal/invitations` - 0.0%
- `internal/members` - 0.0%

Target: 80%+ coverage for critical business logic

## Test Patterns

### Service Tests

```go
func TestService_Method(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        setup   func(*MockRepository)  // Prepare mock state
        wantErr bool
    }{
        // Test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := NewMockRepository()
            tt.setup(repo)
            service := NewService(repo, slog.Default())

            result, err := service.Method(context.Background(), tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("Method() error = %v, wantErr %v", err, tt.wantErr)
            }
            // Additional assertions
        })
    }
}
```

### Handler Tests

```go
func TestHandler_Endpoint(t *testing.T) {
    // Create mock service
    mockService := &MockService{}
    handler := NewHandler(mockService, slog.Default())

    // Create Echo context with test request
    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/", body)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Execute handler
    response, err := handler.Endpoint(c, request)

    // Assert response
}
```

### Integration Tests

```go
func TestPostgresRepository_Operations(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    ctx := context.Background()
    pgContainer := testutil.StartPostgresContainer(ctx, t)

    // Run migrations
    err := pgContainer.RunMigrations("../../migrations")
    if err != nil {
        t.Fatalf("Failed to run migrations: %v", err)
    }

    // Test database operations
}
```

## Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific domain tests
go test ./internal/auth/...

# Run only unit tests (skip integration)
go test -short ./...

# Generate HTML coverage report
make test-coverage-html

# Run benchmarks
make test-bench

# Run with race detection
go test -race ./...
```

## Important Notes

### Type Compatibility

1. **OpenAPI Generated Types**: Always use the exact generated types from `types.gen.go`:
   - `Email` field is `openapi_types.Email`, not `string`
   - Role enums use specific constants like `CreateMemberJSONBodyRoleAdmin`
   - Request/Response types use generated structs

2. **UUID Handling**:
   - User `ID` field is `uuid.UUID`, not `ID`
   - Some fields like `OrganizationID` may be strings in certain contexts

3. **Context Keys**:
   - Auth middleware uses `AuthUserContextKey` and `AuthClaimsContextKey`
   - Not `UserContextKey` or `ClaimsContextKey`

### Mock Generation Checklist

When creating mocks for testing:

1. Never create manual mocks - always use mockery
2. Define interfaces first, then generate mocks
3. Run `make generate-mocks` after interface changes
4. Use generated mocks in `mocks_test.go` files
5. Mock files follow pattern `Mock{InterfaceName}`

### Common Pitfalls

1. **Forgetting Generated Types**: Always check `types.gen.go` for correct field types
2. **Missing Interface Methods**: Use `var _ Interface = (*Mock)(nil)` to verify
3. **Incorrect Error Variables**: Use domain-specific errors like `ErrUserNotFound`
4. **Linting Issues**: Run `make lint` before committing
5. **Format Issues**: Run `gofmt -w` on test files

## Test Data Fixtures

Common test data patterns:

```go
// User fixture
user := &User{
    ID:            uuid.New(),
    Email:         "test@example.com",
    Name:          "Test User",
    EmailVerified: false,
    CreatedAt:     time.Now(),
    UpdatedAt:     time.Now(),
}

// Organization fixture
org := &Organization{
    ID:           uuid.New(),
    Name:         "Test Org",
    BillingEmail: "billing@example.com",
    Plan:         OrganizationPlan(DefaultPlan),
}
```

## Coverage Goals

### Priority 1 (Business Critical)

- [ ] Auth domain - Current: 20.2%, Target: 80%
- [ ] Organizations domain - Current: 3.3%, Target: 70%
- [ ] Pipelines domain - Current: 12.8%, Target: 70%
- [ ] Artifacts domain - Current: 17.2%, Target: 70%
- [ ] Tools domain - Current: 0.0%, Target: 60%
- [ ] Users domain - Current: 12.2%, Target: 60%

### Priority 2 (Infrastructure)

- [x] Config package - Current: 47.2%, Target: 50%
- [ ] Sessions package - Current: 27.2%, Target: 60%
- [ ] Database package - Target: 60%
- [ ] Storage package - Target: 60%

### Priority 3 (Supporting)

- [ ] Server package - Target: 40%
- [ ] Health checks - Target: 50%
- [ ] CLI tools - Target: 30%

## Future Improvements

1. **E2E Tests**: Add end-to-end API tests with full server startup
2. **Performance Tests**: Add benchmarks for critical paths
3. **Fuzz Testing**: Add fuzzing for input validation
4. **Contract Tests**: Ensure API compatibility with OpenAPI spec
5. **Load Tests**: Add k6 or similar for load testing
6. **Mutation Testing**: Consider adding mutation testing for test quality

## Contributing

When adding new tests:

1. Follow established patterns in existing test files
2. Include both success and error cases
3. Use descriptive test names
4. Update this documentation if introducing new patterns
5. Ensure all tests pass with `make test`
6. Check coverage with `make test-coverage`
