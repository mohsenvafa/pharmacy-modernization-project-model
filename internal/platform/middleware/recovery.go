package middleware

import (
	"net/http"

	"pharmacy-modernization-project-model/internal/platform/sanitizer"

	"go.uber.org/zap"
)

// RecoveryMiddleware provides panic recovery for HTTP handlers
func RecoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						zap.Any("error", err),
						zap.String("method", sanitizer.ForLogging(r.Method)),
						zap.String("url", sanitizer.ForLogging(r.URL.String())),
						zap.String("remote_addr", sanitizer.ForLogging(r.RemoteAddr)))

					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
