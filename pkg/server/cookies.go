package server

import (
	"net/http"

	"github.com/archesai/archesai/pkg/auth"
)

// CookieManager handles HTTP cookie operations for authentication.
// This belongs in the HTTP adapter layer, not in business logic.
type CookieManager struct {
	secure   bool
	sameSite http.SameSite
	domain   string
}

// NewCookieManager creates a new cookie manager.
func NewCookieManager(secure bool, domain string) *CookieManager {
	return &CookieManager{
		secure:   secure,
		sameSite: http.SameSiteLaxMode,
		domain:   domain,
	}
}

// SetAuthCookies sets authentication cookies from tokens.
func (c *CookieManager) SetAuthCookies(
	w http.ResponseWriter,
	tokens *auth.Tokens,
	rememberMe bool,
) {
	// Set access token cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		Domain:   c.domain,
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: c.sameSite,
		MaxAge:   tokens.ExpiresIn, // Typically 15 minutes
	}

	// Set refresh token cookie with longer expiry if rememberMe is true
	refreshMaxAge := 7 * 24 * 60 * 60 // 7 days default
	if rememberMe {
		refreshMaxAge = 30 * 24 * 60 * 60 // 30 days if remember me
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		Domain:   c.domain,
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: c.sameSite,
		MaxAge:   refreshMaxAge,
	}

	// Set session ID cookie
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    tokens.SessionID,
		Path:     "/",
		Domain:   c.domain,
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: c.sameSite,
		MaxAge:   refreshMaxAge,
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, sessionCookie)
}

// ClearAuthCookies clears all authentication cookies.
func (c *CookieManager) ClearAuthCookies(w http.ResponseWriter) {
	cookies := []string{"access_token", "refresh_token", "session_id"}

	for _, name := range cookies {
		cookie := &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Domain:   c.domain,
			HttpOnly: true,
			Secure:   c.secure,
			SameSite: c.sameSite,
			MaxAge:   -1, // Delete the cookie
		}
		http.SetCookie(w, cookie)
	}
}

// GetTokenFromCookie retrieves a token from cookies.
func (c *CookieManager) GetTokenFromCookie(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
