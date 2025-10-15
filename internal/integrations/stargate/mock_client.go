package stargate

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// MockClient implements TokenClient with in-memory mock data
type MockClient struct {
	mockToken        string
	mockExpiresIn    int
	mockRefreshToken string
	logger           *zap.Logger
}

// NewMockClient creates a new mock Stargate token client
func NewMockClient(logger *zap.Logger) *MockClient {
	return &MockClient{
		mockToken:        "mock-access-token-12345",
		mockExpiresIn:    3600, // 1 hour
		mockRefreshToken: "mock-refresh-token-67890",
		logger:           logger,
	}
}

// SetMockToken sets a custom mock token for testing
func (c *MockClient) SetMockToken(token string, expiresIn int) {
	c.mockToken = token
	c.mockExpiresIn = expiresIn
}

// GetAccessToken returns a mock access token
func (c *MockClient) GetAccessToken(ctx context.Context) (*TokenResponse, error) {
	c.logger.Debug("mock: returning access token")

	response := &TokenResponse{
		AccessToken:  c.mockToken,
		TokenType:    "Bearer",
		ExpiresIn:    c.mockExpiresIn,
		RefreshToken: c.mockRefreshToken,
		Scope:        "read write",
	}

	c.logger.Info("mock: access token obtained",
		zap.String("token_type", response.TokenType),
		zap.Int("expires_in", response.ExpiresIn),
		zap.Time("expires_at", response.ExpiresAt()),
	)

	return response, nil
}

// RefreshToken returns a new mock token
func (c *MockClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	c.logger.Debug("mock: refreshing access token",
		zap.String("refresh_token", refreshToken),
	)

	// Return a new mock token
	response := &TokenResponse{
		AccessToken:  "mock-refreshed-token-" + time.Now().Format("20060102150405"),
		TokenType:    "Bearer",
		ExpiresIn:    c.mockExpiresIn,
		RefreshToken: c.mockRefreshToken,
		Scope:        "read write",
	}

	c.logger.Info("mock: access token refreshed",
		zap.Int("expires_in", response.ExpiresIn),
	)

	return response, nil
}

// Verify MockClient implements TokenClient
var _ TokenClient = (*MockClient)(nil)
