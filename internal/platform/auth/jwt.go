package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"pharmacy-modernization-project-model/internal/platform/auth/token_identifiers"
	"pharmacy-modernization-project-model/internal/platform/auth/types"
)

// JWTConfig holds JWT configuration
type JWTConfig = types.JWTConfig

var jwtConfig JWTConfig
var tokenManager *TokenManager

// InitJWTConfig initializes the global JWT configuration and token manager
func InitJWTConfig(config JWTConfig) error {
	jwtConfig = config

	// Create token manager
	tokenManager = NewTokenManager(config)

	// Register token identifiers for each configured token type
	for _, tokenType := range config.TokenTypes {
		tokenConfig, exists := config.TokenTypesConfig[tokenType]
		if !exists {
			return fmt.Errorf("configuration not found for token type: %s", string(tokenType))
		}

		switch tokenType {
		case types.TokenTypeAuthPass:
			identifier, err := token_identifiers.NewAuthPassTokenIdentifier(tokenConfig, config)
			if err != nil {
				return fmt.Errorf("failed to create auth_pass token identifier: %w", err)
			}
			tokenManager.RegisterIdentifier(identifier)

		case types.TokenTypeAzureB2C:
			identifier, err := token_identifiers.NewAzureB2CTokenIdentifier(tokenConfig, config)
			if err != nil {
				return fmt.Errorf("failed to create Azure B2C token identifier: %w", err)
			}
			tokenManager.RegisterIdentifier(identifier)

		default:
			identifier, err := token_identifiers.NewDefaultTokenIdentifier(tokenType, tokenConfig, config)
			if err != nil {
				return fmt.Errorf("failed to create default token identifier for token type %s: %w", string(tokenType), err)
			}
			tokenManager.RegisterIdentifier(identifier)
		}
	}

	return nil
}

// ValidateToken validates and parses a JWT token string using the token manager
func ValidateToken(tokenString string) (*types.User, error) {
	if tokenManager == nil {
		return nil, errors.New("token manager not initialized")
	}

	ctx := context.Background()
	return tokenManager.ValidateToken(ctx, tokenString)
}

// ExtractToken extracts JWT token from request based on source
func ExtractToken(r *http.Request, source TokenSource) (string, error) {
	switch source {
	case TokenSourceHeader:
		return extractFromHeaderOnly(r)
	case TokenSourceCookie:
		return extractFromCookieOnly(r)
	case TokenSourceAuto:
		return extractFromAuto(r)
	default:
		return extractFromAuto(r)
	}
}

// extractFromHeaderOnly extracts token from Authorization header only
func extractFromHeaderOnly(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no Authorization header found")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format, expected: Bearer <token>")
	}

	return parts[1], nil
}

// extractFromCookieOnly extracts token from cookie only
func extractFromCookieOnly(r *http.Request) (string, error) {
	cookieName := jwtConfig.CookieName
	if cookieName == "" {
		cookieName = "auth_token" // default
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", errors.New("no auth cookie found")
	}

	if cookie.Value == "" {
		return "", errors.New("auth cookie is empty")
	}

	return cookie.Value, nil
}

// extractFromAuto tries header first, then cookie
func extractFromAuto(r *http.Request) (string, error) {
	// Try header first
	if token, err := extractFromHeaderOnly(r); err == nil {
		return token, nil
	}

	// Fall back to cookie
	if token, err := extractFromCookieOnly(r); err == nil {
		return token, nil
	}

	return "", errors.New("no token found in Authorization header or cookie")
}
