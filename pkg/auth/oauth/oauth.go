// Package oauth implements OAuth providers for authentication
package oauth

import "context"

// Provider defines the interface for OAuth providers
type Provider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*Tokens, error)
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)
}

// Tokens represents OAuth tokens received from providers
type Tokens struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
	IDToken      string
	Scope        string
}

// UserInfo represents user information from OAuth providers
type UserInfo struct {
	ID            string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
	Provider      string
}
