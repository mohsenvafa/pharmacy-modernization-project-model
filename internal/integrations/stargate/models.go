package stargate

import "time"

// TokenRequest represents a request for an access token
type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope,omitempty"`
}

// TokenResponse represents the response from Stargate token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// ExpiresAt calculates when the token expires
func (r *TokenResponse) ExpiresAt() time.Time {
	return time.Now().Add(time.Duration(r.ExpiresIn) * time.Second)
}
