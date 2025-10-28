package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Issuer         []string // Expected token issuers
	Audience       []string // Expected token audiences
	ClientIds      []string // Expected client IDs
	CookieName     string   // Name of the authentication cookie
	JWKSURL        string   // URL to fetch JWKS from
	JWKSCache      int      // Cache duration in minutes (default: 15)
	SigningMethods []string // Allowed signing methods (RS256, ES256, etc.)
}

var jwtConfig JWTConfig
var jwksKeyfunc keyfunc.Keyfunc

// InitJWTConfig initializes the global JWT configuration
func InitJWTConfig(config JWTConfig) error {
	jwtConfig = config

	// Initialize JWKS
	if config.JWKSURL != "" {
		ctx := context.Background()

		var err error
		jwksKeyfunc, err = keyfunc.NewDefaultCtx(ctx, []string{config.JWKSURL})
		if err != nil {
			return fmt.Errorf("failed to initialize JWKS: %w", err)
		}
	} else {
		return errors.New("JWKS URL is required")
	}

	return nil
}

// ValidateToken validates and parses a JWT token string using JWKS
func ValidateToken(tokenString string) (*User, error) {
	// Parse token using JWKS keyfunc
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, jwksKeyfunc.Keyfunc)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if signing method is allowed
	if len(jwtConfig.SigningMethods) > 0 {
		allowed := false
		for _, method := range jwtConfig.SigningMethods {
			if token.Method.Alg() == method {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, errors.New("signing method not allowed: " + token.Method.Alg())
		}
	}

	return extractUserFromClaims(token)
}

// extractUserFromClaims extracts user information from validated token claims
func extractUserFromClaims(token *jwt.Token) (*User, error) {
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Validate issuer if configured
	if len(jwtConfig.Issuer) > 0 {
		validIssuer := false
		for _, issuer := range jwtConfig.Issuer {
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
	if len(jwtConfig.Audience) > 0 {
		validAudience := false
		for _, aud := range claims.Audience {
			for _, expectedAud := range jwtConfig.Audience {
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
			return nil, errors.New("invalid audience")
		}
	}

	// Validate client ID if configured
	if len(jwtConfig.ClientIds) > 0 {
		validClientId := false
		for _, clientId := range jwtConfig.ClientIds {
			if claims.ClientId == clientId {
				validClientId = true
				break
			}
		}
		if !validClientId {
			return nil, errors.New("invalid client ID")
		}
	}

	// Build user from claims
	user := &User{
		ID:              claims.UserID,
		Email:           claims.Email,
		Name:            claims.Name,
		Permissions:     claims.Permissions,
		DataAccessRoles: claims.DataAccessRoles,
		FuncRoles:       claims.FuncRoles,
	}

	return user, nil
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
