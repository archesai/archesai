package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Session provides session storage operations
type Session struct {
	client *redis.Client
}

// NewSession creates a new session store
func NewSession(client *redis.Client) *Session {
	return &Session{
		client: client,
	}
}

// SessionData represents session information
type SessionData struct {
	UserID    string                 `json:"user_id"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt time.Time              `json:"expires_at"`
	Data      map[string]interface{} `json:"data"`
}

// Store saves a session with any data type
func (s *Session) Store(key string, data interface{}, ttl time.Duration) error {
	ctx := context.Background()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	return s.client.Set(ctx, key, jsonData, ttl).Err()
}

// Get retrieves a session and unmarshal to destination
func (s *Session) Get(key string, dest interface{}) error {
	ctx := context.Background()
	jsonData, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(jsonData), dest); err != nil {
		return fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	return nil
}

// GetSessionData retrieves session data using SessionData struct
func (s *Session) GetSessionData(sessionID string) (*SessionData, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	var data SessionData
	err := s.Get(key, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// StoreSessionData stores session data using SessionData struct
func (s *Session) StoreSessionData(sessionID string, data *SessionData, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.Store(key, data, ttl)
}

// Delete removes a session
func (s *Session) Delete(key string) error {
	ctx := context.Background()
	return s.client.Del(ctx, key).Err()
}

// DeleteSession removes a session by ID
func (s *Session) DeleteSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.Delete(key)
}

// Refresh updates session TTL
func (s *Session) Refresh(key string, ttl time.Duration) error {
	ctx := context.Background()
	return s.client.Expire(ctx, key, ttl).Err()
}

// RefreshSession updates session TTL by ID
func (s *Session) RefreshSession(sessionID string, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.Refresh(key, ttl)
}

// Exists checks if a session exists
func (s *Session) Exists(key string) (bool, error) {
	ctx := context.Background()
	n, err := s.client.Exists(ctx, key).Result()
	return n > 0, err
}

// SessionExists checks if a session exists by ID
func (s *Session) SessionExists(sessionID string) (bool, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.Exists(key)
}

// GetUserSessions returns all sessions for a user
func (s *Session) GetUserSessions(userID string) ([]string, error) {
	ctx := context.Background()
	pattern := "session:*"
	var userSessions []string

	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		// Get session data to check user ID
		jsonData, err := s.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var data SessionData
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			continue
		}

		if data.UserID == userID {
			sessionID := key[8:] // Remove "session:" prefix
			userSessions = append(userSessions, sessionID)
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return userSessions, nil
}

// DeleteUserSessions removes all sessions for a user
func (s *Session) DeleteUserSessions(userID string) error {
	sessions, err := s.GetUserSessions(userID)
	if err != nil {
		return err
	}

	for _, sessionID := range sessions {
		if err := s.DeleteSession(sessionID); err != nil {
			return err
		}
	}

	return nil
}
