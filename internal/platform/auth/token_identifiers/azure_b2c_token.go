package token_identifiers

import (
	"context"
	"errors"
	"fmt"

	"pharmacy-modernization-project-model/internal/platform/auth/types"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// AzureB2CTokenIdentifier handles Azure B2C token type detection and validation
type AzureB2CTokenIdentifier struct {
	jwksURL        string
	jwksKeyfunc    keyfunc.Keyfunc
	config         types.JWTConfig
	signingMethods []string // Per-type signing methods
	issuer         []string // Per-type issuer validation
	audience       []string // Per-type audience validation
	clientIds      []string // Per-type client ID validation
}

// NewAzureB2CTokenIdentifier creates a new Azure B2C token identifier
func NewAzureB2CTokenIdentifier(tokenConfig types.TokenTypeConfig, config types.JWTConfig) (*AzureB2CTokenIdentifier, error) {
	ctx := context.Background()
	jwksKeyfunc, err := keyfunc.NewDefaultCtx(ctx, []string{tokenConfig.JWKSURL})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS for Azure B2C token: %w", err)
	}

	return &AzureB2CTokenIdentifier{
		jwksURL:        tokenConfig.JWKSURL,
		jwksKeyfunc:    jwksKeyfunc,
		config:         config,
		signingMethods: tokenConfig.SigningMethods,
		issuer:         tokenConfig.Issuer,
		audience:       tokenConfig.Audience,
		clientIds:      tokenConfig.ClientIds,
	}, nil
}

// DetectTokenType attempts to detect if this is an Azure B2C token based on issuer
func (abti *AzureB2CTokenIdentifier) DetectTokenType(ctx context.Context, tokenString string) (types.TokenType, error) {
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
	for _, expectedIssuer := range abti.issuer {
		if issuer == expectedIssuer {
			return types.TokenTypeAzureB2C, nil
		}
	}

	return "", fmt.Errorf("issuer '%s' does not match Azure B2C token issuers", issuer)
}

// IsValidToken validates the token using JWKS
func (abti *AzureB2CTokenIdentifier) IsValidToken(ctx context.Context, tokenString string) (types.TokenType, error) {
	// Parse token using JWKS keyfunc
	token, err := jwt.ParseWithClaims(tokenString, &types.AzureB2CClaims{}, abti.jwksKeyfunc.Keyfunc)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	// Check if signing method is allowed
	if len(abti.signingMethods) > 0 {
		allowed := false
		for _, method := range abti.signingMethods {
			if token.Method.Alg() == method {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", errors.New("signing method not allowed: " + token.Method.Alg())
		}
	}

	return types.TokenTypeAzureB2C, nil
}

// GetTokenType returns the token type this identifier handles
func (abti *AzureB2CTokenIdentifier) GetTokenType() types.TokenType {
	return types.TokenTypeAzureB2C
}

// GetJWKSURL returns the JWKS URL for this token type
func (abti *AzureB2CTokenIdentifier) GetJWKSURL() string {
	return abti.jwksURL
}

// ExtractUser extracts user information from validated token
func (abti *AzureB2CTokenIdentifier) ExtractUser(token *jwt.Token) (*types.User, error) {
	claims, ok := token.Claims.(*types.AzureB2CClaims)
	if !ok {
		return nil, errors.New("invalid claims for Azure B2C token type")
	}

	// Validate issuer if configured
	if len(abti.issuer) > 0 {
		validIssuer := false
		for _, issuer := range abti.issuer {
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
	if len(abti.audience) > 0 {
		validAudience := false
		for _, expectedAud := range abti.audience {
			if claims.Audience == expectedAud {
				validAudience = true
				break
			}
		}
		if !validAudience {
			return nil, errors.New("invalid audience")
		}
	}

	// Build user from Azure B2C claims
	user := &types.User{
		ID:              claims.ObjectId, // Use Azure B2C object ID as user ID
		Email:           claims.Email,
		Name:            claims.Name,
		Permissions:     claims.Roles,                                     // Map Azure roles to permissions
		DataAccessRoles: claims.Groups,                                    // Map Azure groups to data access roles
		FuncRoles:       convertAzureRolesToFuncRoles(claims.CustomRoles), // Convert custom roles
	}

	return user, nil
}

// convertAzureRolesToFuncRoles converts Azure B2C custom roles to FuncRole format
func convertAzureRolesToFuncRoles(customRoles []string) []types.FuncRole {
	funcRoles := make([]types.FuncRole, len(customRoles))
	for i, role := range customRoles {
		funcRoles[i] = types.FuncRole{
			RoleName: role,
		}
	}
	return funcRoles
}
