package users

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service provides user management business logic
type Service struct {
	repo   Repository
	cache  Cache
	events EventPublisher
	logger *slog.Logger
}

// NewService creates a new user service
func NewService(repo Repository, cache Cache, events EventPublisher, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		cache:  cache,
		events: events,
		logger: logger,
	}
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	// Try cache first
	user, err := s.cache.GetUser(ctx, id)
	if err == nil && user != nil {
		return user, nil
	}

	// Cache miss - get from database
	entity, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update cache for next time
	_ = s.cache.SetUser(ctx, entity, 5*time.Minute)

	return entity, nil
}

// UpdateUser updates user information
func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, req *UpdateUserJSONBody) (*User, error) {
	// Get existing user
	entity, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Email != "" {
		entity.Email = Email(req.Email)
	}
	if req.Image != "" {
		entity.Image = req.Image
	}
	entity.UpdatedAt = time.Now()

	// Save changes
	updatedEntity, err := s.repo.UpdateUser(ctx, id, entity)
	if err != nil {
		return nil, err
	}

	// Update cache
	_ = s.cache.SetUser(ctx, updatedEntity, 5*time.Minute)

	// Publish event
	_ = s.events.PublishUserUpdated(ctx, updatedEntity)

	return updatedEntity, nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Get user first for event publishing
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return err
	}

	// Delete from repository
	err = s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	// Remove from cache
	_ = s.cache.DeleteUser(ctx, id)
	_ = s.cache.DeleteUserByEmail(ctx, string(user.Email))

	// Publish event
	_ = s.events.PublishUserDeleted(ctx, user)

	return nil
}

// ListUsers lists users with pagination
func (s *Service) ListUsers(ctx context.Context, limit, offset int32) ([]*User, error) {
	params := ListUsersParams{
		Limit:  int(limit),
		Offset: int(offset),
	}
	entities, _, err := s.repo.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	users := make([]*User, len(entities))
	copy(users, entities)
	return users, nil
}

// GetUserByEmail retrieves a user by email address
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// Try cache first
	user, err := s.cache.GetUserByEmail(ctx, email)
	if err == nil && user != nil {
		return user, nil
	}

	// Cache miss - get from database
	entity, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Update cache for next time
	_ = s.cache.SetUserByEmail(ctx, email, entity, 5*time.Minute)

	return entity, nil
}
