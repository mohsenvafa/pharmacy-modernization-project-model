package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

// DevMode configuration
var devModeEnabled bool
var mockUsers map[string]*User

// InitDevMode initializes development mode with mock users
func InitDevMode(enabled bool) {
	devModeEnabled = enabled
	if enabled {
		log.Println("⚠️  AUTH DEV MODE ENABLED - Security bypassed with mock users")
		initializeMockUsers()
	}
}

// IsDevModeEnabled returns true if dev mode is active
func IsDevModeEnabled() bool {
	return devModeEnabled
}

// initializeMockUsers creates default mock users for development
func initializeMockUsers() {
	mockUsers = map[string]*User{
		"admin": {
			ID:    "mock-admin-001",
			Email: "admin@dev.local",
			Name:  "Mock Admin",
			Permissions: []string{
				"admin:all",
			},
		},
		"doctor": {
			ID:    "mock-doctor-001",
			Email: "doctor@dev.local",
			Name:  "Dr. Dev",
			Permissions: []string{
				"patient:read",
				"patient:write",
				"prescription:read",
				"prescription:write",
				"prescription:approve",
				"doctor:role",
				"dashboard:view",
			},
		},
		"pharmacist": {
			ID:    "mock-pharmacist-001",
			Email: "pharmacist@dev.local",
			Name:  "Dev Pharmacist",
			Permissions: []string{
				"patient:read",
				"prescription:read",
				"prescription:dispense",
				"pharmacist:role",
				"dashboard:view",
			},
		},
		"nurse": {
			ID:    "mock-nurse-001",
			Email: "nurse@dev.local",
			Name:  "Dev Nurse",
			Permissions: []string{
				"patient:read",
				"prescription:read",
				"nurse:role",
			},
		},
		"readonly": {
			ID:    "mock-readonly-001",
			Email: "readonly@dev.local",
			Name:  "Dev Readonly User",
			Permissions: []string{
				"patient:read",
				"prescription:read",
				"dashboard:view",
			},
		},
	}
}

// AddMockUser allows adding custom mock users for testing
func AddMockUser(key string, user *User) {
	if !devModeEnabled {
		log.Println("Warning: Cannot add mock user - dev mode not enabled")
		return
	}
	if mockUsers == nil {
		mockUsers = make(map[string]*User)
	}
	mockUsers[key] = user
	log.Printf("✓ Added mock user: %s (%s)", key, user.Name)
}

// GetMockUser retrieves a mock user by key
func GetMockUser(key string) *User {
	if !devModeEnabled {
		return nil
	}
	return mockUsers[key]
}

// ListMockUsers returns all available mock users
func ListMockUsers() map[string]*User {
	if !devModeEnabled {
		return nil
	}
	return mockUsers
}

// DevAuthMiddleware bypasses real authentication in dev mode
// Uses mock users based on X-Mock-User header or defaults to admin
func DevAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !devModeEnabled {
				// Dev mode not enabled, continue to real auth
				next.ServeHTTP(w, r)
				return
			}

			// Check for mock user selection via header
			mockUserKey := r.Header.Get("X-Mock-User")
			if mockUserKey == "" {
				// Default to admin user
				mockUserKey = "admin"
			}

			user := mockUsers[mockUserKey]
			if user == nil {
				log.Printf("Warning: Unknown mock user '%s', using admin", mockUserKey)
				user = mockUsers["admin"]
			}

			log.Printf("DEV AUTH: Using mock user '%s' (%s) with permissions: %v",
				mockUserKey, user.Email, user.Permissions)

			// Set user in context
			ctx := SetUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// DevAuthInfo returns information about dev mode and available users
func DevAuthInfo(w http.ResponseWriter, r *http.Request) {
	if !devModeEnabled {
		http.Error(w, "Dev mode not enabled", http.StatusNotFound)
		return
	}

	info := map[string]interface{}{
		"dev_mode": true,
		"message":  "Development mode is ACTIVE - Authentication bypassed",
		"mock_users": func() []map[string]interface{} {
			users := []map[string]interface{}{}
			for key, user := range mockUsers {
				users = append(users, map[string]interface{}{
					"key":         key,
					"id":          user.ID,
					"email":       user.Email,
					"name":        user.Name,
					"permissions": user.Permissions,
				})
			}
			return users
		}(),
		"usage": map[string]string{
			"header":       "X-Mock-User: <key>",
			"default":      "admin",
			"example_curl": "curl -H 'X-Mock-User: doctor' http://localhost:8080/patients",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// RequireAuthWithDevMode is a wrapper that uses dev mode if enabled, otherwise real auth
func RequireAuthWithDevMode() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if devModeEnabled {
				// Use dev mode mock auth
				DevAuthMiddleware()(next).ServeHTTP(w, r)
			} else {
				// Use real JWT auth
				RequireAuth()(next).ServeHTTP(w, r)
			}
		})
	}
}
