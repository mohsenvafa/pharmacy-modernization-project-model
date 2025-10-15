package stargate

// Config holds the configuration for the Stargate authentication service
type Config struct {
	// OAuth/Token Endpoints (full URLs from YAML config)
	TokenURL        string
	RefreshTokenURL string

	// Client credentials
	ClientID     string
	ClientSecret string
	Scope        string
}

// EndpointsConfig defines the interface for Stargate endpoints configuration
type EndpointsConfig interface {
	TokenEndpoint() string
	RefreshTokenEndpoint() string
}

// Verify Config implements EndpointsConfig
var _ EndpointsConfig = (*Config)(nil)

// TokenEndpoint returns the full URL for getting a token
func (c *Config) TokenEndpoint() string {
	return c.TokenURL
}

// RefreshTokenEndpoint returns the full URL for refreshing a token
func (c *Config) RefreshTokenEndpoint() string {
	return c.RefreshTokenURL
}
