package stargate

import (
	"context"
	"fmt"

	"pharmacy-modernization-project-model/internal/platform/httpclient"

	"go.uber.org/zap"
)

// HTTPClient implements TokenClient using HTTP requests
type HTTPClient struct {
	client    *httpclient.Client
	config    Config
	endpoints EndpointsConfig
	logger    *zap.Logger
}

// NewHTTPClient creates a new HTTP-based Stargate token client
func NewHTTPClient(cfg Config, client *httpclient.Client, logger *zap.Logger) *HTTPClient {
	return &HTTPClient{
		client:    client,
		config:    cfg,
		endpoints: &cfg,
		logger:    logger,
	}
}

// GetAccessToken retrieves an access token using client credentials
func (c *HTTPClient) GetAccessToken(ctx context.Context) (*TokenResponse, error) {
	c.logger.Debug("requesting access token from Stargate",
		zap.String("client_id", c.config.ClientID),
		zap.String("scope", c.config.Scope),
	)

	// Build token request
	tokenReq := TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scope:        c.config.Scope,
	}

	var response TokenResponse
	err := c.client.PostJSON(ctx, c.endpoints.TokenEndpoint(), tokenReq, &response)
	if err != nil {
		c.logger.Error("failed to get access token",
			zap.String("token_url", c.endpoints.TokenEndpoint()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	c.logger.Info("access token obtained from Stargate",
		zap.String("token_type", response.TokenType),
		zap.Int("expires_in", response.ExpiresIn),
		zap.Time("expires_at", response.ExpiresAt()),
	)

	return &response, nil
}

// RefreshToken refreshes an existing token
func (c *HTTPClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	c.logger.Debug("refreshing access token from Stargate")

	// Build refresh request
	refreshReq := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     c.config.ClientID,
		"client_secret": c.config.ClientSecret,
	}

	var response TokenResponse
	err := c.client.PostJSON(ctx, c.endpoints.RefreshTokenEndpoint(), refreshReq, &response)
	if err != nil {
		c.logger.Error("failed to refresh token",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	c.logger.Info("access token refreshed from Stargate",
		zap.Int("expires_in", response.ExpiresIn),
	)

	return &response, nil
}

// Verify HTTPClient implements TokenClient
var _ TokenClient = (*HTTPClient)(nil)
