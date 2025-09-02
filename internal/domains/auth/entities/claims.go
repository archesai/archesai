// Package entities defines domain entities and related types for the auth module.
package entities

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents JWT claims
type Claims struct {
	UserID         uuid.UUID  `json:"user_id"`
	Email          string     `json:"email"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Role           string     `json:"role,omitempty"`
	jwt.RegisteredClaims
}
