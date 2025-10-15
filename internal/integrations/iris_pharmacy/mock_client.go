package iris_pharmacy

import (
	"context"

	"go.uber.org/zap"
)

// MockClient implements PharmacyClient with in-memory mock data
type MockClient struct {
	data   map[string]PrescriptionResponse
	logger *zap.Logger
}

// NewMockClient creates a new mock pharmacy client
func NewMockClient(logger *zap.Logger) *MockClient {
	return &MockClient{
		data:   make(map[string]PrescriptionResponse),
		logger: logger,
	}
}

// SeedPrescription adds a mock prescription (useful for testing)
func (c *MockClient) SeedPrescription(prescription PrescriptionResponse) {
	c.data[prescription.ID] = prescription
}

// GetPrescription retrieves a mock prescription for a given prescription ID
func (c *MockClient) GetPrescription(ctx context.Context, prescriptionID string) (*PrescriptionResponse, error) {
	if prescription, ok := c.data[prescriptionID]; ok {
		c.logger.Debug("mock prescription found",
			zap.String("prescription_id", prescriptionID),
			zap.String("drug", prescription.Drug),
			zap.String("status", prescription.Status),
		)
		return &prescription, nil
	}

	c.logger.Warn("mock prescription not found, returning default",
		zap.String("prescription_id", prescriptionID),
	)

	// Return a default prescription for unknown prescription IDs
	defaultPrescription := &PrescriptionResponse{
		ID:           prescriptionID,
		Status:       "unknown",
		PharmacyName: "CVS",
		PharmacyType: "Retail",
	}

	return defaultPrescription, nil
}

// Verify MockClient implements PharmacyClient
var _ PharmacyClient = (*MockClient)(nil)
