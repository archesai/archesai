# Testing Standards

This document defines the testing standards, patterns, and best practices for the Arches codebase.

## Table of Contents

- [Test Organization](#test-organization)
- [Table-Driven Tests](#table-driven-tests)
- [Test Naming](#test-naming)
- [Assertions](#assertions)
- [Test Helpers](#test-helpers)
- [Mocking](#mocking)
- [Test Fixtures](#test-fixtures)
- [Error Testing](#error-testing)
- [Coverage](#coverage)

---

## Test Organization

### File Structure

Test files are placed alongside the code they test:

```text
pkg/validation/
    errors.go
    errors_test.go
    rules.go
    rules_test.go
    validator.go
```

### Test File Naming

- Test files must end with `_test.go`
- Name matches the file being tested: `foo.go` → `foo_test.go`

### Package Declaration

Use the same package name for white-box testing (access to unexported members):

```go
package validation

import "testing"

func TestErrors_Add(t *testing.T) { ... }
```

Use `_test` suffix for black-box testing (only exported members):

```go
package validation_test

import (
    "testing"
    "github.com/archesai/archesai/pkg/validation"
)

func TestValidateStruct(t *testing.T) { ... }
```

---

## Table-Driven Tests

**All tests should use table-driven format.** This is the standard pattern for the codebase.

### Standard Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
    }{
        {"descriptive name", input1, expected1},
        {"another case", input2, expected2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionName(tt.input)
            if result != tt.expected {
                t.Errorf("FunctionName(%v) = %v, want %v", tt.input, result, tt.expected)
            }
        })
    }
}
```

### With Error Expectations

```go
func TestRequired(t *testing.T) {
    tests := []struct {
        name      string
        value     *string
        wantError bool
    }{
        {"nil value", nil, true},
        {"empty string", ptr(""), true},
        {"valid string", ptr("hello"), false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var errs Errors
            Required(tt.value, "field", &errs)
            if tt.wantError && !errs.HasErrors() {
                t.Error("expected error, got none")
            }
            if !tt.wantError && errs.HasErrors() {
                t.Errorf("unexpected error: %v", errs)
            }
        })
    }
}
```

### With Multiple Inputs

```go
func TestMinLength(t *testing.T) {
    tests := []struct {
        name      string
        value     *string
        min       int
        wantError bool
    }{
        {"nil value", nil, 3, false},
        {"too short", ptr("ab"), 3, true},
        {"exact length", ptr("abc"), 3, false},
        {"longer", ptr("abcdef"), 3, false},
        {"unicode", ptr("日本"), 3, true},
        {"unicode valid", ptr("日本語"), 3, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var errs Errors
            MinLength(tt.value, tt.min, "field", &errs)
            if tt.wantError && !errs.HasErrors() {
                t.Error("expected error, got none")
            }
            if !tt.wantError && errs.HasErrors() {
                t.Errorf("unexpected error: %v", errs)
            }
        })
    }
}
```

### Complex Test Cases

For tests with complex setup or expected values, use named struct fields:

```go
func TestErrors_Error(t *testing.T) {
    tests := []struct {
        name   string
        errors Errors
        want   string
    }{
        {
            name:   "empty errors",
            errors: Errors{},
            want:   "",
        },
        {
            name: "single error",
            errors: Errors{
                {Field: "name", Message: "is required"},
            },
            want: "validation failed: name is required",
        },
        {
            name: "multiple errors",
            errors: Errors{
                {Field: "name", Message: "is required"},
                {Field: "email", Message: "must be a valid email address"},
            },
            want: "validation failed: name is required, email must be a valid email address",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.errors.Error(); got != tt.want {
                t.Errorf("Errors.Error() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## Test Naming

### Test Function Names

Format: `Test<Type>_<Method>` or `Test<Function>`

```go
// Testing a method on a type
func TestErrors_Error(t *testing.T) { ... }
func TestErrors_HasErrors(t *testing.T) { ... }
func TestErrors_Add(t *testing.T) { ... }

// Testing a standalone function
func TestKebabCase(t *testing.T) { ... }
func TestValidateStruct(t *testing.T) { ... }

// Testing a complex scenario
func TestParser_Parse_BasicSpec(t *testing.T) { ... }
func TestParser_Parse_AutoDiscoveryPaths(t *testing.T) { ... }
```

### Subtest Names (t.Run)

Use lowercase, descriptive names:

```go
// Good - descriptive, lowercase
t.Run("nil value", func(t *testing.T) { ... })
t.Run("empty string", func(t *testing.T) { ... })
t.Run("valid email", func(t *testing.T) { ... })
t.Run("acronym at start", func(t *testing.T) { ... })

// Bad - too vague or wrong format
t.Run("test1", func(t *testing.T) { ... })
t.Run("TestNilValue", func(t *testing.T) { ... })
```

---

## Assertions

### Standard Assertion Pattern

Use the standard library's `testing` package. No external assertion libraries.

```go
// Good - clear error message with got/want pattern
if got := result; got != tt.want {
    t.Errorf("Function() = %v, want %v", got, tt.want)
}

// Good - with input context
if result != tt.expected {
    t.Errorf("KebabCase(%q) = %q, want %q", tt.input, result, tt.expected)
}
```

### Error Assertions

```go
// Checking for expected error
if tt.wantError && err == nil {
    t.Error("expected error, got none")
}
if !tt.wantError && err != nil {
    t.Errorf("unexpected error: %v", err)
}

// Checking error absence
if err != nil {
    t.Fatalf("Function() error = %v", err)
}
```

### Fatal vs Error

- Use `t.Fatal` or `t.Fatalf` when the test cannot continue
- Use `t.Error` or `t.Errorf` when the test can continue with other checks

```go
// Fatal - test cannot continue without this
doc, err := NewOpenAPIDocumentFromFS(fsys, "openapi.yaml")
if err != nil {
    t.Fatalf("NewOpenAPIDocumentFromFS() error = %v", err)
}

// Error - test can check other things
if s.ProjectName != "github.com/example/myapi" {
    t.Errorf("ProjectName = %q, want %q", s.ProjectName, "github.com/example/myapi")
}
if s.Title != "My API" {
    t.Errorf("Title = %q, want %q", s.Title, "My API")
}
```

### Checking Multiple Conditions

```go
// Good - check each condition separately for clear failure messages
if len(s.Tags) != 1 {
    t.Errorf("len(Tags) = %d, want 1", len(s.Tags))
} else if s.Tags[0].Name != "User" {
    t.Errorf("Tags[0].Name = %q, want %q", s.Tags[0].Name, "User")
}
```

---

## Test Helpers

### Helper Functions

Mark helper functions with `t.Helper()`:

```go
// Good - helper marked properly
func newTestParser(t *testing.T, files map[string]string) *Parser {
    t.Helper()
    fsys := fstest.MapFS{}
    for path, content := range files {
        fsys[path] = &fstest.MapFile{Data: []byte(content)}
    }
    doc, err := NewOpenAPIDocumentFromFS(fsys, "openapi.yaml")
    if err != nil {
        t.Fatalf("NewOpenAPIDocumentFromFS() error = %v", err)
    }
    return NewParser(doc)
}
```

### Pointer Helper

A common helper for creating pointers to values:

```go
func ptr[T any](v T) *T {
    return &v
}

// Usage
tests := []struct {
    name  string
    value *string
}{
    {"nil value", nil},
    {"empty string", ptr("")},
    {"valid string", ptr("hello")},
}
```

### Test Data Builders

For complex test data, use inline maps or structs:

```go
func TestParser_Parse_BasicSpec(t *testing.T) {
    files := map[string]string{
        "openapi.yaml": `
openapi: 3.1.0
x-project-name: github.com/example/myapi
info:
  title: My API
  version: v1.0.0
`,
        "components/schemas/User.yaml": `
title: User
type: object
properties:
  id:
    type: string
    format: uuid
`,
    }

    p := newTestParser(t, files)
    // ... test logic
}
```

---

## Mocking

### Use Mockery

**Never create manual mocks.** Always use mockery to generate mocks from interfaces.

```go
//go:generate mockery --name=ExecutorRepository
```

### Mock Usage Pattern

```go
func TestService_Execute(t *testing.T) {
    mockRepo := mocks.NewExecutorRepository(t)
    mockRepo.On("Get", mock.Anything, expectedID).Return(executor, nil)

    service := NewService(mockRepo)
    result, err := service.Execute(ctx, executorID, input)

    mockRepo.AssertExpectations(t)
    // ... assertions
}
```

---

## Test Fixtures

### Inline Test Data

Prefer inline test data for clarity:

```go
files := map[string]string{
    "openapi.yaml": `
openapi: 3.1.0
info:
  title: My API
  version: v1.0.0
`,
}
```

### Using testdata Directory

For larger fixtures, use the `testdata` directory:

```text
pkg/spec/
    parse.go
    parse_test.go
    testdata/
        basic_spec.yaml
        complex_spec.yaml
```

```go
func TestParser_ParseFile(t *testing.T) {
    data, err := os.ReadFile("testdata/basic_spec.yaml")
    if err != nil {
        t.Fatalf("failed to read test data: %v", err)
    }
    // ... test logic
}
```

### Using testing/fstest

For filesystem-based tests:

```go
import "testing/fstest"

func TestFunction(t *testing.T) {
    fsys := fstest.MapFS{
        "file.yaml": &fstest.MapFile{Data: []byte("content")},
        "dir/nested.yaml": &fstest.MapFile{Data: []byte("nested content")},
    }
    // ... test logic using fsys
}
```

---

## Error Testing

### Testing Error Conditions

```go
func TestFunction_Errors(t *testing.T) {
    tests := []struct {
        name      string
        input     InputType
        wantError bool
    }{
        {"valid input", validInput, false},
        {"invalid input", invalidInput, true},
        {"nil input", nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := Function(tt.input)
            if tt.wantError && err == nil {
                t.Error("expected error, got none")
            }
            if !tt.wantError && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

### Testing Specific Error Types

```go
func TestFunction_ErrorTypes(t *testing.T) {
    _, err := Function(invalidInput)

    var badReq BadRequestError
    if !errors.As(err, &badReq) {
        t.Errorf("expected BadRequestError, got %T", err)
    }
}
```

### Testing Error Messages

```go
func TestErrors_Error(t *testing.T) {
    errs := Errors{
        {Field: "name", Message: "is required"},
    }

    got := errs.Error()
    want := "validation failed: name is required"

    if got != want {
        t.Errorf("Error() = %q, want %q", got, want)
    }
}
```

---

## Coverage

### Running Tests with Coverage

```bash
# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Check coverage percentage
go test -cover ./... | grep -E "coverage:|ok"
```

### Coverage Goals

- Aim for high coverage on business logic
- Focus on testing edge cases and error conditions
- Don't write tests just to increase coverage numbers

### What to Test

**Must test:**

- Public APIs
- Error handling paths
- Edge cases (nil, empty, boundary values)
- Business logic

**Optional:**

- Simple getters/setters
- Trivial constructors
- Generated code (test the generator, not the output)

---

## Common Patterns

### Setup and Teardown

For tests that need setup/teardown:

```go
func TestWithSetup(t *testing.T) {
    // Setup
    tmpDir, err := os.MkdirTemp("", "test-*")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tmpDir)

    // Test logic
    // ...
}
```

### Parallel Tests

Mark independent tests as parallel when appropriate:

```go
func TestIndependentCases(t *testing.T) {
    tests := []struct {
        name string
        // ...
    }{
        // ...
    }

    for _, tt := range tests {
        tt := tt // capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // test logic
        })
    }
}
```

### Testing with Context

```go
func TestServiceMethod(t *testing.T) {
    ctx := context.Background()
    // or with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := service.Method(ctx, input)
    // ... assertions
}
```

---

## Anti-Patterns to Avoid

| Anti-Pattern                   | Problem                   | Solution                                |
| ------------------------------ | ------------------------- | --------------------------------------- |
| Hard-coded expected values     | Brittle, unclear intent   | Use named constants or compute expected |
| Testing implementation details | Breaks on refactor        | Test behavior, not implementation       |
| Large test functions           | Hard to maintain          | Break into table-driven subtests        |
| Shared mutable state           | Flaky tests               | Isolate test data per test              |
| Sleep-based synchronization    | Slow, flaky               | Use channels or sync primitives         |
| Ignoring errors in tests       | Hides bugs                | Always check errors                     |
| Manual mocks                   | Inconsistent, error-prone | Use mockery                             |

```go
// Bad - hard-coded magic value
if len(result) != 42 {
    t.Error("wrong length")
}

// Good - explain the expectation
expectedCount := len(inputItems) * 2 // each item produces two outputs
if len(result) != expectedCount {
    t.Errorf("len(result) = %d, want %d", len(result), expectedCount)
}
```
