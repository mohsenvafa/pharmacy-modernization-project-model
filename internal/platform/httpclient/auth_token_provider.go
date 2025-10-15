package httpclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// TokenProvider is an interface for services that provide access tokens
type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
}

// CachedTokenProvider caches tokens and refreshes them when expired
type CachedTokenProvider struct {
	provider      TokenProvider
	logger        *zap.Logger
	mu            sync.RWMutex
	cachedToken   string
	expiresAt     time.Time
	refreshBefore time.Duration // Refresh token this much before expiry
}

// NewCachedTokenProvider creates a token provider with caching
func NewCachedTokenProvider(provider TokenProvider, refreshBefore time.Duration, logger *zap.Logger) *CachedTokenProvider {
	if refreshBefore == 0 {
		refreshBefore = 5 * time.Minute // Default: refresh 5 min before expiry
	}

	return &CachedTokenProvider{
		provider:      provider,
		logger:        logger,
		refreshBefore: refreshBefore,
	}
}

// GetToken returns a cached token or fetches a new one if expired
func (p *CachedTokenProvider) GetToken(ctx context.Context) (string, error) {
	p.mu.RLock()
	if p.cachedToken != "" && time.Now().Before(p.expiresAt.Add(-p.refreshBefore)) {
		token := p.cachedToken
		p.mu.RUnlock()
		p.logger.Debug("using cached access token",
			zap.Time("expires_at", p.expiresAt),
		)
		return token, nil
	}
	p.mu.RUnlock()

	// Need to refresh token
	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock
	if p.cachedToken != "" && time.Now().Before(p.expiresAt.Add(-p.refreshBefore)) {
		return p.cachedToken, nil
	}

	p.logger.Info("fetching new access token")

	token, err := p.provider.GetToken(ctx)
	if err != nil {
		p.logger.Error("failed to fetch access token", zap.Error(err))
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	p.cachedToken = token
	p.expiresAt = time.Now().Add(55 * time.Minute) // Assume 1-hour tokens

	p.logger.Info("access token refreshed",
		zap.Time("expires_at", p.expiresAt),
	)

	return token, nil
}

// SetTokenExpiry allows setting a custom expiry time (useful when token has explicit expiry)
func (p *CachedTokenProvider) SetTokenExpiry(expiresAt time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.expiresAt = expiresAt
}

// InvalidateToken forces the next request to fetch a fresh token
func (p *CachedTokenProvider) InvalidateToken() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.logger.Info("invalidating cached token")
	p.cachedToken = ""
	p.expiresAt = time.Time{}
}

// AuthHeaderProvider provides Authorization headers using a token provider
type AuthHeaderProvider struct {
	tokenProvider TokenProvider
	authType      string // "Bearer", "Token", etc.
	logger        *zap.Logger
}

// NewAuthHeaderProvider creates a header provider for Authorization headers
func NewAuthHeaderProvider(tokenProvider TokenProvider, authType string, logger *zap.Logger) *AuthHeaderProvider {
	if authType == "" {
		authType = "Bearer"
	}

	return &AuthHeaderProvider{
		tokenProvider: tokenProvider,
		authType:      authType,
		logger:        logger,
	}
}

func (p *AuthHeaderProvider) GetHeaders(ctx context.Context) (map[string]string, error) {
	token, err := p.tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}

	return map[string]string{
		"Authorization": p.authType + " " + token,
	}, nil
}
