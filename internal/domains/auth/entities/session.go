package entities

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user session
type Session struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	UserID               uuid.UUID  `json:"user_id" db:"user_id"`
	Token                string     `json:"token" db:"token"`
	ActiveOrganizationID *uuid.UUID `json:"active_organization_id,omitempty" db:"active_organization_id"`
	IPAddress            *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent            *string    `json:"user_agent,omitempty" db:"user_agent"`
	ExpiresAt            time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}
