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
func (b *Builder) WithJWTConfig(cookieName string) *Builder {
	b.jwtConfig = JWTConfig{
		CookieName: cookieName,
	}
	return b
}

// WithTokenTypesConfig sets configuration for each token type
func (b *Builder) WithTokenTypesConfig(tokenTypesConfig map[TokenType]TokenTypeConfig, cacheMinutes int) *Builder {
	b.jwtConfig.TokenTypesConfig = tokenTypesConfig
	b.jwtConfig.JWKSCache = cacheMinutes
	return b
}

// WithTokenTypes sets supported token types
func (b *Builder) WithTokenTypes(tokenTypes []TokenType) *Builder {
	b.jwtConfig.TokenTypes = tokenTypes
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
			b.logger.Info("Authentication initialized")
		} else {
			b.logger.Info("Authentication initialized")
		}
	}

	return nil
}
