package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/email"
	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
)

// SetDatabaseQueries sets the database queries for the service
func (s *Service) SetDatabaseQueries(queries *postgresql.Queries) {
	s.dbQueries = queries
}

// SetEmailService sets the email service for the service
func (s *Service) SetEmailService(emailService *email.Service) {
	s.emailService = emailService
}

// SetAPIKeyService sets the API key service for the service
func (s *Service) SetAPIKeyService(apiKeyService *APIKeyService) {
	s.apiKeyService = apiKeyService
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*users.User, *TokenResponse, error) {
	// Check if user already exists
	existingUser, err := s.usersRepo.GetUserByEmail(ctx, string(req.Email))
	if err == nil && existingUser != nil {
		return nil, nil, ErrUserExists
	}

	// Validate password strength
	if err := s.validatePassword(req.Password); err != nil {
		return nil, nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Hash the password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user with embedded User
	now := time.Now()
	user := &users.User{
		Id:            uuid.New(),
		Email:         req.Email,
		Name:          req.Name,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Save user to database - repository expects User
	createdEntity, err := s.usersRepo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", "error", err)
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}
	// Update user with created entity (in case DB added fields)
	user = createdEntity

	// Create local account with password
	account := &Account{
		Id:         uuid.New(),
		UserId:     user.Id,
		ProviderId: Local,
		AccountId:  string(user.Email), // Use email as account ID for local auth
		Password:   hashedPassword,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	_, err = s.repo.CreateAccount(ctx, account)
	if err != nil {
		s.logger.Error("failed to create account", "error", err)
		// Try to clean up the created user
		_ = s.usersRepo.DeleteUser(ctx, user.Id)
		return nil, nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Generate email verification token if email service is configured
	if s.emailService != nil && s.dbQueries != nil {
		verificationToken, err := s.generateVerificationToken()
		if err != nil {
			s.logger.Error("failed to generate verification token", "error", err)
			// Continue without email verification
		} else {
			// Store verification token in database
			_, err = s.dbQueries.CreateVerificationToken(ctx, postgresql.CreateVerificationTokenParams{
				Id:         uuid.New(),
				Identifier: string(user.Email),
				Value:      verificationToken,
				ExpiresAt:  time.Now().Add(24 * time.Hour), // Token expires in 24 hours
			})
			if err != nil {
				s.logger.Error("failed to store verification token", "error", err)
				// Continue without email verification
			} else {
				// Send verification email
				err = s.emailService.SendVerificationEmail(ctx, string(user.Email), user.Name, verificationToken)
				if err != nil {
					s.logger.Error("failed to send verification email", "error", err)
					// Continue - user can request resend later
				}
			}
		}
	}

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		s.logger.Error("failed to generate tokens", "error", err)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session - use SessionManager if available
	var session *Session
	if s.sessionManager != nil {
		session, err = s.sessionManager.CreateSession(ctx, user.Id, uuid.Nil, "", "")
		if err != nil {
			s.logger.Error("failed to create session", "error", err)
			return nil, nil, fmt.Errorf("failed to create session: %w", err)
		}
		// Update session with refresh token
		session.Token = tokens.RefreshToken
		_, err = s.sessionManager.UpdateSession(ctx, session.Id, session)
		if err != nil {
			s.logger.Error("failed to update session token", "error", err)
		}
	} else {
		// Fallback to direct repository
		sessionNow := time.Now()
		session = &Session{
			Id:        uuid.New(),
			UserId:    user.Id,
			Token:     tokens.RefreshToken,
			ExpiresAt: sessionNow.Add(s.config.SessionTokenExpiry).Format(time.RFC3339),
			CreatedAt: sessionNow,
			UpdatedAt: sessionNow,
			// Required fields with empty defaults
			ActiveOrganizationId: uuid.Nil,
			IpAddress:            "",
			UserAgent:            "",
		}
		_, err = s.repo.CreateSession(ctx, session)
		if err != nil {
			s.logger.Error("failed to create session", "error", err)
			return nil, nil, fmt.Errorf("failed to create session: %w", err)
		}
	}

	s.logger.Info("user signed up successfully", "user_id", user.Id.String())
	return user, tokens, nil
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*users.User, *TokenResponse, error) {
	// Check if IP is locked out due to brute force attempts
	if s.config.MaxLoginAttempts > 0 {
		if s.isIPLockedOut(ctx, ipAddress) {
			s.logger.Warn("IP address locked out due to brute force attempts", "ip", ipAddress)
			return nil, nil, fmt.Errorf("too many failed login attempts, try again later")
		}
	}

	// Get user by email
	userEntity, err := s.usersRepo.GetUserByEmail(ctx, string(req.Email))
	if err != nil {
		// Track failed attempt
		s.trackFailedLoginAttempt(ctx, ipAddress, string(req.Email))
		s.logger.Warn("user not found", "email", req.Email)
		return nil, nil, ErrInvalidCredentials
	}

	// Get the user's local account to verify password
	account, err := s.repo.GetAccountByProviderAndProviderID(ctx, string(Local), string(req.Email))
	if err != nil {
		// Track failed attempt
		s.trackFailedLoginAttempt(ctx, ipAddress, string(req.Email))
		s.logger.Warn("account not found", "email", req.Email)
		return nil, nil, ErrInvalidCredentials
	}

	// Verify password
	if account.Password != "" {
		if err := s.verifyPassword(req.Password, account.Password); err != nil {
			// Track failed attempt
			s.trackFailedLoginAttempt(ctx, ipAddress, string(req.Email))
			s.logger.Warn("invalid password", "user_id", userEntity.Id.String(), "ip", ipAddress)
			return nil, nil, ErrInvalidCredentials
		}
	}

	user := userEntity

	// Clear any failed login attempts on successful authentication
	s.clearFailedAttempts(ctx, ipAddress, string(req.Email))

	// Check concurrent session limits
	if s.config.MaxConcurrentSessions > 0 {
		activeSessions, err := s.ListUserSessions(ctx, user.Id)
		if err == nil && len(activeSessions) >= s.config.MaxConcurrentSessions {
			// Remove oldest session if limit reached
			s.logger.Info("concurrent session limit reached, removing oldest session",
				"user_id", user.Id,
				"max_sessions", s.config.MaxConcurrentSessions,
				"active_sessions", len(activeSessions))

			// Find and remove the oldest session
			if len(activeSessions) > 0 {
				oldestSession := activeSessions[0]
				for _, session := range activeSessions {
					if session.CreatedAt.Before(oldestSession.CreatedAt) {
						oldestSession = session
					}
				}
				_ = s.repo.DeleteSession(ctx, oldestSession.Id)
			}
		}
	}

	// Generate tokens with extended refresh token if remember me is enabled
	var tokens *TokenResponse
	if req.RememberMe {
		// Use extended refresh token expiry for remember me
		extendedConfig := s.config
		extendedConfig.RefreshTokenExpiry = 30 * 24 * time.Hour // 30 days for remember me
		tokens, err = s.generateTokensWithConfig(user, extendedConfig)
		s.logger.Info("generating extended session for remember me", "user_id", user.Id)
	} else {
		tokens, err = s.generateTokens(user)
	}
	if err != nil {
		s.logger.Error("failed to generate tokens", "error", err)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session - use SessionManager if available
	var session *Session
	if s.sessionManager != nil {
		// TODO: Get organization ID from user's default org
		session, err = s.sessionManager.CreateSession(ctx, user.Id, uuid.Nil, ipAddress, userAgent)
		if err != nil {
			s.logger.Error("failed to create session", "error", err)
			return nil, nil, fmt.Errorf("failed to create session: %w", err)
		}
		// Update session with refresh token
		session.Token = tokens.RefreshToken
		_, err = s.sessionManager.UpdateSession(ctx, session.Id, session)
		if err != nil {
			s.logger.Error("failed to update session token", "error", err)
		}
	} else {
		// Fallback to direct repository
		sessionNow := time.Now()
		session = &Session{
			Id:                   uuid.New(),
			UserId:               user.Id,
			Token:                tokens.RefreshToken,
			ExpiresAt:            sessionNow.Add(s.config.SessionTokenExpiry).Format(time.RFC3339),
			CreatedAt:            sessionNow,
			UpdatedAt:            sessionNow,
			ActiveOrganizationId: uuid.Nil, // TODO: Set proper organization ID
			IpAddress:            ipAddress,
			UserAgent:            userAgent,
		}
		_, err = s.repo.CreateSession(ctx, session)
		if err != nil {
			s.logger.Error("failed to create session", "error", err)
			return nil, nil, fmt.Errorf("failed to create session: %w", err)
		}
	}

	s.logger.Info("user signed in successfully", "user_id", userEntity.Id.String())
	return userEntity, tokens, nil
}

// Logout invalidates a user session
func (s *Service) Logout(ctx context.Context, token string) error {
	// Use SessionManager if available
	if s.sessionManager != nil {
		err := s.sessionManager.DeleteSessionByToken(ctx, token)
		if err != nil {
			s.logger.Error("failed to delete session", "error", err)
			return ErrInvalidToken
		}
		s.logger.Info("user signed out successfully")
		return nil
	}

	// Fallback to direct repository
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return ErrInvalidToken
	}

	// Delete session
	if err := s.repo.DeleteSession(ctx, session.Id); err != nil {
		s.logger.Error("failed to delete session", "error", err)
		return fmt.Errorf("failed to delete session: %w", err)
	}

	s.logger.Info("user signed out successfully", "user_id", session.UserId)
	return nil
}

// GetUserByID retrieves a user by their ID
func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*users.User, error) {
	return s.usersRepo.GetUser(ctx, userID)
}
