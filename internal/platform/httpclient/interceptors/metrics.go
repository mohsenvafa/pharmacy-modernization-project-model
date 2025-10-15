package interceptors

import (
	"context"
	"net/http"

	"pharmacy-modernization-project-model/internal/platform/httpclient"

	"go.uber.org/zap"
)

// MetricsInterceptor collects metrics about HTTP requests
type MetricsInterceptor struct {
	logger *zap.Logger
}

// NewMetricsInterceptor creates a new metrics interceptor
func NewMetricsInterceptor(logger *zap.Logger) *MetricsInterceptor {
	return &MetricsInterceptor{
		logger: logger,
	}
}

func (m *MetricsInterceptor) Before(ctx context.Context, req *http.Request) error {
	return nil
}

func (m *MetricsInterceptor) After(ctx context.Context, resp *http.Response, response *httpclient.Response) error {
	// Log metrics
	m.logger.Info("http metrics",
		zap.String("method", resp.Request.Method),
		zap.String("url", resp.Request.URL.String()),
		zap.Int("status", resp.StatusCode),
		zap.Duration("duration", response.Duration),
		zap.Int("response_bytes", len(response.Body)),
	)

	// Here you could also send metrics to a metrics service like Prometheus
	// prometheus.HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	// prometheus.HTTPRequestTotal.WithLabelValues(method, path, status).Inc()

	return nil
}
