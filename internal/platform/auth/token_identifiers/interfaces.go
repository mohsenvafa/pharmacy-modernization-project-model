package token_identifiers

import (
	"context"

	"pharmacy-modernization-project-model/internal/platform/auth/types"

	"github.com/golang-jwt/jwt/v5"
)

// TokenTypeIdentifier defines the interface for token type detection and validation
type TokenTypeIdentifier interface {
	// DetectTokenType attempts to detect the token type from a token string
	DetectTokenType(ctx context.Context, tokenString string) (types.TokenType, error)

	// IsValidToken validates the token and returns the token type if valid
	IsValidToken(ctx context.Context, tokenString string) (types.TokenType, error)

	// GetTokenType returns the token type this identifier handles
	GetTokenType() types.TokenType

	// GetJWKSURL returns the JWKS URL for this token type
	GetJWKSURL() string

	// ExtractUser extracts user information from validated token
	ExtractUser(token *jwt.Token) (*types.User, error)
}
