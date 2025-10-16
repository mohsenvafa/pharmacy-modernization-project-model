package iris_billing

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"pharmacy-modernization-project-model/internal/platform/httpclient"

	"go.uber.org/zap"
)

// HTTPClient implements BillingClient using HTTP requests
type HTTPClient struct {
	client    *httpclient.Client
	endpoints EndpointsConfig
	logger    *zap.Logger
}

// NewHTTPClient creates a new HTTP-based billing client
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

// generateIdempotencyKey creates a deterministic idempotency key based on prescription ID
// This ensures the same invoice creation request always generates the same key
func generateIdempotencyKey(prescriptionID string) string {
	hash := sha256.Sum256([]byte("create-invoice-" + prescriptionID))
	return fmt.Sprintf("%x", hash[:16]) // Use first 16 bytes of hash
}

// GetInvoice retrieves an invoice for a given prescription ID
func (c *HTTPClient) GetInvoice(ctx context.Context, prescriptionID string) (*InvoiceResponse, error) {
	url := replacePathParams(c.endpoints.GetInvoiceEndpoint(), map[string]string{
		"prescriptionID": prescriptionID,
	})

	c.logger.Debug("fetching invoice",
		zap.String("prescription_id", prescriptionID),
		zap.String("url", url),
	)

	// ✅ Example: Add endpoint-specific header for THIS endpoint only
	resp, err := c.client.Get(ctx, url, map[string]string{
		"Content-Type":    "application/json",
		"Accept":          "application/json",
		"X-IRIS-Env-Name": "IRIS_stage", // ✅ Only for GetInvoice endpoint
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: request failed", resp.StatusCode)
	}

	var response InvoiceResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		c.logger.Error("failed to decode invoice response",
			zap.String("prescription_id", prescriptionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode invoice response: %w", err)
	}

	c.logger.Debug("invoice retrieved successfully",
		zap.String("prescription_id", prescriptionID),
		zap.String("invoice_id", response.ID),
		zap.Float64("amount", response.Amount),
		zap.String("status", response.Status),
	)

	return &response, nil
}

// GetInvoicesByPatientID retrieves all invoices for a given patient ID
func (c *HTTPClient) GetInvoicesByPatientID(ctx context.Context, patientID string) (*InvoiceListResponse, error) {
	url := replacePathParams(c.endpoints.GetInvoicesByPatientEndpoint(), map[string]string{
		"patientID": patientID,
	})

	c.logger.Debug("fetching invoices by patient ID",
		zap.String("patient_id", patientID),
		zap.String("url", url),
	)

	resp, err := c.client.Get(ctx, url, map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get invoices by patient ID: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: request failed", resp.StatusCode)
	}

	var response InvoiceListResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		c.logger.Error("failed to decode invoice list response",
			zap.String("patient_id", patientID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode invoice list response: %w", err)
	}

	c.logger.Debug("invoices retrieved successfully",
		zap.String("patient_id", patientID),
		zap.Int("count", response.Total),
	)

	return &response, nil
}

// CreateInvoice creates a new invoice
func (c *HTTPClient) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	url := c.endpoints.CreateInvoiceEndpoint()

	// Generate idempotency key to prevent duplicate invoice creation
	idempotencyKey := generateIdempotencyKey(req.PrescriptionID)

	c.logger.Debug("creating invoice",
		zap.String("prescription_id", req.PrescriptionID),
		zap.Float64("amount", req.Amount),
		zap.String("url", url),
		zap.String("idempotency_key", idempotencyKey),
	)

	// Prepare request body
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make POST request with idempotency key header
	resp, err := c.client.Post(ctx, url, bytes.NewReader(bodyBytes), map[string]string{
		"Content-Type":      "application/json",
		"Accept":            "application/json",
		"X-Idempotency-Key": idempotencyKey, // ✅ Prevents duplicate invoices
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: request failed", resp.StatusCode)
	}

	var response CreateInvoiceResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		c.logger.Error("failed to decode create invoice response",
			zap.String("prescription_id", req.PrescriptionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Info("invoice created successfully",
		zap.String("invoice_id", response.ID),
		zap.String("prescription_id", response.PrescriptionID),
		zap.Float64("amount", response.Amount),
	)

	return &response, nil
}

// AcknowledgeInvoice acknowledges an invoice
func (c *HTTPClient) AcknowledgeInvoice(ctx context.Context, invoiceID string, req AcknowledgeInvoiceRequest) (*AcknowledgeInvoiceResponse, error) {
	url := replacePathParams(c.endpoints.AcknowledgeInvoiceEndpoint(), map[string]string{
		"invoiceID": invoiceID,
	})

	c.logger.Debug("acknowledging invoice",
		zap.String("invoice_id", invoiceID),
		zap.String("acknowledged_by", req.AcknowledgedBy),
		zap.String("url", url),
	)

	var response AcknowledgeInvoiceResponse
	err := c.client.PostJSON(ctx, url, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to acknowledge invoice: %w", err)
	}

	c.logger.Info("invoice acknowledged successfully",
		zap.String("invoice_id", invoiceID),
	)

	return &response, nil
}

// GetInvoicePayment retrieves payment details for an invoice
func (c *HTTPClient) GetInvoicePayment(ctx context.Context, invoiceID string) (*InvoicePaymentResponse, error) {
	url := replacePathParams(c.endpoints.GetInvoicePaymentEndpoint(), map[string]string{
		"invoiceID": invoiceID,
	})

	c.logger.Debug("fetching invoice payment",
		zap.String("invoice_id", invoiceID),
		zap.String("url", url),
	)

	var response InvoicePaymentResponse
	err := c.client.GetJSON(ctx, url, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice payment: %w", err)
	}

	c.logger.Debug("invoice payment retrieved successfully",
		zap.String("invoice_id", invoiceID),
		zap.String("payment_id", response.PaymentID),
		zap.String("status", response.Status),
	)

	return &response, nil
}

// Verify HTTPClient implements BillingClient
var _ BillingClient = (*HTTPClient)(nil)
