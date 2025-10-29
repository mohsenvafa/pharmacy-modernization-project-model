package token_identifiers

import (
	"context"
	"errors"
	"fmt"

	"pharmacy-modernization-project-model/internal/platform/auth/types"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// AuthPassTokenIdentifier handles auth_pass token type detection and validation
type AuthPassTokenIdentifier struct {
	jwksURL        string
	jwksKeyfunc    keyfunc.Keyfunc
	config         types.JWTConfig
	signingMethods []string // Per-type signing methods
	issuer         []string // Per-type issuer validation
	audience       []string // Per-type audience validation
	clientIds      []string // Per-type client ID validation
}

// NewAuthPassTokenIdentifier creates a new auth_pass token identifier
func NewAuthPassTokenIdentifier(tokenConfig types.TokenTypeConfig, config types.JWTConfig) (*AuthPassTokenIdentifier, error) {
	ctx := context.Background()
	jwksKeyfunc, err := keyfunc.NewDefaultCtx(ctx, []string{tokenConfig.JWKSURL})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS for auth_pass token: %w", err)
	}

	return &AuthPassTokenIdentifier{
		jwksURL:        tokenConfig.JWKSURL,
		jwksKeyfunc:    jwksKeyfunc,
		config:         config,
		signingMethods: tokenConfig.SigningMethods,
		issuer:         tokenConfig.Issuer,
		audience:       tokenConfig.Audience,
		clientIds:      tokenConfig.ClientIds,
	}, nil
}

// DetectTokenType attempts to detect if this is an auth_pass token based on issuer
func (apti *AuthPassTokenIdentifier) DetectTokenType(ctx context.Context, tokenString string) (types.TokenType, error) {
	// Parse token without validation to extract issuer
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("dummy"), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	issuer, ok := claims["iss"].(string)
	if !ok {
		return "", fmt.Errorf("token does not contain issuer claim")
	}

	// Check if issuer matches any of the configured issuers for this token type
	for _, expectedIssuer := range apti.issuer {
		if issuer == expectedIssuer {
			return types.TokenTypeAuthPass, nil
		}
	}

	return "", fmt.Errorf("issuer '%s' does not match auth_pass token issuers", issuer)
}

// IsValidToken validates the token using JWKS
func (apti *AuthPassTokenIdentifier) IsValidToken(ctx context.Context, tokenString string) (types.TokenType, error) {
	// Parse token using JWKS keyfunc
	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, apti.jwksKeyfunc.Keyfunc)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	// Check if signing method is allowed
	if len(apti.signingMethods) > 0 {
		allowed := false
		for _, method := range apti.signingMethods {
			if token.Method.Alg() == method {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", errors.New("signing method not allowed: " + token.Method.Alg())
		}
	}

	return types.TokenTypeAuthPass, nil
}

// GetTokenType returns the token type this identifier handles
func (apti *AuthPassTokenIdentifier) GetTokenType() types.TokenType {
	return types.TokenTypeAuthPass
}

// GetJWKSURL returns the JWKS URL for this token type
func (apti *AuthPassTokenIdentifier) GetJWKSURL() string {
	return apti.jwksURL
}

// ExtractUser extracts user information from validated token
func (apti *AuthPassTokenIdentifier) ExtractUser(token *jwt.Token) (*types.User, error) {
	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		return nil, errors.New("invalid claims for auth_pass token type")
	}

	// Validate issuer if configured
	if len(apti.issuer) > 0 {
		validIssuer := false
		for _, issuer := range apti.issuer {
			if claims.Issuer == issuer {
				validIssuer = true
				break
			}
		}
		if !validIssuer {
			return nil, errors.New("invalid issuer")
		}
	}

	// Validate audience if configured
	if len(apti.audience) > 0 {
		validAudience := false
		for _, aud := range claims.Audience {
			for _, expectedAud := range apti.audience {
				if aud == expectedAud {
					validAudience = true
					break
				}
			}
			if validAudience {
				break
			}
		}
		if !validAudience {
			return nil, fmt.Errorf("invalid audience for auth_pass token type: expected one of %v, got %v", apti.audience, claims.Audience)
		}
	}

	// Validate client ID if configured
	if len(apti.clientIds) > 0 {
		validClientId := false
		for _, clientId := range apti.clientIds {
			if claims.ClientId == clientId {
				validClientId = true
				break
			}
		}
		if !validClientId {
			return nil, fmt.Errorf("invalid client ID for auth_pass token type: expected one of %v, got %s", apti.clientIds, claims.ClientId)
		}
	}

	// Build user from claims
	user := &types.User{
		ID:              claims.UserID,
		Email:           claims.Email,
		Name:            claims.Name,
		Permissions:     claims.Permissions,
		DataAccessRoles: claims.DataAccessRoles,
		FuncRoles:       claims.FuncRoles,
	}

	return user, nil
}
