package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// ==========================================
// AUTHENTICATION MIDDLEWARE (401 errors)
// ==========================================

// RequireAuth validates JWT from header or cookie (auto-detect)
func RequireAuth() func(http.Handler) http.Handler {
	return requireAuthWithSource(TokenSourceAuto)
}

// RequireAuthFromHeader validates JWT from Authorization header only
func RequireAuthFromHeader() func(http.Handler) http.Handler {
	return requireAuthWithSource(TokenSourceHeader)
}

// RequireAuthFromCookie validates JWT from cookie only
func RequireAuthFromCookie() func(http.Handler) http.Handler {
	return requireAuthWithSource(TokenSourceCookie)
}

// requireAuthWithSource is the internal implementation
func requireAuthWithSource(source TokenSource) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := ExtractToken(r, source)
			if err != nil {
				log.Printf("AUTH 401: %s %s - No token: %v", r.Method, r.URL.Path, err)
				handleUnauthorized(w, r, "Authentication required", source)
				return
			}

			user, err := ValidateToken(tokenString)
			if err != nil {
				log.Printf("AUTH 401: %s %s - Invalid token: %v", r.Method, r.URL.Path, err)
				handleUnauthorized(w, r, "Invalid or expired token", source)
				return
			}

			log.Printf("AUTH OK: %s %s - User: %s (%s)", r.Method, r.URL.Path, user.Name, user.Email)
			ctx := SetUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// handleUnauthorized returns 401 Unauthorized response
func handleUnauthorized(w http.ResponseWriter, r *http.Request, message string, source TokenSource) {
	switch source {
	case TokenSourceHeader:
		// API - return JSON with 401
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "unauthorized",
			"message": message,
			"status":  401,
		})

	case TokenSourceCookie:
		// Browser - redirect to login
		redirectURL := "/login"
		if r.URL.Path != "" && r.URL.Path != "/login" {
			redirectURL = "/login?redirect=" + r.URL.Path
		}
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)

	case TokenSourceAuto:
		if isAPIRequest(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "unauthorized",
				"message": message,
				"status":  401,
			})
		} else {
			redirectURL := "/login"
			if r.URL.Path != "" && r.URL.Path != "/login" {
				redirectURL = "/login?redirect=" + r.URL.Path
			}
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}
	}
}

// ==========================================
// PERMISSION MIDDLEWARE (403 errors)
// ==========================================

// RequirePermission checks for a single permission
func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := GetCurrentUser(r.Context())
			if err != nil {
				handleUnauthorized(w, r, "No authenticated user", TokenSourceAuto)
				return
			}

			if !hasPermission(user.Permissions, permission) {
				log.Printf("AUTHZ 403: %s %s - User %s lacks permission: %s",
					r.Method, r.URL.Path, user.Email, permission)
				handleForbidden(w, r, []string{permission}, "all")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermissionsMatchAll checks user has ALL permissions
func RequirePermissionsMatchAll(permissions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := GetCurrentUser(r.Context())
			if err != nil {
				handleUnauthorized(w, r, "No authenticated user", TokenSourceAuto)
				return
			}

			if !HasAllPermissions(user.Permissions, permissions) {
				log.Printf("AUTHZ 403: %s %s - User %s lacks all permissions: %v",
					r.Method, r.URL.Path, user.Email, permissions)
				handleForbidden(w, r, permissions, "all")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermissionsMatchAny checks user has ANY of the permissions
func RequirePermissionsMatchAny(permissions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := GetCurrentUser(r.Context())
			if err != nil {
				handleUnauthorized(w, r, "No authenticated user", TokenSourceAuto)
				return
			}

			if !HasAnyPermission(user.Permissions, permissions) {
				log.Printf("AUTHZ 403: %s %s - User %s lacks any permission: %v",
					r.Method, r.URL.Path, user.Email, permissions)
				handleForbidden(w, r, permissions, "any")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// handleForbidden returns 403 Forbidden response
func handleForbidden(w http.ResponseWriter, r *http.Request, requiredPermissions []string, match string) {
	// Check if HTMX request
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Retarget", "#error-message")
		w.Header().Set("HX-Reswap", "innerHTML")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`
			<div class="alert alert-error">
				<p>You don't have permission to perform this action.</p>
			</div>
		`))
		return
	}

	if isAPIRequest(r) {
		// API - return JSON with 403
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)

		var message string
		if match == "all" {
			message = "User requires all of the following permissions"
		} else {
			message = "User requires at least one of the following permissions"
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":                "forbidden",
			"message":              message,
			"required_permissions": requiredPermissions,
			"match":                match,
			"status":               403,
		})
	} else {
		// Browser - show error page
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Access Denied</title>
				<style>
					body { font-family: system-ui; max-width: 600px; margin: 100px auto; padding: 20px; }
					h1 { color: #d32f2f; }
					a { color: #1976d2; text-decoration: none; }
					a:hover { text-decoration: underline; }
				</style>
			</head>
			<body>
				<h1>403 - Access Denied</h1>
				<p>You don't have permission to access this resource.</p>
				<p><a href="/dashboard">← Go to Dashboard</a></p>
			</body>
			</html>
		`))
	}
}

// ==========================================
// HELPERS
// ==========================================

// isAPIRequest determines if the request is from an API client vs browser
func isAPIRequest(r *http.Request) bool {
	// Check if Authorization header present (API client)
	if r.Header.Get("Authorization") != "" {
		return true
	}

	// Check Accept header for JSON
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return true
	}

	// Check for HTMX request (browser)
	if r.Header.Get("HX-Request") != "" {
		return false
	}

	// Check if path starts with /api/
	if strings.HasPrefix(r.URL.Path, "/api/") {
		return true
	}

	// Default to browser
	return false
}
