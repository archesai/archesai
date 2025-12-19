package executor

import "context"

// Executor is the core abstraction for business logic operations.
// Every operation (CreatePipeline, GetUser, etc.) implements this interface
// with typed input and output.
//
// Example usage:
//
//	// Generated type alias in your package
//	type CreatePipeline = executor.Executor[CreatePipelineInput, CreatePipelineOutput]
//
//	// Your implementation (unexported struct)
//	type createPipeline struct {
//	    repo PipelineRepository
//	}
//
//	// Exported constructor returning concrete type
//	func NewCreatePipeline(repo PipelineRepository) *createPipeline {
//	    return &createPipeline{repo: repo}
//	}
//
//	// Implement the interface
//	func (c *createPipeline) Execute(ctx context.Context, input *CreatePipelineInput) (*CreatePipelineOutput, error) {
//	    // business logic here
//	}
//
// Callers assign to the interface type:
//
//	var uc pipelines.CreatePipeline = pipelines.NewCreatePipeline(repo)
type Executor[I, O any] interface {
	Execute(ctx context.Context, input *I) (*O, error)
}

// VoidExecutor is used for operations that return no data (HTTP 204 No Content).
// Typically used for DELETE operations.
//
// Example usage:
//
//	// Generated type alias
//	type DeletePipeline = executor.VoidExecutor[DeletePipelineInput]
//
//	// Implementation
//	type deletePipeline struct {
//	    repo PipelineRepository
//	}
//
//	func NewDeletePipeline(repo PipelineRepository) *deletePipeline {
//	    return &deletePipeline{repo: repo}
//	}
//
//	func (d *deletePipeline) Execute(ctx context.Context, input *DeletePipelineInput) error {
//	    return d.repo.Delete(ctx, input.ID)
//	}
type VoidExecutor[I any] interface {
	Execute(ctx context.Context, input *I) error
}
