package auth

import (
	"pharmacy-modernization-project-model/internal/platform/auth/types"
)

// User represents an authenticated user
type User = types.User

// FuncRole represents a functional role
type FuncRole = types.FuncRole

// JWTClaims represents the claims structure for default tokens
type JWTClaims = types.JWTClaims

// AzureB2CClaims represents the claims structure for Azure B2C tokens
type AzureB2CClaims = types.AzureB2CClaims

// TokenType represents the type of JWT token
type TokenType = types.TokenType

// TokenTypeConfig represents configuration for a specific token type
type TokenTypeConfig = types.TokenTypeConfig

const (
	TokenTypeAuthPass = types.TokenTypeAuthPass
	TokenTypeAzureB2C = types.TokenTypeAzureB2C
)

// TokenSource defines where to extract the JWT token from
type TokenSource int

const (
	TokenSourceAuto   TokenSource = iota // Auto-detect (tries header, then cookie)
	TokenSourceHeader                    // Only Authorization header
	TokenSourceCookie                    // Only cookie
)
