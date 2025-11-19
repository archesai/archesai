package oauth

// Tokens represents OAuth tokens received from providers
type Tokens struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	ExpiresIn    int
	TokenType    string
	Scope        string
}

// UserInfo represents user information from OAuth providers
type UserInfo struct {
	ID            string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
	Provider      string
}
