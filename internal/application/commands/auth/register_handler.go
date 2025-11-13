package auth

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/core/models"
	"github.com/archesai/archesai/internal/core/repositories"
	"github.com/archesai/archesai/internal/core/services"
)

// RegisterCommandHandler handles registration commands.
type RegisterCommandHandler struct {
	authService services.AuthService
	userRepo    repositories.UserRepository
	accountRepo repositories.AccountRepository
}

// NewRegisterCommandHandler creates a new register command handler.
func NewRegisterCommandHandler(
	authService services.AuthService,
	userRepo repositories.UserRepository,
	accountRepo repositories.AccountRepository,
) *RegisterCommandHandler {
	return &RegisterCommandHandler{
		authService: authService,
		userRepo:    userRepo,
		accountRepo: accountRepo,
	}
}

// Handle executes the register command.
func (h *RegisterCommandHandler) Handle(
	ctx context.Context,
	cmd *RegisterCommand,
) (*models.AuthTokens, error) {
	if cmd.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if cmd.Password == "" {
		return nil, fmt.Errorf("password is required")
	}
	if cmd.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Register the user
	tokens, err := h.authService.Register(ctx, cmd.Email, cmd.Password, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	return tokens, nil
}
