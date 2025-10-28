package auth

import "github.com/golang-jwt/jwt/v5"

// User represents an authenticated user extracted from JWT
type User struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	Name            string     `json:"name"`
	Permissions     []string   `json:"permissions"`
	DataAccessRoles []string   `json:"dataAccessRoles"`
	FuncRoles       []FuncRole `json:"func-roles"`
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

// GetFuncRoleNames returns a slice of all functional role names
func (u *User) GetFuncRoleNames() []string {
	roleNames := make([]string, len(u.FuncRoles))
	for i, role := range u.FuncRoles {
		roleNames[i] = role.RoleName
	}
	return roleNames
}

// FuncRole represents a functional role in the JWT claims
type FuncRole struct {
	RoleName string `json:"roleName"`
}

// JWTClaims represents the JWT token claims
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

// TokenSource defines where to extract the JWT token from
type TokenSource int

const (
	TokenSourceAuto   TokenSource = iota // Auto-detect (tries header, then cookie)
	TokenSourceHeader                    // Only Authorization header
	TokenSourceCookie                    // Only cookie
)
