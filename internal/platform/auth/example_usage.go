package auth

import (
	"net/http"
)

// Example usage of the new JWT claims fields

// ExampleHandler shows how to use the new DataAccessRoles and FuncRoles
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, err := GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user has a specific functional role
	if user.HasFuncRole("NONPROD_RxIntakeUITechnicianAccess") {
		// User has technician access
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Technician access granted"))
		return
	}

	// Check data access roles
	for _, role := range user.DataAccessRoles {
		if role == "PHARMACY_DATA" {
			// User has pharmacy data access
			break
		}
	}

	// Get all functional role names
	roleNames := user.GetFuncRoleNames()
	// roleNames will contain: ["NONPROD_RxIntakeUITechnicianAccess"]

	// Check permissions (existing functionality)
	for _, permission := range user.Permissions {
		if permission == "patient:read" {
			// User can read patient data
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Access granted"))
}

// ExampleMiddleware shows how to create middleware based on functional roles
func RequireFuncRole(roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := GetCurrentUser(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !user.HasFuncRole(roleName) {
				http.Error(w, "Forbidden: Missing required role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Example: Use the middleware
// router.Use(RequireFuncRole("NONPROD_RxIntakeUITechnicianAccess"))
