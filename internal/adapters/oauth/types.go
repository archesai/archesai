package oauth

// OAuthTokens represents OAuth tokens received from providers
type OAuthTokens struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	ExpiresIn    int
	TokenType    string
	Scope        string
}

// OAuthUserInfo represents user information from OAuth providers
type OAuthUserInfo struct {
	ID            string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
	Provider      string
}
