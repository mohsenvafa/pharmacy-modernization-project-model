package stargate

import (
	"context"

	"go.uber.org/zap"
)

// TokenProviderAdapter adapts the Stargate TokenClient to httpclient.TokenProvider interface
// This allows Stargate to be used as a token provider for other HTTP clients
type TokenProviderAdapter struct {
	client TokenClient
	logger *zap.Logger
}

// NewTokenProviderAdapter creates an adapter that bridges Stargate TokenClient to TokenProvider
func NewTokenProviderAdapter(client TokenClient, logger *zap.Logger) *TokenProviderAdapter {
	return &TokenProviderAdapter{
		client: client,
		logger: logger,
	}
}

// GetToken implements the httpclient.TokenProvider interface
func (a *TokenProviderAdapter) GetToken(ctx context.Context) (string, error) {
	a.logger.Debug("fetching access token via Stargate")

	tokenResp, err := a.client.GetAccessToken(ctx)
	if err != nil {
		a.logger.Error("failed to get token from Stargate", zap.Error(err))
		return "", err
	}

	a.logger.Debug("access token obtained from Stargate",
		zap.String("token_type", tokenResp.TokenType),
		zap.Int("expires_in", tokenResp.ExpiresIn),
	)

	return tokenResp.AccessToken, nil
}
