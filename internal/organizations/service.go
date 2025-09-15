package organizations

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Service provides organization business logic
type Service struct {
	repository Repository
	logger     *slog.Logger
}

// NewService creates a new organization service
func NewService(repository Repository, logger *slog.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

// Create creates a new organization
func (s *Service) Create(ctx context.Context, req *CreateOrganizationRequest, creatorUserID string) (*Organization, error) {
	s.logger.Debug("creating organization", "id", req.OrganizationId, "creator", creatorUserID)

	// Set default plan
	plan := OrganizationPlan(DefaultPlan)

	org := &Organization{
		Id:           req.OrganizationId,
		Name:         "", // Name should be set from somewhere else
		BillingEmail: openapi_types.Email(req.BillingEmail),
		Plan:         plan,
		Credits:      0.0, // Start with 0 credits
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdOrg, err := s.repository.Create(ctx, org)
	if err != nil {
		s.logger.Error("failed to create organization", "error", err)
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	s.logger.Info("organization created successfully", "id", createdOrg.Id, "name", createdOrg.Name)
	return createdOrg, nil
}

// Get retrieves an organization by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Organization, error) {
	org, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return org, nil
}

// Update updates an organization
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateOrganizationRequest) (*Organization, error) {
	org, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Update fields that were provided
	if req.BillingEmail != "" {
		org.BillingEmail = openapi_types.Email(req.BillingEmail)
	}
	org.UpdatedAt = time.Now()

	updatedOrg, err := s.repository.Update(ctx, org.Id, org)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return updatedOrg, nil
}

// Delete deletes an organization
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Add additional checks (e.g., organization has no active resources)
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}

// List retrieves a list of organizations
func (s *Service) List(ctx context.Context, limit, offset int) ([]*Organization, int, error) {
	orgs, totalInt64, err := s.repository.List(ctx, ListOrganizationsParams{
		Page: PageQuery{
			Number: offset/limit + 1,
			Size:   limit,
		},
	})
	total := int(totalInt64)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list organizations: %w", err)
	}
	return orgs, total, nil
}
