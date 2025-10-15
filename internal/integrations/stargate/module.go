package stargate

import (
	"time"

	"pharmacy-modernization-project-model/internal/platform/httpclient"

	"go.uber.org/zap"
)

// ModuleDependencies holds all dependencies required to initialize the Stargate token module
type ModuleDependencies struct {
	Config     Config
	Logger     *zap.Logger
	HTTPClient *httpclient.Client
	UseMock    bool
	Timeout    time.Duration
}

// ModuleExport contains the exported services from the Stargate token module
type ModuleExport struct {
	TokenClient TokenClient
}

// Module initializes and returns the Stargate token module with its dependencies
func Module(deps ModuleDependencies) ModuleExport {
	// Use mock client if configured
	if deps.UseMock {
		deps.Logger.Info("initializing mock Stargate token client")
		return ModuleExport{
			TokenClient: NewMockClient(deps.Logger),
		}
	}

	// Create HTTP client if not provided (fallback for tests/edge cases)
	if deps.HTTPClient == nil {
		timeout := deps.Timeout
		if timeout <= 0 {
			timeout = 30 * time.Second
		}

		deps.Logger.Warn("no shared http client provided, creating dedicated client for Stargate",
			zap.Duration("timeout", timeout),
			zap.String("note", "consider passing shared client for better connection pooling"),
		)

		deps.HTTPClient = httpclient.NewClient(
			httpclient.Config{
				Timeout:     timeout,
				ServiceName: "stargate_auth",
			},
			deps.Logger,
		)
	}

	deps.Logger.Info("initializing HTTP Stargate token client",
		zap.String("token_url", deps.Config.TokenURL),
	)

	client := NewHTTPClient(deps.Config, deps.HTTPClient, deps.Logger)
	return ModuleExport{TokenClient: client}
}
