package auth

import (
	"fmt"

	"go.uber.org/zap"
)

// Builder helps configure and initialize the authentication system
type Builder struct {
	jwtConfig JWTConfig
	devMode   bool
	env       string
	logger    *zap.Logger
}

// NewBuilder creates a new auth builder
func NewBuilder() *Builder {
	return &Builder{}
}

// WithJWTConfig sets JWT configuration
func (b *Builder) WithJWTConfig(issuer, audience, clientIds []string, cookieName string) *Builder {
	b.jwtConfig = JWTConfig{
		Issuer:     issuer,
		Audience:   audience,
		ClientIds:  clientIds,
		CookieName: cookieName,
	}
	return b
}

// WithJWKSConfig sets JWKS configuration for RSA/ECDSA signature validation
func (b *Builder) WithJWKSConfig(jwksURL string, cacheMinutes int, signingMethods []string) *Builder {
	b.jwtConfig.JWKSURL = jwksURL
	b.jwtConfig.JWKSCache = cacheMinutes
	b.jwtConfig.SigningMethods = signingMethods
	return b
}

// WithDevMode enables or disables development mode
func (b *Builder) WithDevMode(enabled bool) *Builder {
	b.devMode = enabled
	return b
}

// WithEnvironment sets the application environment (dev, prod, etc.)
func (b *Builder) WithEnvironment(env string) *Builder {
	b.env = env
	return b
}

// WithLogger sets the logger for auth messages
func (b *Builder) WithLogger(logger *zap.Logger) *Builder {
	b.logger = logger
	return b
}

// Build initializes the authentication system with all configured options
// Returns an error if configuration is invalid or unsafe
func (b *Builder) Build() error {
	// Initialize JWT configuration
	if err := InitJWTConfig(b.jwtConfig); err != nil {
		return fmt.Errorf("failed to initialize JWT config: %w", err)
	}

	// Safety check: prevent dev mode in production
	if b.env == "prod" && b.devMode {
		return fmt.Errorf("FATAL: Dev mode cannot be enabled in production environment")
	}

	// Initialize dev mode
	InitDevMode(b.devMode)

	// Log warnings if dev mode is active
	if b.devMode {
		if b.logger != nil {
			b.logger.Warn("⚠️  AUTH DEV MODE ACTIVE - Do not use in production!")
		}
	}

	// Log successful initialization
	if b.logger != nil {
		if b.devMode {
			b.logger.Info("Authentication initialized",
				zap.String("mode", "development"),
				zap.Bool("dev_mode", true),
				zap.Strings("issuers", b.jwtConfig.Issuer),
				zap.Strings("audiences", b.jwtConfig.Audience),
				zap.Strings("client_ids", b.jwtConfig.ClientIds),
				zap.String("jwks_url", b.jwtConfig.JWKSURL),
				zap.Int("jwks_cache_minutes", b.jwtConfig.JWKSCache),
				zap.Strings("signing_methods", b.jwtConfig.SigningMethods),
			)
		} else {
			b.logger.Info("Authentication initialized",
				zap.String("mode", "production"),
				zap.Bool("dev_mode", false),
				zap.Strings("issuers", b.jwtConfig.Issuer),
				zap.Strings("audiences", b.jwtConfig.Audience),
				zap.Strings("client_ids", b.jwtConfig.ClientIds),
				zap.String("jwks_url", b.jwtConfig.JWKSURL),
				zap.Int("jwks_cache_minutes", b.jwtConfig.JWKSCache),
				zap.Strings("signing_methods", b.jwtConfig.SigningMethods),
			)
		}
	}

	return nil
}

// MustBuild initializes authentication and panics on error
// Use this only when you're certain the configuration is valid
func (b *Builder) MustBuild() {
	if err := b.Build(); err != nil {
		panic(fmt.Sprintf("Failed to initialize authentication: %v", err))
	}
}
