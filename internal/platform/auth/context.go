package auth

import (
	"context"
	"errors"
)

type contextKey string

const userContextKey contextKey = "auth_user"

// SetUser stores the authenticated user in the context
func SetUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetCurrentUser retrieves the authenticated user from the context
func GetCurrentUser(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok || user == nil {
		return nil, errors.New("no authenticated user in context")
	}
	return user, nil
}

// MustGetCurrentUser retrieves the user from context and panics if not found
// Use this only in handlers where authentication is guaranteed by middleware
func MustGetCurrentUser(ctx context.Context) *User {
	user, err := GetCurrentUser(ctx)
	if err != nil {
		panic("expected authenticated user but none found")
	}
	return user
}
