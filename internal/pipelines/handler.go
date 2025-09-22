package pipelines

import (
	"context"
	"fmt"
)

// CreatePipelineStep handles the create pipeline step endpoint
func (s *Service) CreatePipelineStep(
	_ context.Context,
	_ CreatePipelineStepRequestObject,
) (CreatePipelineStepResponseObject, error) {
	return nil, fmt.Errorf("CreatePipelineStep not yet implemented")
}

// GetPipelineExecutionPlan handles the create pipeline endpoint
func (s *Service) GetPipelineExecutionPlan(
	_ context.Context,
	_ GetPipelineExecutionPlanRequestObject,
) (GetPipelineExecutionPlanResponseObject, error) {
	return nil, fmt.Errorf("CreatePipelineStep not yet implemented")
}

// GetPipelineSteps handles the get pipeline steps endpoint
func (s *Service) GetPipelineSteps(
	_ context.Context,
	_ GetPipelineStepsRequestObject,
) (GetPipelineStepsResponseObject, error) {
	return nil, fmt.Errorf("GetPipelineSteps not yet implemented")
}

// ValidatePipelineExecutionPlan handles the validate pipeline execution plan endpoint
func (s *Service) ValidatePipelineExecutionPlan(
	_ context.Context,
	_ ValidatePipelineExecutionPlanRequestObject,
) (ValidatePipelineExecutionPlanResponseObject, error) {
	return nil, fmt.Errorf("ValidatePipelineExecutionPlan not yet implemented")
}
