package sessions

// TokenPair represents access and refresh tokens for OAuth
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
