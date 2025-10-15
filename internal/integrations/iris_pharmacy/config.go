package iris_pharmacy

// Config holds the configuration for the IRIS pharmacy service
type Config struct {
	// API Endpoints (full URLs from YAML config)
	GetPrescriptionURL string
}

// EndpointsConfig defines the interface for pharmacy endpoints configuration
type EndpointsConfig interface {
	GetPrescriptionEndpoint() string
}

// Verify Config implements EndpointsConfig
var _ EndpointsConfig = (*Config)(nil)

// GetPrescriptionEndpoint returns the full URL for getting a prescription
func (c *Config) GetPrescriptionEndpoint() string {
	return c.GetPrescriptionURL
}
