package integrations

import (
	"time"

	"go.uber.org/zap"

	irisbilling "pharmacy-modernization-project-model/internal/integrations/iris_billing"
	irispharmacy "pharmacy-modernization-project-model/internal/integrations/iris_pharmacy"
	"pharmacy-modernization-project-model/internal/platform/config"
	"pharmacy-modernization-project-model/internal/platform/httpclient"
	"pharmacy-modernization-project-model/internal/platform/httpclient/interceptors"
)

// Dependencies holds all required dependencies for the integrations layer
type Dependencies struct {
	Config *config.Config
	Logger *zap.Logger
}

// Export contains all integration services exported by this package
type Export struct {
	PharmacyClient irispharmacy.PharmacyClient
	BillingClient  irisbilling.BillingClient
}

// New initializes all integration services with their dependencies
func New(deps Dependencies) Export {
	if deps.Config == nil {
		deps.Logger.Warn("config is nil, returning empty integrations export")
		return Export{}
	}

	logger := deps.Logger.With(zap.String("layer", "integrations"))

	// Create metrics interceptor to track API call timing
	metricsInterceptor := interceptors.NewMetricsInterceptor(logger)

	// Create global header provider for all API requests
	// These headers will be added to ALL requests across all integrations
	globalHeaderProvider := httpclient.NewStaticHeaderProvider(map[string]string{
		"X-IRIS-User-ID": "xyz", // ✅ Example: Add user ID to all API calls
		// Add more global headers here as needed:
		// "X-Client-Version": "1.0.0",
		// "X-Request-Source": "rxintake-app",
	})

	// Create shared HTTP client for all external API integrations
	// This client is reused across all integration services for efficient connection pooling
	sharedHTTPClient := httpclient.NewClient(
		httpclient.Config{
			Timeout:        30 * time.Second,     // Default timeout for all external APIs
			MaxIdleConns:   100,                  // Connection pool size
			ServiceName:    "external_apis",      // For observability/logging
			HeaderProvider: globalHeaderProvider, // ✅ Global headers for ALL requests
			// For auth tokens, see integration_wire_with_auth_example.go (Stargate example)
		},
		logger,
		metricsInterceptor, // ✅ Track timing for all API calls
	)

	logger.Info("shared http client created with global headers",
		zap.String("X-IRIS-User-ID", "xyz"),
	)

	// Initialize pharmacy client
	pharmacy := irispharmacy.Module(irispharmacy.ModuleDependencies{
		Config: irispharmacy.Config{
			GetPrescriptionURL: deps.Config.External.Pharmacy.Endpoints.GetPrescription,
		},
		Logger:     logger.With(zap.String("service", "pharmacy")),
		HTTPClient: sharedHTTPClient, // Use the shared client
		UseMock:    deps.Config.External.Pharmacy.UseMock,
		Timeout:    parseDuration(deps.Config.External.Pharmacy.Timeout, 30*time.Second),
	}).PharmacyClient

	// Initialize billing client
	billing := irisbilling.Module(irisbilling.ModuleDependencies{
		Config: irisbilling.Config{
			GetInvoiceURL:           deps.Config.External.Billing.Endpoints.GetInvoice,
			GetInvoicesByPatientURL: deps.Config.External.Billing.Endpoints.GetInvoicesByPatient,
			CreateInvoiceURL:        deps.Config.External.Billing.Endpoints.CreateInvoice,
			AcknowledgeInvoiceURL:   deps.Config.External.Billing.Endpoints.AcknowledgeInvoice,
			GetInvoicePaymentURL:    deps.Config.External.Billing.Endpoints.GetInvoicePayment,
		},
		Logger:     logger.With(zap.String("service", "billing")),
		HTTPClient: sharedHTTPClient, // Use the shared client
		UseMock:    deps.Config.External.Billing.UseMock,
		Timeout:    parseDuration(deps.Config.External.Billing.Timeout, 30*time.Second),
	}).BillingClient

	logger.Info("integrations layer initialized successfully")

	return Export{
		PharmacyClient: pharmacy,
		BillingClient:  billing,
	}
}

// parseDuration safely parses a duration string with a fallback
func parseDuration(value string, fallback time.Duration) time.Duration {
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	return fallback
}
