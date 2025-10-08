package auth

import "context"

// hasPermission checks if user has a single permission
func hasPermission(userPermissions []string, required string) bool {
	for _, perm := range userPermissions {
		if perm == required {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if user has ALL of the required permissions
func HasAllPermissions(userPermissions []string, required []string) bool {
	for _, req := range required {
		if !hasPermission(userPermissions, req) {
			return false
		}
	}
	return true
}

// HasAnyPermission checks if user has ANY of the required permissions
func HasAnyPermission(userPermissions []string, required []string) bool {
	for _, req := range required {
		if hasPermission(userPermissions, req) {
			return true
		}
	}
	return false
}

// HasAllPermissionsCtx checks if the current user in context has ALL permissions
func HasAllPermissionsCtx(ctx context.Context, required []string) bool {
	user, err := GetCurrentUser(ctx)
	if err != nil {
		return false
	}
	return HasAllPermissions(user.Permissions, required)
}

// HasAnyPermissionCtx checks if the current user in context has ANY of the permissions
func HasAnyPermissionCtx(ctx context.Context, required []string) bool {
	user, err := GetCurrentUser(ctx)
	if err != nil {
		return false
	}
	return HasAnyPermission(user.Permissions, required)
}
