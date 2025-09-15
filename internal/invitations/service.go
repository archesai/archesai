package invitations

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Invitation status constants
const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusDeclined = "declined"
)

// InvitationService handles business logic for invitations
type InvitationService interface {
	Create(ctx context.Context, invitation *Invitation) (*Invitation, error)
	Get(ctx context.Context, id uuid.UUID) (*Invitation, error)
	GetByEmail(ctx context.Context, email string, organizationID string) (*Invitation, error)
	List(ctx context.Context, params ListInvitationsParams) ([]*Invitation, int64, error)
	ListByOrganization(ctx context.Context, organizationID string) ([]*Invitation, error)
	ListByInviter(ctx context.Context, inviterID string) ([]*Invitation, error)
	Update(ctx context.Context, id uuid.UUID, invitation *Invitation) (*Invitation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Accept(ctx context.Context, id uuid.UUID) (*Invitation, error)
	Decline(ctx context.Context, id uuid.UUID) (*Invitation, error)
	Resend(ctx context.Context, id uuid.UUID) (*Invitation, error)
}

// Service implements the InvitationService interface
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new invitation service
func NewService(repo Repository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new invitation
func (s *Service) Create(ctx context.Context, invitation *Invitation) (*Invitation, error) {
	// Check if invitation already exists for this email in the organization
	existing, err := s.repo.GetByEmail(ctx, invitation.Email, invitation.OrganizationId)
	if err == nil && existing != nil {
		if existing.Status == StatusPending {
			return nil, ErrInvitationAlreadyExists
		}
	}

	// Set default values
	if invitation.Id == uuid.Nil {
		invitation.Id = uuid.New()
	}
	invitation.Status = StatusPending
	invitation.CreatedAt = time.Now()
	invitation.UpdatedAt = time.Now()

	// Set expiration to 7 days from now if not set
	if invitation.ExpiresAt == "" {
		invitation.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	}

	created, err := s.repo.Create(ctx, invitation)
	if err != nil {
		s.logger.Error("failed to create invitation",
			"error", err,
			"email", invitation.Email,
			"organizationId", invitation.OrganizationId)
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	return created, nil
}

// Get retrieves an invitation by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Invitation, error) {
	invitation, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Debug("invitation not found", "id", id, "error", err)
		return nil, ErrInvitationNotFound
	}
	return invitation, nil
}

// GetByEmail retrieves an invitation by email and organization ID
func (s *Service) GetByEmail(ctx context.Context, email string, organizationID string) (*Invitation, error) {
	invitation, err := s.repo.GetByEmail(ctx, email, organizationID)
	if err != nil {
		return nil, ErrInvitationNotFound
	}
	return invitation, nil
}

// List retrieves invitations with pagination
func (s *Service) List(ctx context.Context, params ListInvitationsParams) ([]*Invitation, int64, error) {
	// Set defaults
	if params.Page.Size == 0 {
		params.Page.Size = 10
	}
	if params.Page.Size > 100 {
		params.Page.Size = 100
	}

	invitations, total, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.Error("failed to list invitations", "error", err)
		return nil, 0, fmt.Errorf("failed to list invitations: %w", err)
	}

	return invitations, total, nil
}

// ListByOrganization retrieves all invitations for an organization
func (s *Service) ListByOrganization(ctx context.Context, organizationID string) ([]*Invitation, error) {
	invitations, err := s.repo.ListByOrganization(ctx, organizationID)
	if err != nil {
		s.logger.Error("failed to list invitations by organization",
			"error", err,
			"organizationId", organizationID)
		return nil, fmt.Errorf("failed to list invitations by organization: %w", err)
	}
	return invitations, nil
}

// ListByInviter retrieves all invitations sent by a specific user
func (s *Service) ListByInviter(ctx context.Context, inviterID string) ([]*Invitation, error) {
	invitations, err := s.repo.ListByInviter(ctx, inviterID)
	if err != nil {
		s.logger.Error("failed to list invitations by inviter",
			"error", err,
			"inviterId", inviterID)
		return nil, fmt.Errorf("failed to list invitations by inviter: %w", err)
	}
	return invitations, nil
}

