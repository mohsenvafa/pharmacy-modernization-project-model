package auth

import (
	"context"
	"fmt"

	"pharmacy-modernization-project-model/internal/platform/auth/token_identifiers"
	"pharmacy-modernization-project-model/internal/platform/auth/types"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager manages different token type identifiers
type TokenManager struct {
	identifiers map[types.TokenType]token_identifiers.TokenTypeIdentifier
	config      types.JWTConfig
}

// NewTokenManager creates a new token manager
func NewTokenManager(config types.JWTConfig) *TokenManager {
	return &TokenManager{
		identifiers: make(map[types.TokenType]token_identifiers.TokenTypeIdentifier),
		config:      config,
	}
}

// RegisterIdentifier registers a token type identifier
func (tm *TokenManager) RegisterIdentifier(identifier token_identifiers.TokenTypeIdentifier) {
	tm.identifiers[identifier.GetTokenType()] = identifier
}

// DetectTokenType attempts to detect the token type using all registered identifiers
func (tm *TokenManager) DetectTokenType(ctx context.Context, tokenString string) (types.TokenType, error) {
	for _, identifier := range tm.identifiers {
		tokenType, err := identifier.DetectTokenType(ctx, tokenString)
		if err == nil {
			return tokenType, nil
		}
	}
	return "", fmt.Errorf("unable to detect token type")
}

// ValidateToken validates a token using the appropriate identifier
func (tm *TokenManager) ValidateToken(ctx context.Context, tokenString string) (*types.User, error) {
	// First detect the token type
	tokenType, err := tm.DetectTokenType(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to detect token type: %w", err)
	}

	// Check if this token type is supported
	if len(tm.config.TokenTypes) > 0 {
		supported := false
		for _, supportedType := range tm.config.TokenTypes {
			if tokenType == supportedType {
				supported = true
				break
			}
		}
		if !supported {
			return nil, fmt.Errorf("token type not supported: %s", string(tokenType))
		}
	}

	// Get the appropriate identifier
	identifier, exists := tm.identifiers[tokenType]
	if !exists {
		return nil, fmt.Errorf("no identifier registered for token type: %s", string(tokenType))
	}

	// Validate the token
	validatedTokenType, err := identifier.IsValidToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if validatedTokenType != tokenType {
		return nil, fmt.Errorf("token type mismatch: expected %s, got %s", string(tokenType), string(validatedTokenType))
	}

	// Parse the token to extract user information
	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// This will be handled by the specific identifier
		return []byte("dummy"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract user information
	return identifier.ExtractUser(token)
}

// GetSupportedTokenTypes returns all supported token types
func (tm *TokenManager) GetSupportedTokenTypes() []types.TokenType {
	var types []types.TokenType
	for tokenType := range tm.identifiers {
		types = append(types, tokenType)
	}
	return types
}

// IsTokenTypeSupported checks if a token type is supported
func (tm *TokenManager) IsTokenTypeSupported(tokenType types.TokenType) bool {
	_, exists := tm.identifiers[tokenType]
	return exists
}

// DetectTokenType is a convenience function that uses the global token manager
func DetectTokenType(tokenString string) (types.TokenType, error) {
	if tokenManager == nil {
		return "", fmt.Errorf("token manager not initialized")
	}
	return tokenManager.DetectTokenType(context.Background(), tokenString)
}
