package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Issuer     []string
	Audience   []string
	CookieName string
}

var jwtConfig JWTConfig

// InitJWTConfig initializes the global JWT configuration
func InitJWTConfig(config JWTConfig) {
	jwtConfig = config
}

// ValidateToken validates and parses a JWT token string
func ValidateToken(tokenString string) (*User, error) {
	if jwtConfig.Secret == "" {
		return nil, errors.New("JWT config not initialized")
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

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

	// Build user from claims
	user := &User{
		ID:          claims.UserID,
		Email:       claims.Email,
		Name:        claims.Name,
		Permissions: claims.Permissions,
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

// CreateToken creates a new JWT token for a user (helper for testing/login)
func CreateToken(user *User, expirationHours int) (string, error) {
	if jwtConfig.Secret == "" {
		return "", errors.New("JWT config not initialized")
	}

	// Use first issuer and audience from arrays, or empty strings if arrays are empty
	var issuer string
	var audience jwt.ClaimStrings

	if len(jwtConfig.Issuer) > 0 {
		issuer = jwtConfig.Issuer[0]
	}

	if len(jwtConfig.Audience) > 0 {
		audience = jwt.ClaimStrings(jwtConfig.Audience)
	}

	claims := JWTClaims{
		UserID:      user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Permissions: user.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.Secret))
}
