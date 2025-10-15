package executor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/infrastructure/executor"
)

// TestLocalExecutorBasic tests basic local execution
func TestLocalExecutorBasic(t *testing.T) {
	type Input struct {
		Value int `json:"value"`
	}
	type Output struct {
		Doubled int `json:"doubled"`
	}

	// Define execution function
	executeFunc := func(_ context.Context, input Input) (Output, error) {
		return Output{Doubled: input.Value * 2}, nil
	}

	// Create local executor
	exec, err := executor.NewLocalExecutor[Input, Output](
		executor.LocalConfig{},
		executeFunc,
	)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Execute
	ctx := context.Background()
	output, err := exec.Execute(ctx, Input{Value: 21})
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Verify
	if output.Doubled != 42 {
		t.Errorf("Expected doubled=42, got %d", output.Doubled)
	}
}

// TestLocalExecutorNilFunction tests that nil function is rejected
func TestLocalExecutorNilFunction(t *testing.T) {
	type Input struct{}
	type Output struct{}

	_, err := executor.NewLocalExecutor[Input, Output](
		executor.LocalConfig{}, nil,
	)
	if err == nil {
		t.Fatal("Expected error for nil function")
	}
}

// TestLocalExecutorSchemaValidation tests schema validation
func TestLocalExecutorSchemaValidation(t *testing.T) {
	type Input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type Output struct {
		Message string `json:"message"`
	}

	// Define execution function
	executeFunc := func(_ context.Context, input Input) (Output, error) {
		return Output{Message: "Hello, " + input.Name}, nil
	}

	// Define strict schemas
	inputSchema := []byte(`{
		"type": "object",
		"required": ["name", "email"],
		"properties": {
			"name": {"type": "string", "minLength": 1},
			"email": {"type": "string", "format": "email"}
		},
		"additionalProperties": false
	}`)

	outputSchema := []byte(`{
		"type": "object",
		"required": ["message"],
		"properties": {
			"message": {"type": "string"}
		}
	}`)

	// Create executor with schemas
	exec, err := executor.NewLocalExecutor[Input, Output](
		executor.LocalConfig{
			Config: executor.Config{
				SchemaIn:  inputSchema,
				SchemaOut: outputSchema,
			},
		},
		executeFunc,
	)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	ctx := context.Background()

	// Test with valid input
	output, err := exec.Execute(ctx, Input{
		Name:  "John Doe",
		Email: "john@example.com",
	})
	if err != nil {
		t.Fatalf("Execution with valid input failed: %v", err)
	}
	if output.Message != "Hello, John Doe" {
		t.Errorf("Unexpected output: %s", output.Message)
	}

	// Test with invalid input (empty name)
	_, err = exec.Execute(ctx, Input{
		Name:  "",
		Email: "john@example.com",
	})
	if err == nil {
		t.Fatal("Expected validation error for empty name")
	}
	t.Logf("Got expected validation error: %v", err)
}

// TestLocalExecutorInvalidOutputSchema tests output validation failure
func TestLocalExecutorInvalidOutputSchema(t *testing.T) {
	type Input struct {
		Value int `json:"value"`
	}
	type Output struct {
		Result string `json:"result"`
	}

	// Function that returns invalid output (empty string)
	executeFunc := func(_ context.Context, _ Input) (Output, error) {
		return Output{Result: ""}, nil
	}

	outputSchema := []byte(`{
		"type": "object",
		"required": ["result"],
		"properties": {
			"result": {"type": "string", "minLength": 1}
		}
	}`)

	exec, err := executor.NewLocalExecutor[Input, Output](
		executor.LocalConfig{
			Config: executor.Config{
				SchemaIn:  []byte(`{"type": "object"}`),
				SchemaOut: outputSchema,
			},
		},
		executeFunc,
	)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	ctx := context.Background()
	_, err = exec.Execute(ctx, Input{Value: 1})
	if err == nil {
		t.Fatal("Expected output validation error")
	}
	t.Logf("Got expected output validation error: %v", err)
}

