package patientprescriptions

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

type PrescriptionListComponent struct {
	provider patientproviders.PatientPrescriptionProvider
	log      *zap.Logger
}

// Lazy loaded component
func NewPrescriptionListComponent(deps *contracts.UiDependencies) *PrescriptionListComponent {
	return &PrescriptionListComponent{provider: deps.PrescriptionProvider, log: deps.Log}
}

func (h *PrescriptionListComponent) Handler(w http.ResponseWriter, r *http.Request) {
	// Bind and validate query parameters
	req, _, err := bind.Query[request.PatientComponentRequest](r)
	if err != nil {
		h.log.Error("failed to bind query parameters", zap.Error(err))
		helper.WriteUIError(w, "Invalid patient ID parameter", http.StatusBadRequest)
		return
	}

	patientID := req.PatientID
	view, err := h.componentView(r.Context(), patientID)
	if err != nil {
		helper.WriteUIInternalError(w, "Failed to load patient prescriptions")
		return
	}
	if !helper.WaitOrContext(r.Context(), 3) {
		helper.WriteUIError(w, "Request canceled", http.StatusRequestTimeout)
		return
	}
	if err := view.Render(r.Context(), w); err != nil {
		helper.WriteUIInternalError(w, "Failed to render patient prescriptions")
	}
}

func (h *PrescriptionListComponent) componentView(ctx context.Context, patientID string) (templ.Component, error) {
	if patientID == "" {
		return nil, errors.New("patient id is required")
	}
	if h.provider == nil {
		return nil, errors.New("prescription provider is missing")
	}

	prescriptions, err := h.provider.PatientPrescriptionListByPatientID(ctx, patientID)
	if err != nil {
		if h.log != nil {
			h.log.Error("failed to load patient prescriptions", zap.Error(err), zap.String("patient_id", patientID))
		}
		return nil, err
	}

	params := PrescriptionListParams{
		Title:         "Prescriptions",
		EmptyMessage:  "No prescriptions found for this patient.",
		Prescriptions: prescriptions,
	}

	return PrescriptionListComponentView(params), nil
}
