package sessions

import (
	"log/slog"
)

// Service provides session business logic
type Service struct {
	repo           Repository
	cache          Cache
	sessionManager *SessionManager
	logger         *slog.Logger
}

// NewService creates a new session service
func NewService(repo Repository, cache Cache, logger *slog.Logger) *Service {
	sessionManager := NewSessionManager(repo, cache, 0) // Use default TTL
	return &Service{
		repo:           repo,
		cache:          cache,
		sessionManager: sessionManager,
		logger:         logger,
	}
}
