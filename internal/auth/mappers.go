// Package auth conversion functions for database entities
// This file contains manual implementations of conversion functions
// between database models and API/domain models.
package auth

import (
	"github.com/archesai/archesai/internal/database/postgresql"
)

// UserDBToAPI converts a database User to an API User
func UserDBToAPI(db interface{}) *User {
	v, ok := db.(*postgresql.User)
	if !ok || v == nil {
		return nil
	}

	user := &User{
		Id:            v.Id,
		Email:         Email(v.Email),
		Name:          v.Name,
		EmailVerified: v.EmailVerified,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
	}
	if v.Image != nil {
		user.Image = *v.Image
	}
	return user
}

// SessionDBToAPI converts a database Session to an API Session
func SessionDBToAPI(db interface{}) *Session {
	v, ok := db.(*postgresql.Session)
	if !ok || v == nil {
		return nil
	}

	session := &Session{
		Id:        v.Id,
		Token:     v.Token,
		UserId:    v.UserId,
		ExpiresAt: v.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
	if v.ActiveOrganizationId != nil {
		session.ActiveOrganizationId = *v.ActiveOrganizationId
	}
	if v.IpAddress != nil {
		session.IpAddress = *v.IpAddress
	}
	if v.UserAgent != nil {
		session.UserAgent = *v.UserAgent
	}
	return session
}

// AccountDBToAPI converts a database Account to an API Account
func AccountDBToAPI(db interface{}) *Account {
	v, ok := db.(*postgresql.Account)
	if !ok || v == nil {
		return nil
	}

	account := &Account{
		Id:         v.Id,
		UserId:     v.UserId,
		ProviderId: AccountProviderId(v.ProviderId),
		AccountId:  v.AccountId,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
	}
	if v.AccessToken != nil {
		account.AccessToken = *v.AccessToken
	}
	if v.RefreshToken != nil {
		account.RefreshToken = *v.RefreshToken
	}
	if v.IdToken != nil {
		account.IdToken = *v.IdToken
	}
	if v.Password != nil {
		account.Password = *v.Password
	}
	if v.AccessTokenExpiresAt != nil {
		account.AccessTokenExpiresAt = *v.AccessTokenExpiresAt
	}
	return account
}