// Update updates an invitation
func (s *Service) Update(ctx context.Context, id uuid.UUID, invitation *Invitation) (*Invitation, error) {
	// Get existing invitation
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, ErrInvitationNotFound
	}

	// Check if invitation can be updated
	if existing.Status == StatusAccepted {
		return nil, ErrInvitationAlreadyAccepted
	}
	if existing.Status == StatusDeclined {
		return nil, ErrInvitationAlreadyDeclined
	}

	// Check if expired
	expiresAt, err := time.Parse(time.RFC3339, existing.ExpiresAt)
	if err == nil && time.Now().After(expiresAt) {
		return nil, ErrInvitationExpired
	}

	// Update fields
	invitation.UpdatedAt = time.Now()

	updated, err := s.repo.Update(ctx, id, invitation)
	if err != nil {
		s.logger.Error("failed to update invitation", "error", err, "id", id)
		return nil, fmt.Errorf("failed to update invitation: %w", err)
	}

	return updated, nil
}

// Delete deletes an invitation
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if invitation exists
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return ErrInvitationNotFound
	}

	// Only pending invitations can be deleted
	if existing.Status != StatusPending {
		return fmt.Errorf("cannot delete invitation with status %s", existing.Status)
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete invitation", "error", err, "id", id)
		return fmt.Errorf("failed to delete invitation: %w", err)
	}

	return nil
}

// Accept accepts an invitation
func (s *Service) Accept(ctx context.Context, id uuid.UUID) (*Invitation, error) {
	invitation, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, ErrInvitationNotFound
	}

	// Check status
	if invitation.Status == StatusAccepted {
		return nil, ErrInvitationAlreadyAccepted
	}
	if invitation.Status == StatusDeclined {
		return nil, ErrInvitationAlreadyDeclined
	}

	// Check expiration
	expiresAt, err := time.Parse(time.RFC3339, invitation.ExpiresAt)
	if err == nil && time.Now().After(expiresAt) {
		return nil, ErrInvitationExpired
	}

	// Update status
	invitation.Status = StatusAccepted
	invitation.UpdatedAt = time.Now()

	updated, err := s.repo.Update(ctx, id, invitation)
	if err != nil {
		s.logger.Error("failed to accept invitation", "error", err, "id", id)
		return nil, fmt.Errorf("failed to accept invitation: %w", err)
	}

	return updated, nil
}

// Decline declines an invitation
func (s *Service) Decline(ctx context.Context, id uuid.UUID) (*Invitation, error) {
	invitation, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, ErrInvitationNotFound
	}

	// Check status
	if invitation.Status == StatusAccepted {
		return nil, ErrInvitationAlreadyAccepted
	}
	if invitation.Status == StatusDeclined {
		return nil, ErrInvitationAlreadyDeclined
	}

	// Update status
	invitation.Status = StatusDeclined
	invitation.UpdatedAt = time.Now()

	updated, err := s.repo.Update(ctx, id, invitation)
	if err != nil {
		s.logger.Error("failed to decline invitation", "error", err, "id", id)
		return nil, fmt.Errorf("failed to decline invitation: %w", err)
	}

	return updated, nil
}

// Resend resends an invitation
func (s *Service) Resend(ctx context.Context, id uuid.UUID) (*Invitation, error) {
	invitation, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, ErrInvitationNotFound
	}

	// Only pending invitations can be resent
	if invitation.Status != StatusPending {
		return nil, fmt.Errorf("cannot resend invitation with status %s", invitation.Status)
	}

	// Extend expiration
	invitation.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	invitation.UpdatedAt = time.Now()

	updated, err := s.repo.Update(ctx, id, invitation)
	if err != nil {
		s.logger.Error("failed to resend invitation", "error", err, "id", id)
		return nil, fmt.Errorf("failed to resend invitation: %w", err)
	}

	return updated, nil
}
