package pipelines

import (
	"context"
	"fmt"
	"log/slog"
)

// Handler implements StrictServerInterface by delegating to the service.
type Handler struct {
	service ServiceInterface
	logger  *slog.Logger
}

// NewHandler creates a new handler that implements StrictServerInterface.
func NewHandler(service ServiceInterface, logger *slog.Logger) StrictServerInterface {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// CreatePipeline handles the create pipeline endpoint.
func (h *Handler) CreatePipeline(
	ctx context.Context,
	request CreatePipelineRequestObject,
) (CreatePipelineResponseObject, error) {
	return h.service.Create(ctx, request)
}

// GetPipeline handles the get pipeline endpoint.
func (h *Handler) GetPipeline(
	ctx context.Context,
	request GetPipelineRequestObject,
) (GetPipelineResponseObject, error) {
	return h.service.Get(ctx, request)
}

// UpdatePipeline handles the update pipeline endpoint.
func (h *Handler) UpdatePipeline(
	ctx context.Context,
	request UpdatePipelineRequestObject,
) (UpdatePipelineResponseObject, error) {
	return h.service.Update(ctx, request)
}

// DeletePipeline handles the delete pipeline endpoint.
func (h *Handler) DeletePipeline(
	ctx context.Context,
	request DeletePipelineRequestObject,
) (DeletePipelineResponseObject, error) {
	return h.service.Delete(ctx, request)
}

// ListPipelines handles the list pipelines endpoint.
func (h *Handler) ListPipelines(
	ctx context.Context,
	request ListPipelinesRequestObject,
) (ListPipelinesResponseObject, error) {
	return h.service.List(ctx, request)
}

// GetPipelineExecutionPlan handles the get pipeline execution plan endpoint.
func (h *Handler) GetPipelineExecutionPlan(
	_ context.Context,
	_ GetPipelineExecutionPlanRequestObject,
) (GetPipelineExecutionPlanResponseObject, error) {
	// TODO: Implement execution plan logic
	return nil, fmt.Errorf("not implemented")
}

// ValidatePipelineExecutionPlan handles the validate pipeline execution plan endpoint.
func (h *Handler) ValidatePipelineExecutionPlan(
	_ context.Context,
	_ ValidatePipelineExecutionPlanRequestObject,
) (ValidatePipelineExecutionPlanResponseObject, error) {
	// TODO: Implement validation logic
	return nil, fmt.Errorf("not implemented")
}

// GetPipelineSteps handles the get pipeline steps endpoint.
func (h *Handler) GetPipelineSteps(
	_ context.Context,
	_ GetPipelineStepsRequestObject,
) (GetPipelineStepsResponseObject, error) {
	// TODO: Implement get steps logic
	return nil, fmt.Errorf("not implemented")
}

// CreatePipelineStep handles the create pipeline step endpoint.
func (h *Handler) CreatePipelineStep(
	_ context.Context,
	_ CreatePipelineStepRequestObject,
) (CreatePipelineStepResponseObject, error) {
	// TODO: Implement create step logic
	return nil, fmt.Errorf("not implemented")
}
