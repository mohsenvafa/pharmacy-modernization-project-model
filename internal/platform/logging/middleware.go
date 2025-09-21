package logging

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ctxKey string
const correlationKey ctxKey = "correlation-id"

func CorrelationID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cid := r.Header.Get("X-Correlation-Id")
			if cid == "" { cid = uuid.New().String() }
			w.Header().Set("X-Correlation-Id", cid)
			r = r.WithContext(context.WithValue(r.Context(), correlationKey, cid))
			next.ServeHTTP(w, r)
		})
	}
}

func GetCorrelationID(ctx context.Context) string {
	if v, ok := ctx.Value(correlationKey).(string); ok { return v }
	return ""
}

func ZapRequestLogger(l *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			next.ServeHTTP(ww, r)
			l.Info("http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.Duration("duration", time.Since(start)),
				zap.String("request_id", middleware.GetReqID(r.Context())),
				zap.String("correlation_id", GetCorrelationID(r.Context())),
				zap.String("remote_ip", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}
