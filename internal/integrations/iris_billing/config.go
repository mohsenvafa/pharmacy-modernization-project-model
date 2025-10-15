package iris_billing

// Config holds the configuration for the IRIS billing service
type Config struct {
	// API Endpoints (full URLs from YAML config)
	GetInvoiceURL         string
	CreateInvoiceURL      string
	AcknowledgeInvoiceURL string
	GetInvoicePaymentURL  string
}

// EndpointsConfig defines the interface for billing endpoints configuration
type EndpointsConfig interface {
	GetInvoiceEndpoint() string
	CreateInvoiceEndpoint() string
	AcknowledgeInvoiceEndpoint() string
	GetInvoicePaymentEndpoint() string
}

// Verify Config implements EndpointsConfig
var _ EndpointsConfig = (*Config)(nil)

// GetInvoiceEndpoint returns the full URL for getting an invoice
func (c *Config) GetInvoiceEndpoint() string {
	return c.GetInvoiceURL
}

// CreateInvoiceEndpoint returns the full URL for creating an invoice
func (c *Config) CreateInvoiceEndpoint() string {
	return c.CreateInvoiceURL
}

// AcknowledgeInvoiceEndpoint returns the full URL for acknowledging an invoice
func (c *Config) AcknowledgeInvoiceEndpoint() string {
	return c.AcknowledgeInvoiceURL
}

// GetInvoicePaymentEndpoint returns the full URL for getting invoice payment
func (c *Config) GetInvoicePaymentEndpoint() string {
	return c.GetInvoicePaymentURL
}
