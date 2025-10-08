package auth

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// AuthDirective implements @auth directive for GraphQL
// Returns UNAUTHENTICATED error (401) if user is not authenticated
func AuthDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		user, err := GetCurrentUser(ctx)
		if err != nil || user == nil {
			return nil, &gqlerror.Error{
				Message: "Unauthenticated",
				Extensions: map[string]interface{}{
					"code":   "UNAUTHENTICATED",
					"status": 401,
				},
			}
		}

		return next(ctx)
	}
}

// PermissionAnyDirective implements @permissionAny directive for GraphQL
// Returns FORBIDDEN error (403) if user lacks any of the required permissions
func PermissionAnyDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver, requires []string) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, requires []string) (interface{}, error) {
		user, err := GetCurrentUser(ctx)
		if err != nil || user == nil {
			return nil, &gqlerror.Error{
				Message: "Unauthenticated",
				Extensions: map[string]interface{}{
					"code":   "UNAUTHENTICATED",
					"status": 401,
				},
			}
		}

		if !HasAnyPermission(user.Permissions, requires) {
			return nil, &gqlerror.Error{
				Message: fmt.Sprintf("Forbidden: requires at least one of: %v", requires),
				Extensions: map[string]interface{}{
					"code":                 "FORBIDDEN",
					"status":               403,
					"required_permissions": requires,
					"match":                "any",
				},
			}
		}

		return next(ctx)
	}
}

// PermissionAllDirective implements @permissionAll directive for GraphQL
// Returns FORBIDDEN error (403) if user lacks all of the required permissions
func PermissionAllDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver, requires []string) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, requires []string) (interface{}, error) {
		user, err := GetCurrentUser(ctx)
		if err != nil || user == nil {
			return nil, &gqlerror.Error{
				Message: "Unauthenticated",
				Extensions: map[string]interface{}{
					"code":   "UNAUTHENTICATED",
					"status": 401,
				},
			}
		}

		if !HasAllPermissions(user.Permissions, requires) {
			return nil, &gqlerror.Error{
				Message: fmt.Sprintf("Forbidden: requires all of: %v", requires),
				Extensions: map[string]interface{}{
					"code":                 "FORBIDDEN",
					"status":               403,
					"required_permissions": requires,
					"match":                "all",
				},
			}
		}

		return next(ctx)
	}
}
