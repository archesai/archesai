package users

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GetByEmail gets a user by email
func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email         string
	Name          string
	EmailVerified bool
}

// Create creates a new user
func (s *Service) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	user := &User{
		ID:            uuid.New(),
		Email:         req.Email,
		Name:          req.Name,
		EmailVerified: req.EmailVerified,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return created, nil
}