// TestLocalExecutorTimeout tests timeout handling
func TestLocalExecutorTimeout(t *testing.T) {
	type Input struct {
		Data string `json:"data"`
	}
	type Output struct {
		Result string `json:"result"`
	}

	// Function that sleeps longer than timeout
	executeFunc := func(ctx context.Context, _ Input) (Output, error) {
		select {
		case <-time.After(5 * time.Second):
			return Output{Result: "completed"}, nil
		case <-ctx.Done():
			return Output{}, ctx.Err()
		}
	}

	// Create executor with short timeout
	exec, err := executor.NewLocalExecutor[Input, Output](
		executor.LocalConfig{
			Config: executor.Config{
				Timeout: 100 * time.Millisecond,
			},
		},
		executeFunc,
	)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	ctx := context.Background()
	_, err = exec.Execute(ctx, Input{Data: "test"})
	if err == nil {
		t.Fatal("Expected timeout error")
	}
	t.Logf("Got expected timeout error: %v", err)
}

// TestLocalExecutorError tests error propagation
func TestLocalExecutorError(t *testing.T) {
	type Input struct {
		Value int `json:"value"`
	}
	type Output struct {
		Result int `json:"result"`
	}

	expectedErr := errors.New("custom error")

	// Function that returns an error
	executeFunc := func(_ context.Context, _ Input) (Output, error) {
		return Output{}, expectedErr
	}

	exec, err := executor.NewLocalExecutor[Input, Output](
		executor.LocalConfig{},
		executeFunc,
	)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	ctx := context.Background()
	_, err = exec.Execute(ctx, Input{Value: 1})
	if err == nil {
		t.Fatal("Expected error from execution")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error to wrap custom error, got: %v", err)
	}
}

// TestLocalExecutorContextCancellation tests context cancellation
func TestLocalExecutorContextCancellation(t *testing.T) {
	type Input struct {
		Data string `json:"data"`
	}
	type Output struct {
		Result string `json:"result"`
	}

	// Function that respects context
	executeFunc := func(ctx context.Context, _ Input) (Output, error) {
		select {
		case <-time.After(5 * time.Second):
			return Output{Result: "completed"}, nil
		case <-ctx.Done():
			return Output{}, ctx.Err()
		}
	}

	exec, err := executor.NewLocalExecutor[Input, Output](
		executor.LocalConfig{
			Config: executor.Config{
				Timeout: 10 * time.Second, // Long timeout
			},
		},
		executeFunc,
	)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	_, err = exec.Execute(ctx, Input{Data: "test"})
	if err == nil {
		t.Fatal("Expected cancellation error")
	}
	t.Logf("Got expected cancellation error: %v", err)
}

// BenchmarkLocalExecution benchmarks local execution
func BenchmarkLocalExecution(b *testing.B) {
	type Input struct {
		Values []int `json:"values"`
	}
	type Output struct {
		Sum int `json:"sum"`
	}

	executeFunc := func(_ context.Context, input Input) (Output, error) {
		sum := 0
		for _, v := range input.Values {
			sum += v
		}
		return Output{Sum: sum}, nil
	}

	exec, err := executor.NewLocalExecutor[Input, Output](
		executor.LocalConfig{},
		executeFunc,
	)
	if err != nil {
		b.Fatalf("Failed to create executor: %v", err)
	}

	ctx := context.Background()
	input := Input{Values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := exec.Execute(ctx, input)
		if err != nil {
			b.Fatalf("Execution failed: %v", err)
		}
	}
}

// BenchmarkLocalExecutionWithValidation benchmarks local execution with schema validation
func BenchmarkLocalExecutionWithValidation(b *testing.B) {
	type Input struct {
		Values []int `json:"values"`
	}
	type Output struct {
		Sum int `json:"sum"`
	}

	executeFunc := func(_ context.Context, input Input) (Output, error) {
		sum := 0
		for _, v := range input.Values {
			sum += v
		}
		return Output{Sum: sum}, nil
	}

	inputSchema := []byte(`{
		"type": "object",
		"required": ["values"],
		"properties": {
			"values": {
				"type": "array",
				"items": {"type": "integer"}
			}
		}
	}`)

	outputSchema := []byte(`{
		"type": "object",
		"required": ["sum"],
		"properties": {
			"sum": {"type": "integer"}
		}
	}`)

	exec, err := executor.NewLocalExecutor[Input, Output](
		executor.LocalConfig{
			Config: executor.Config{
				SchemaIn:  inputSchema,
				SchemaOut: outputSchema,
			},
		},
		executeFunc,
	)
	if err != nil {
		b.Fatalf("Failed to create executor: %v", err)
	}

	ctx := context.Background()
	input := Input{Values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := exec.Execute(ctx, input)
		if err != nil {
			b.Fatalf("Execution failed: %v", err)
		}
	}
}
