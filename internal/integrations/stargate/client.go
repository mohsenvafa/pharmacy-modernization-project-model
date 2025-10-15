package stargate

import "context"

// TokenClient defines the interface for interacting with Stargate authentication API
type TokenClient interface {
	// GetAccessToken retrieves an access token using client credentials
	GetAccessToken(ctx context.Context) (*TokenResponse, error)

	// RefreshToken refreshes an existing token (if refresh tokens are supported)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
}
