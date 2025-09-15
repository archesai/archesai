package accounts

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrAccountNotFound is returned when an account is not found
	ErrAccountNotFound = errors.New("account not found")
	// ErrInvalidProvider is returned when an invalid provider is specified
	ErrInvalidProvider = errors.New("invalid provider")
	// ErrDuplicateAccount is returned when an account already exists
	ErrDuplicateAccount = errors.New("account already exists")
)

// Service handles account business logic
type Service struct {
	repo     Repository
	cache    Cache
	events   EventPublisher
	logger   *slog.Logger
	cacheTTL time.Duration
}

// NewService creates a new account service
func NewService(repo Repository, cache Cache, events EventPublisher, logger *slog.Logger) *Service {
	return &Service{
		repo:     repo,
		cache:    cache,
		events:   events,
		logger:   logger,
		cacheTTL: 5 * time.Minute,
	}
}

// Create creates a new account
func (s *Service) Create(ctx context.Context, account *Account) (*Account, error) {
	if account.Id == uuid.Nil {
		account.Id = uuid.New()
	}

	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	created, err := s.repo.Create(ctx, account)
	if err != nil {
		s.logger.Error("failed to create account", "error", err)
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	if s.cache != nil {
		if err := s.cache.Set(ctx, created, s.cacheTTL); err != nil {
			s.logger.Warn("failed to cache account", "error", err)
		}
	}

	if s.events != nil {
		if err := s.events.PublishAccountCreated(ctx, created); err != nil {
			s.logger.Warn("failed to publish account created event", "error", err)
		}
	}

	return created, nil
}

// Get retrieves an account by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Account, error) {
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, id); err == nil && cached != nil {
			return cached, nil
		}
	}

	account, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	if s.cache != nil && account != nil {
		if err := s.cache.Set(ctx, account, s.cacheTTL); err != nil {
			s.logger.Warn("failed to cache account", "error", err)
		}
	}

	return account, nil
}

// GetByProviderID retrieves an account by provider and provider account ID
func (s *Service) GetByProviderID(ctx context.Context, provider string, providerAccountID string) (*Account, error) {
	if s.cache != nil {
		if cached, err := s.cache.GetByProviderId(ctx, provider, providerAccountID); err == nil && cached != nil {
			return cached, nil
		}
	}

	account, err := s.repo.GetByProviderId(ctx, provider, providerAccountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	if s.cache != nil && account != nil {
		if err := s.cache.Set(ctx, account, s.cacheTTL); err != nil {
			s.logger.Warn("failed to cache account", "error", err)
		}
	}

	return account, nil
}

// List retrieves a paginated list of accounts
func (s *Service) List(ctx context.Context, params ListAccountsParams) ([]*Account, int64, error) {
	accounts, total, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.Error("failed to list accounts", "error", err)
		return nil, 0, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, total, nil
}

// ListByUserID retrieves all accounts for a specific user
func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Account, error) {
	if s.cache != nil {
		if cached, err := s.cache.ListByUserId(ctx, userID); err == nil && cached != nil {
			return cached, nil
		}
	}

	accounts, err := s.repo.ListByUserId(ctx, userID)
	if err != nil {
		s.logger.Error("failed to list accounts by user", "error", err, "userId", userID)
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, nil
}

// Update updates an existing account
func (s *Service) Update(ctx context.Context, id uuid.UUID, account *Account) (*Account, error) {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	account.Id = existing.Id
	account.CreatedAt = existing.CreatedAt
	account.UpdatedAt = time.Now()

	updated, err := s.repo.Update(ctx, id, account)
	if err != nil {
		s.logger.Error("failed to update account", "error", err)
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	if s.cache != nil {
		if err := s.cache.Delete(ctx, id); err != nil {
			s.logger.Warn("failed to invalidate cache", "error", err)
		}
		if err := s.cache.Set(ctx, updated, s.cacheTTL); err != nil {
			s.logger.Warn("failed to cache updated account", "error", err)
		}
	}

	if s.events != nil {
		if err := s.events.PublishAccountUpdated(ctx, updated); err != nil {
			s.logger.Warn("failed to publish account updated event", "error", err)
		}
	}

	return updated, nil
}

// Delete removes an account by ID
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	account, err := s.repo.Get(ctx, id)
	if err != nil {
		return ErrAccountNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete account", "error", err)
		return fmt.Errorf("failed to delete account: %w", err)
	}

	if s.cache != nil {
		if err := s.cache.Delete(ctx, id); err != nil {
			s.logger.Warn("failed to invalidate cache", "error", err)
		}
	}

	if s.events != nil {
		if err := s.events.PublishAccountDeleted(ctx, account); err != nil {
			s.logger.Warn("failed to publish account deleted event", "error", err)
		}
	}

	return nil
}

// LinkAccount links an account to a user
func (s *Service) LinkAccount(ctx context.Context, userID uuid.UUID, account *Account) (*Account, error) {
	account.UserId = userID

	existing, err := s.repo.GetByProviderId(ctx, string(account.ProviderId), account.AccountId)
	if err == nil && existing != nil {
		if existing.UserId != userID {
			return nil, ErrDuplicateAccount
		}
		return s.Update(ctx, existing.Id, account)
	}

	created, err := s.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	if s.events != nil {
		if err := s.events.PublishAccountLinked(ctx, created); err != nil {
			s.logger.Warn("failed to publish account linked event", "error", err)
		}
	}

	return created, nil
}

// UnlinkAccount unlinks an account from a user
func (s *Service) UnlinkAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	account, err := s.Get(ctx, accountID)
	if err != nil {
		return err
	}

	if account.UserId != userID {
		return ErrAccountNotFound
	}

	if err := s.Delete(ctx, accountID); err != nil {
		return err
	}

	if s.events != nil {
		if err := s.events.PublishAccountUnlinked(ctx, account); err != nil {
			s.logger.Warn("failed to publish account unlinked event", "error", err)
		}
	}

	return nil
}
