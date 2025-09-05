package redis

import (
	"encoding/json"
	"fmt"
	"time"
)

// Session provides session storage operations
type Session struct{}

// NewSession creates a new session store
func NewSession() *Session {
	return &Session{}
}

// SessionData represents session information
type SessionData struct {
	UserID    string                 `json:"user_id"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt time.Time              `json:"expires_at"`
	Data      map[string]interface{} `json:"data"`
}

// Store saves a session
func (s *Session) Store(sessionID string, data *SessionData, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	key := fmt.Sprintf("session:%s", sessionID)
	return client.Set(ctx, key, jsonData, ttl).Err()
}

// Get retrieves a session
func (s *Session) Get(sessionID string) (*SessionData, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	jsonData, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var data SessionData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	return &data, nil
}

// Delete removes a session
func (s *Session) Delete(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return client.Del(ctx, key).Err()
}

// Refresh updates session TTL
func (s *Session) Refresh(sessionID string, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return client.Expire(ctx, key, ttl).Err()
}

// Exists checks if a session exists
func (s *Session) Exists(sessionID string) (bool, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	n, err := client.Exists(ctx, key).Result()
	return n > 0, err
}

// GetUserSessions returns all sessions for a user
func (s *Session) GetUserSessions(userID string) ([]string, error) {
	pattern := "session:*"
	var userSessions []string

	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		// Get session data to check user ID
		jsonData, err := client.Get(ctx, key).Result()
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
		if err := s.Delete(sessionID); err != nil {
			return err
		}
	}

	return nil
}
