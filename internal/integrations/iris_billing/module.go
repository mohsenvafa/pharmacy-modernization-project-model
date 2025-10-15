package iris_billing

import (
	"time"

	"pharmacy-modernization-project-model/internal/platform/httpclient"

	"go.uber.org/zap"
)

// ModuleDependencies holds all dependencies required to initialize the billing module
type ModuleDependencies struct {
	Config     Config
	Logger     *zap.Logger
	HTTPClient *httpclient.Client
	UseMock    bool
	Timeout    time.Duration
}

// ModuleExport contains the exported services from the billing module
type ModuleExport struct {
	BillingClient BillingClient
}

// Module initializes and returns the billing module with its dependencies
func Module(deps ModuleDependencies) ModuleExport {
	// Use mock client if configured
	if deps.UseMock {
		deps.Logger.Info("initializing mock billing client")
		return ModuleExport{
			BillingClient: NewMockClient(deps.Logger),
		}
	}

	// Create HTTP client if not provided (fallback for tests/edge cases)
	if deps.HTTPClient == nil {
		timeout := deps.Timeout
		if timeout <= 0 {
			timeout = 30 * time.Second
		}

		deps.Logger.Warn("no shared http client provided, creating dedicated client for billing service",
			zap.Duration("timeout", timeout),
			zap.String("note", "consider passing shared client for better connection pooling"),
		)

		deps.HTTPClient = httpclient.NewClient(
			httpclient.Config{
				Timeout:     timeout,
				ServiceName: "iris_billing",
			},
			deps.Logger,
		)
	}

	deps.Logger.Info("initializing HTTP billing client",
		zap.String("get_invoice_url", deps.Config.GetInvoiceURL),
		zap.String("create_invoice_url", deps.Config.CreateInvoiceURL),
	)

	client := NewHTTPClient(deps.Config, deps.HTTPClient, deps.Logger)
	return ModuleExport{BillingClient: client}
}
