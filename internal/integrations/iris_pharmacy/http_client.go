package iris_pharmacy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type HTTPClient struct {
	endpoint   string
	httpClient *http.Client
	log        *zap.Logger
}

func NewHTTPClient(cfg Config, httpClient *http.Client, log *zap.Logger) *HTTPClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	endpoint := strings.TrimSuffix(cfg.BaseURL, "/") + "/" + strings.Trim(cfg.Path, "/") + "/"
	return &HTTPClient{endpoint: endpoint, httpClient: httpClient, log: log}
}

func (c *HTTPClient) GetPrescription(ctx context.Context, prescriptionID string) (GetPrescriptionResponse, error) {
	var result GetPrescriptionResponse
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+prescriptionID, nil)
	if err != nil {
		if c.log != nil {
			c.log.Error("iris_pharmacy request", zap.Error(err), zap.String("prescriptionID", prescriptionID))
		}
		return result, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.log != nil {
			c.log.Error("iris_pharmacy response", zap.Error(err), zap.String("prescriptionID", prescriptionID))
		}
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("iris pharmacy: unexpected status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

var _ Client = (*HTTPClient)(nil)
