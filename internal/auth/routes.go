package auth

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all auth routes
func RegisterRoutes(e *echo.Echo, handler *Handler) {
	auth := e.Group("/auth")

	// Password authentication
	auth.POST("/login", handler.Login)
	auth.POST("/logout", handler.Logout)
	auth.POST("/logout-all", handler.LogoutAll)

	// Magic link authentication
	auth.POST("/magic-link", handler.RequestMagicLink)
	auth.GET("/magic-link/verify", handler.VerifyMagicLink)

	// OAuth authentication
	auth.GET("/oauth/:provider", handler.OAuthLogin)
	auth.GET("/oauth/:provider/callback", handler.OAuthCallback)

	// Token management
	auth.POST("/refresh", handler.RefreshToken)
	auth.POST("/api-tokens", handler.CreateAPIToken)
}
