package iris_pharmacy

import (
	"context"
	"fmt"
	"strings"

	"pharmacy-modernization-project-model/internal/platform/httpclient"

	"go.uber.org/zap"
)

// HTTPClient implements PharmacyClient using HTTP requests
type HTTPClient struct {
	client    *httpclient.Client
	endpoints EndpointsConfig
	logger    *zap.Logger
}

// NewHTTPClient creates a new HTTP-based pharmacy client
func NewHTTPClient(cfg Config, client *httpclient.Client, logger *zap.Logger) *HTTPClient {
	return &HTTPClient{
		client:    client,
		endpoints: &cfg,
		logger:    logger,
	}
}

// replacePathParams replaces {paramName} in URL with actual values
func replacePathParams(url string, params map[string]string) string {
	for key, value := range params {
		placeholder := "{" + key + "}"
		url = strings.ReplaceAll(url, placeholder, value)
	}
	return url
}

// GetPrescription retrieves a prescription for a given prescription ID
func (c *HTTPClient) GetPrescription(ctx context.Context, prescriptionID string) (*PrescriptionResponse, error) {
	url := replacePathParams(c.endpoints.GetPrescriptionEndpoint(), map[string]string{
		"prescriptionID": prescriptionID,
	})

	c.logger.Debug("fetching prescription",
		zap.String("prescription_id", prescriptionID),
		zap.String("url", url),
	)

	var response PrescriptionResponse
	err := c.client.GetJSON(ctx, url, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get prescription: %w", err)
	}

	c.logger.Debug("prescription retrieved successfully",
		zap.String("prescription_id", prescriptionID),
		zap.String("drug", response.Drug),
		zap.String("status", response.Status),
		zap.String("pharmacy_name", response.PharmacyName),
	)

	return &response, nil
}

// Verify HTTPClient implements PharmacyClient
var _ PharmacyClient = (*HTTPClient)(nil)
