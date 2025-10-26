package patientinvoices

import (
	"context"
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/domain/patient/contracts/request"
	patientproviders "pharmacy-modernization-project-model/domain/patient/providers"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/internal/bind"
	helper "pharmacy-modernization-project-model/internal/helper"
)

type InvoiceListComponent struct {
	provider patientproviders.PatientInvoiceProvider
	log      *zap.Logger
}

// Lazy loaded component
func NewInvoiceListComponent(deps *contracts.UiDependencies) *InvoiceListComponent {
	return &InvoiceListComponent{
		provider: deps.InvoiceProvider,
		log:      deps.Log,
	}
}

func (h *InvoiceListComponent) Handler(w http.ResponseWriter, r *http.Request) {
	// Bind and validate query parameters
	req, _, err := bind.Query[request.PatientComponentRequest](r)
	if err != nil {
		h.log.Error("failed to bind query parameters", zap.Error(err))
		helper.WriteUIError(w, "Invalid patient ID parameter", http.StatusBadRequest)
		return
	}

	patientID := req.PatientID

	// Wait for 3 seconds (simulating loading)
	if !helper.WaitOrContext(r.Context(), 3) {
		// Request was canceled - render error view
		errorView := ErrorView(patientID, "Request was canceled or timed out.")
		if err := errorView.Render(r.Context(), w); err != nil {
			http.Error(w, "failed to render error view", http.StatusInternalServerError)
		}
		return
	}

	view, err := h.componentView(r.Context(), patientID)
	if err != nil {
		// Error loading data - render error view instead of HTTP error
		errorView := ErrorView(patientID, "Unable to load invoice data. Please try again.")
		if renderErr := errorView.Render(r.Context(), w); renderErr != nil {
			http.Error(w, "failed to render error view", http.StatusInternalServerError)
		}
		return
	}

	if err := view.Render(r.Context(), w); err != nil {
		// Error rendering view - render error view
		errorView := ErrorView(patientID, "Unable to display invoice data. Please try again.")
		if renderErr := errorView.Render(r.Context(), w); renderErr != nil {
			http.Error(w, "failed to render error view", http.StatusInternalServerError)
		}
	}
}

func (h *InvoiceListComponent) componentView(ctx context.Context, patientID string) (templ.Component, error) {
	if patientID == "" {
		return nil, errors.New("patient id is required")
	}
	if h.provider == nil {
		return nil, errors.New("invoice provider is missing")
	}

	invoiceResponse, err := h.provider.GetInvoicesByPatientID(ctx, patientID)
	if err != nil {
		if h.log != nil {
			h.log.Error("failed to load patient invoices", zap.Error(err))
		}
		return nil, err
	}

	params := InvoiceListParams{
		Title:        "Invoices",
		EmptyMessage: "No invoices found for this patient.",
		PatientID:    patientID,
		Invoices:     invoiceResponse.Invoices,
		Total:        invoiceResponse.Total,
	}

	return InvoiceListComponentView(params), nil
}
