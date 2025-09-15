package tools

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service implements the Service interface for tool operations
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new tool service
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves tools for an organization
func (s *Service) List(ctx context.Context, _ string, limit, offset int) ([]*Tool, int64, error) {
	params := ListToolsParams{
		Page: PageQuery{
			Number: offset/limit + 1,
			Size:   limit,
		},
	}

	tools, total, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.Error("failed to list tools", "error", err)
		return nil, 0, err
	}

	return tools, total, nil
}

// Create creates a new tool
func (s *Service) Create(ctx context.Context, req *CreateToolJSONRequestBody, orgID UUID) (*Tool, error) {
	tool := &Tool{
		Id:             uuid.New(),
		OrganizationId: orgID,
		Name:           req.Name,
		Description:    req.Description,
	}

	createdTool, err := s.repo.Create(ctx, tool)
	if err != nil {
		s.logger.Error("failed to create tool", "error", err)
		return nil, err
	}

	return createdTool, nil
}

// Get retrieves a tool by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Tool, error) {
	tool, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get tool", "id", id, "error", err)
		return nil, err
	}

	return tool, nil
}

// Update updates a tool
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateToolJSONRequestBody) (*Tool, error) {
	tool, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get tool for update", "id", id, "error", err)
		return nil, err
	}

	if req.Name != "" {
		tool.Name = req.Name
	}
	if req.Description != "" {
		tool.Description = req.Description
	}

	updatedTool, err := s.repo.Update(ctx, id, tool)
	if err != nil {
		s.logger.Error("failed to update tool", "id", id, "error", err)
		return nil, err
	}

	return updatedTool, nil
}

// Delete deletes a tool by ID
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete tool", "id", id, "error", err)
		return err
	}

	return nil
}

// Domain type aliases
type (
	// CreateToolRequest represents a request to create a tool
	CreateToolRequest = CreateToolJSONBody

	// UpdateToolRequest represents a request to update a tool
	UpdateToolRequest = UpdateToolJSONBody
)
