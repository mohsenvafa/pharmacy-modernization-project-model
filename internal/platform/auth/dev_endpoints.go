package auth

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// RegisterDevEndpoints registers development mode endpoints on the router
// Only registers endpoints when dev mode is enabled
func RegisterDevEndpoints(r chi.Router, logger *zap.Logger) {
	if !devModeEnabled {
		// Dev mode not enabled, skip endpoint registration
		return
	}

	// Register dev mode endpoints
	r.Get("/__dev/auth", DevAuthInfo)
	r.Get("/__dev/switch", SetMockUserCookie)

	logger.Info("Dev mode endpoints registered",
		zap.String("auth_path", "/__dev/auth"),
		zap.String("switch_path", "/__dev/switch"))
}
