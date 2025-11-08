package token_identifiers

import (
	"context"
	"errors"
	"fmt"

	"pharmacy-modernization-project-model/internal/platform/auth/types"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// DefaultTokenIdentifier handles generic token type detection and validation
type DefaultTokenIdentifier struct {
	tokenType      types.TokenType
	jwksURL        string
	jwksKeyfunc    keyfunc.Keyfunc
	signingMethods []string
	issuer         []string
	audience       []string
	clientIds      []string
}

// NewDefaultTokenIdentifier creates a new default token identifier for the supplied token type
func NewDefaultTokenIdentifier(tokenType types.TokenType, tokenConfig types.TokenTypeConfig, _ types.JWTConfig) (*DefaultTokenIdentifier, error) {
	ctx := context.Background()
	jwksKeyfunc, err := keyfunc.NewDefaultCtx(ctx, []string{tokenConfig.JWKSURL})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS for token type %s: %w", tokenType, err)
	}

	return &DefaultTokenIdentifier{
		tokenType:      tokenType,
		jwksURL:        tokenConfig.JWKSURL,
		jwksKeyfunc:    jwksKeyfunc,
		signingMethods: tokenConfig.SigningMethods,
		issuer:         tokenConfig.Issuer,
		audience:       tokenConfig.Audience,
		clientIds:      tokenConfig.ClientIds,
	}, nil
}

// DetectTokenType attempts to detect if this is the configured token based on issuer claim
func (dti *DefaultTokenIdentifier) DetectTokenType(ctx context.Context, tokenString string) (types.TokenType, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if len(dti.issuer) == 0 {
		return dti.tokenType, nil
	}

	claims := token.Claims.(jwt.MapClaims)
	issuer, ok := claims["iss"].(string)
	if !ok {
		return "", fmt.Errorf("token does not contain issuer claim")
	}

	for _, expectedIssuer := range dti.issuer {
		if issuer == expectedIssuer {
			return dti.tokenType, nil
		}
	}

	return "", fmt.Errorf("issuer '%s' does not match configured issuers for token type %s", issuer, dti.tokenType)
}

// IsValidToken validates the token using JWKS
func (dti *DefaultTokenIdentifier) IsValidToken(ctx context.Context, tokenString string) (types.TokenType, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, dti.jwksKeyfunc.Keyfunc)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	if len(dti.signingMethods) > 0 {
		allowed := false
		for _, method := range dti.signingMethods {
			if token.Method.Alg() == method {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", fmt.Errorf("signing method not allowed: %s", token.Method.Alg())
		}
	}

	return dti.tokenType, nil
}

// GetTokenType returns the token type this identifier handles
func (dti *DefaultTokenIdentifier) GetTokenType() types.TokenType {
	return dti.tokenType
}

// GetJWKSURL returns the JWKS URL for this token type
func (dti *DefaultTokenIdentifier) GetJWKSURL() string {
	return dti.jwksURL
}

// ExtractUser extracts user information from validated token
func (dti *DefaultTokenIdentifier) ExtractUser(token *jwt.Token) (*types.User, error) {
	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		return nil, errors.New("invalid claims for token type " + string(dti.tokenType))
	}

	if len(dti.issuer) > 0 {
		validIssuer := false
		for _, issuer := range dti.issuer {
			if claims.Issuer == issuer {
				validIssuer = true
				break
			}
		}
		if !validIssuer {
			return nil, errors.New("invalid issuer")
		}
	}

	if len(dti.audience) > 0 {
		validAudience := false
		for _, aud := range claims.Audience {
			for _, expectedAud := range dti.audience {
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
			return nil, fmt.Errorf("invalid audience for token type %s: expected one of %v, got %v", dti.tokenType, dti.audience, claims.Audience)
		}
	}

	if len(dti.clientIds) > 0 {
		validClientId := false
		for _, clientId := range dti.clientIds {
			if claims.ClientId == clientId {
				validClientId = true
				break
			}
		}
		if !validClientId {
			return nil, fmt.Errorf("invalid client ID for token type %s: expected one of %v, got %s", dti.tokenType, dti.clientIds, claims.ClientId)
		}
	}

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
