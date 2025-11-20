package auth

// Tokens represents the tokens returned after successful authentication.
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	SessionID    string `json:"sessionId"`
}
