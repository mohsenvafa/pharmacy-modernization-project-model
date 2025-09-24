package iris_billing

import (
    "context"

    "go.uber.org/zap"
)

type MockClient struct {
	data map[string]GetInvoiceResponse
	log  *zap.Logger
}

func NewMockClient(seed map[string]GetInvoiceResponse, log *zap.Logger) *MockClient {
	data := make(map[string]GetInvoiceResponse)
	for k, v := range seed {
		data[k] = v
	}
	return &MockClient{data: data, log: log}
}

func (c *MockClient) GetInvoice(ctx context.Context, prescriptionID string) (GetInvoiceResponse, error) {
	if val, ok := c.data[prescriptionID]; ok {
		return val, nil
	}
	if c.log != nil {
		c.log.Warn("mock iris billing: invoice not found", zap.String("prescriptionID", prescriptionID))
	}
	return GetInvoiceResponse{PrescriptionID: prescriptionID, Status: "unbilled"}, nil
}

var _ Client = (*MockClient)(nil)
