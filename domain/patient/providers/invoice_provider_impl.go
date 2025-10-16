package providers

import (
	"context"

	irisbilling "pharmacy-modernization-project-model/internal/integrations/iris_billing"

	"go.uber.org/zap"
)

type InvoiceProviderImpl struct {
	billingClient irisbilling.BillingClient
	logger        *zap.Logger
}

func NewInvoiceProvider(billingClient irisbilling.BillingClient, logger *zap.Logger) PatientInvoiceProvider {
	return &InvoiceProviderImpl{
		billingClient: billingClient,
		logger:        logger,
	}
}

func (p *InvoiceProviderImpl) GetInvoicesByPatientID(ctx context.Context, patientID string) (*irisbilling.InvoiceListResponse, error) {
	p.logger.Debug("fetching invoices for patient",
		zap.String("patient_id", patientID),
	)

	response, err := p.billingClient.GetInvoicesByPatientID(ctx, patientID)
	if err != nil {
		p.logger.Error("failed to fetch invoices",
			zap.String("patient_id", patientID),
			zap.Error(err),
		)
		return nil, err
	}

	return response, nil
}
