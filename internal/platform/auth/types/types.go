package types

import (
	"github.com/golang-jwt/jwt/v5"
)

// TokenType represents the type of JWT token
type TokenType string

const (
	TokenTypeAuthPass TokenType = "auth_pass" // Your custom token format
	TokenTypeAzureB2C TokenType = "azure_b2c" // Azure B2C token format
)

// User represents an authenticated user
type User struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	Name            string     `json:"name"`
	Permissions     []string   `json:"permissions"`
	DataAccessRoles []string   `json:"dataAccessRoles"`
	FuncRoles       []FuncRole `json:"func-roles"`
}

// FuncRole represents a functional role
type FuncRole struct {
	RoleName string `json:"roleName"`
}

// HasFuncRole checks if the user has a specific functional role
func (u *User) HasFuncRole(roleName string) bool {
	for _, role := range u.FuncRoles {
		if role.RoleName == roleName {
			return true
		}
	}
	return false
}

// GetFuncRoleNames returns all functional role names
func (u *User) GetFuncRoleNames() []string {
	roleNames := make([]string, len(u.FuncRoles))
	for i, role := range u.FuncRoles {
		roleNames[i] = role.RoleName
	}
	return roleNames
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	CookieName       string                        // Name of the authentication cookie
	TokenTypesConfig map[TokenType]TokenTypeConfig // Configuration for each token type
	JWKSCache        int                           // Cache duration in minutes (default: 15)
	TokenTypes       []TokenType                   // Supported token types (auth_pass, azure_b2c)
}

// TokenTypeConfig represents configuration for a specific token type
type TokenTypeConfig struct {
	JWKSURL        string   // JWKS URL for this token type
	SigningMethods []string // Allowed signing methods for this token type
	Issuer         []string // Expected token issuers for this token type
	Audience       []string // Expected token audiences for this token type
	ClientIds      []string // Expected client IDs for this token type
}

// JWTClaims represents the claims structure for default tokens
type JWTClaims struct {
	UserID          string     `json:"user_id"`
	Email           string     `json:"email"`
	Name            string     `json:"name"`
	Permissions     []string   `json:"permissions"`
	ClientId        string     `json:"client_id"`
	SubjectIdType   string     `json:"subjectIdType"`
	DataAccessRoles []string   `json:"dataAccessRoles"`
	FuncRoles       []FuncRole `json:"func-roles"`
	jwt.RegisteredClaims
}

// AzureB2CClaims represents the claims structure for Azure B2C tokens
type AzureB2CClaims struct {
	// Standard Azure B2C claims
	Audience  string `json:"aud"`
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	NotBefore int64  `json:"nbf"`

	// Azure B2C specific claims
	Email          string `json:"email"`
	Name           string `json:"name"`
	GivenName      string `json:"given_name"`
	FamilyName     string `json:"family_name"`
	ObjectId       string `json:"oid"`
	TenantId       string `json:"tid"`
	AppId          string `json:"appid"`
	AppIdAcr       string `json:"appidacr"`
	Idp            string `json:"idp"`
	IdpAccessToken string `json:"idp_access_token"`

	// Custom claims (you can add more based on your Azure B2C setup)
	Roles       []string `json:"roles"`
	Groups      []string `json:"groups"`
	CustomRoles []string `json:"custom_roles"`

	jwt.RegisteredClaims
}
