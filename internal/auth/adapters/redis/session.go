// Package redis provides Redis-based adapters for auth domain
package redis

import (
	"fmt"
	"time"

	"github.com/archesai/archesai/internal/redis"
)

// SessionStore handles session storage in Redis
type SessionStore struct {
	client *redis.Client
}

// NewSessionStore creates a new Redis session store
func NewSessionStore(client *redis.Client) *SessionStore {
	return &SessionStore{
		client: client,
	}
}

// Store saves a session in Redis with the specified TTL
func (s *SessionStore) Store(sessionID string, userID string, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.client.Session.Store(key, userID, ttl)
}

// Get retrieves a session from Redis
func (s *SessionStore) Get(sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	var userID string
	err := s.client.Session.Get(key, &userID)
	return userID, err
}

// Delete removes a session from Redis
func (s *SessionStore) Delete(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.client.Session.Delete(key)
}

// Refresh extends the TTL of a session
func (s *SessionStore) Refresh(sessionID string, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.client.Session.Refresh(key, ttl)
}
