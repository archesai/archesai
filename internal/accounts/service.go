package accounts

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
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

// ResendVerificationEmail resends the email verification email
func (s *Service) ResendVerificationEmail(_ context.Context, email string) error {
	// TODO: Implement email verification resend
	// This would:
	// 1. Check if user exists by email
	// 2. Generate verification token
	// 3. Store token in database with expiry
	// 4. Send verification email
	s.logger.Info("verification email resend requested", "email", email)
	return fmt.Errorf("email verification resend not yet implemented")
}

// VerifyEmail verifies a user's email address using a verification token
func (s *Service) VerifyEmail(_ context.Context, token string) error {
	// TODO: Implement email verification
	// This would:
	// 1. Validate token and check expiry
	// 2. Find user associated with token
	// 3. Mark user email as verified
	// 4. Invalidate verification token
	s.logger.Info("email verification requested", "token", token)
	return fmt.Errorf("email verification not yet implemented")
}

// RequestPasswordReset initiates a password reset process
func (s *Service) RequestPasswordReset(_ context.Context, email string) error {
	// TODO: Implement password reset request
	// This would:
	// 1. Check if user exists by email
	// 2. Generate reset token
	// 3. Store token in database with expiry
	// 4. Send password reset email
	s.logger.Info("password reset requested", "email", email)
	return fmt.Errorf("password reset not yet implemented")
}

// ConfirmPasswordReset completes the password reset process
func (s *Service) ConfirmPasswordReset(_ context.Context, token, _ string) error {
	// TODO: Implement password reset confirmation
	// This would:
	// 1. Validate token and check expiry
	// 2. Validate new password strength
	// 3. Hash new password
	// 4. Update user password
	// 5. Invalidate reset token
	s.logger.Info("password reset confirmation requested", "token", token)
	return fmt.Errorf("password reset confirmation not yet implemented")
}

// RequestEmailChange initiates an email change process
func (s *Service) RequestEmailChange(_ context.Context, userID uuid.UUID, newEmail string) error {
	// TODO: Implement email change request
	// This would:
	// 1. Validate new email format
	// 2. Check if new email is already taken
	// 3. Generate change token
	// 4. Store token in database with expiry
	// 5. Send confirmation email to new address
	s.logger.Info("email change requested", "user_id", userID, "new_email", newEmail)
	return fmt.Errorf("email change not yet implemented")
}

// ConfirmEmailChange completes the email change process
func (s *Service) ConfirmEmailChange(_ context.Context, token string) error {
	// TODO: Implement email change confirmation
	// This would:
	// 1. Validate token and check expiry
	// 2. Update user email address
	// 3. Mark new email as verified
	// 4. Invalidate change token
	s.logger.Info("email change confirmation requested", "token", token)
	return fmt.Errorf("email change confirmation not yet implemented")
}
