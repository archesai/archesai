package auth

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// OAuth state storage using a simple in-memory map
// In production, this should use Redis or another distributed cache
var oauthStateStore = struct {
	sync.RWMutex
	states map[string]oauthState
}{
	states: make(map[string]oauthState),
}

type oauthState struct {
	Provider    string
	RedirectURI string
	ExpiresAt   time.Time
}

// StoreOAuthState stores OAuth state data temporarily
func (sm *SessionManager) StoreOAuthState(_ context.Context, state, provider, redirectURI string, ttl time.Duration) error {
	oauthStateStore.Lock()
	defer oauthStateStore.Unlock()

	// Clean up expired states
	now := time.Now()
	for k, v := range oauthStateStore.states {
		if v.ExpiresAt.Before(now) {
			delete(oauthStateStore.states, k)
		}
	}

	oauthStateStore.states[state] = oauthState{
		Provider:    provider,
		RedirectURI: redirectURI,
		ExpiresAt:   now.Add(ttl),
	}

	return nil
}

// GetOAuthRedirectURI retrieves the redirect URI for a state
func (sm *SessionManager) GetOAuthRedirectURI(_ context.Context, state string) (string, error) {
	oauthStateStore.RLock()
	defer oauthStateStore.RUnlock()

	if s, ok := oauthStateStore.states[state]; ok {
		if s.ExpiresAt.After(time.Now()) {
			return s.RedirectURI, nil
		}
	}

	return "", fmt.Errorf("state not found or expired")
}

// DeleteOAuthState removes OAuth state data
func (sm *SessionManager) DeleteOAuthState(_ context.Context, state string) error {
	oauthStateStore.Lock()
	defer oauthStateStore.Unlock()

	delete(oauthStateStore.states, state)
	return nil
}
