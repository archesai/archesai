package oauth

import (
	"context"
	"fmt"
	"log/slog"
)

// StrictServer implements StrictServerInterface for OAuth
type StrictServer struct {
	service *Service
	logger  *slog.Logger
}

// NewStrictServer creates a new strict server implementation
func NewStrictServer(service *Service, logger *slog.Logger) *StrictServer {
	return &StrictServer{
		service: service,
		logger:  logger,
	}
}

// OauthAuthorize starts the OAuth authorization flow
func (s *StrictServer) OauthAuthorize(
	ctx context.Context,
	request OauthAuthorizeRequestObject,
) (OauthAuthorizeResponseObject, error) {
	providerName := string(request.Provider)

	// Get the redirect URI from params or use default
	redirectURI := request.Params.RedirectURI
	if redirectURI == "" {
		// Default redirect URI based on provider
		redirectURI = fmt.Sprintf("http://localhost:8080/auth/oauth/%s/callback", providerName)
	}

	// Get authorization URL from service
	authURL, err := s.service.GetAuthorizationURL(ctx, providerName, redirectURI)
	if err != nil {
		s.logger.Error("Failed to get authorization URL",
			"provider", providerName,
			"error", err)

		// Return error redirect
		errorURL := s.service.BuildCallbackURL(providerName, nil, err)
		return OauthAuthorize302Response{
			Headers: OauthAuthorize302ResponseHeaders{
				Location: errorURL,
			},
		}, nil
	}

	// Redirect to provider's authorization page
	return OauthAuthorize302Response{
		Headers: OauthAuthorize302ResponseHeaders{
			Location: authURL,
		},
	}, nil
}

// OauthCallback handles the OAuth callback from the provider
func (s *StrictServer) OauthCallback(
	ctx context.Context,
	request OauthCallbackRequestObject,
) (OauthCallbackResponseObject, error) {
	providerName := string(request.Provider)

	// Extract code and state from query params
	code := request.Params.Code
	state := request.Params.State

	// Check for errors from provider
	if request.Params.Error != "" {
		err := fmt.Errorf("provider error: %s", request.Params.Error)
		s.logger.Error("OAuth provider returned error",
			"provider", providerName,
			"error", request.Params.Error)

		errorURL := s.service.BuildCallbackURL(providerName, nil, err)
		return OauthCallback302Response{
			Headers: OauthCallback302ResponseHeaders{
				Location: errorURL,
			},
		}, nil
	}

	// Validate code is present
	if code == "" {
		err := fmt.Errorf("authorization code not provided")
		errorURL := s.service.BuildCallbackURL(providerName, nil, err)
		return OauthCallback302Response{
			Headers: OauthCallback302ResponseHeaders{
				Location: errorURL,
			},
		}, nil
	}

	// Get the redirect URI (should match what was sent in authorize)
	redirectURI := fmt.Sprintf("http://localhost:8080/auth/oauth/%s/callback", providerName)

	// Handle the callback through service
	tokenPair, err := s.service.HandleCallback(ctx, providerName, code, state, redirectURI)
	if err != nil {
		s.logger.Error("Failed to handle OAuth callback",
			"provider", providerName,
			"error", err)

		errorURL := s.service.BuildCallbackURL(providerName, nil, err)
		return OauthCallback302Response{
			Headers: OauthCallback302ResponseHeaders{
				Location: errorURL,
			},
		}, nil
	}

	// Success! Redirect to frontend with tokens
	successURL := s.service.BuildCallbackURL(providerName, tokenPair, nil)
	return OauthCallback302Response{
		Headers: OauthCallback302ResponseHeaders{
			Location: successURL,
		},
	}, nil
}
