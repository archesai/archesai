package valueobjects

// AuthTokens represents the tokens returned after successful authentication.
type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	SessionID    string `json:"sessionId"`
}
