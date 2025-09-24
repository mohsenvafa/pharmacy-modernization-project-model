package iris_pharmacy

import (
	"context"

	"go.uber.org/zap"
)

type MockClient struct {
	data map[string]GetPrescriptionResponse
	log  *zap.Logger
}

func NewMockClient(seed map[string]GetPrescriptionResponse, log *zap.Logger) *MockClient {
	data := make(map[string]GetPrescriptionResponse)
	for k, v := range seed {
		data[k] = v
	}
	return &MockClient{data: data, log: log}
}

func (c *MockClient) GetPrescription(ctx context.Context, prescriptionID string) (GetPrescriptionResponse, error) {
	if val, ok := c.data[prescriptionID]; ok {
		return val, nil
	}
	if c.log != nil {
		c.log.Warn("mock iris pharmacy: prescription not found", zap.String("prescriptionID", prescriptionID))
	}
	return GetPrescriptionResponse{ID: prescriptionID, Status: "unknown", PharmacyName: "CVS", PharmacyType: "Retail"}, nil
}

var _ Client = (*MockClient)(nil)
