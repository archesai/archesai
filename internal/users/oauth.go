package users

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// OAuthUserInfo represents user information from an OAuth provider
type OAuthUserInfo struct {
	ProviderAccountID string
	Email             string
	EmailVerified     bool
	Name              string
	Picture           string
	Locale            string
}

// FindOrCreateFromOAuth finds or creates a user from OAuth provider information
func (s *Service) FindOrCreateFromOAuth(
	ctx context.Context,
	provider string,
	userInfo *OAuthUserInfo,
) (*User, error) {
	// Try to find existing user by email
	existingUser, err := s.repo.GetByEmail(ctx, userInfo.Email)
	if err == nil {
		// User exists, update provider info if needed
		s.logger.Info("Found existing user for OAuth login",
			"email", userInfo.Email,
			"provider", provider)

		// Update last login time
		existingUser.UpdatedAt = time.Now()

		// If email wasn't verified before but OAuth provider says it is
		if !existingUser.EmailVerified && userInfo.EmailVerified {
			existingUser.EmailVerified = true
		}

		// Update image if provided and different
		if userInfo.Picture != "" {
			existingUser.Image = &userInfo.Picture
		}

		// Update the user
		_, err = s.repo.Update(ctx, existingUser.ID, existingUser)
		if err != nil {
			s.logger.Error("Failed to update user after OAuth login", "error", err)
			// Continue anyway, user can still login
		}

		return existingUser, nil
	}

	// User doesn't exist, create new one
	s.logger.Info("Creating new user from OAuth",
		"email", userInfo.Email,
		"provider", provider)

	newUser := &User{
		ID:            uuid.New(),
		Email:         userInfo.Email,         // Convert string to Email type
		EmailVerified: userInfo.EmailVerified, // OAuth providers usually verify emails
		Name:          userInfo.Name,
		Image:         &userInfo.Picture, // Map Picture to Image field
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create the user
	_, err = s.repo.Create(ctx, newUser)
	if err != nil {
		s.logger.Error("Failed to create user from OAuth",
			"error", err,
			"email", userInfo.Email)
		return nil, err
	}

	return newUser, nil
}
