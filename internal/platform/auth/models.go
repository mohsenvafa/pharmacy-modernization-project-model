package auth

import "github.com/golang-jwt/jwt/v5"

// User represents an authenticated user extracted from JWT
type User struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// TokenSource defines where to extract the JWT token from
type TokenSource int

const (
	TokenSourceAuto   TokenSource = iota // Auto-detect (tries header, then cookie)
	TokenSourceHeader                    // Only Authorization header
	TokenSourceCookie                    // Only cookie
)
