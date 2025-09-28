package patientprescriptions

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"go.uber.org/zap"

	patientproviders "pharmacy-modernization-project-model/domain/patient/providers"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
)

type PrescriptionListComponent struct {
	provider patientproviders.PatientPrescriptionProvider
	log      *zap.Logger
}

func NewPrescriptionListComponent(deps *contracts.UiDependencies) *PrescriptionListComponent {
	return &PrescriptionListComponent{provider: deps.PrescriptionProvider, log: deps.Log}
}

func (h *PrescriptionListComponent) Handler(w http.ResponseWriter, r *http.Request) {
	patientID := r.URL.Query().Get("patientId")
	view, err := h.componentView(r.Context(), patientID)
	if err != nil {
		http.Error(w, "failed to load patient prescriptions", http.StatusInternalServerError)
		return
	}
	select {
	case <-time.After(3 * time.Second):
	case <-r.Context().Done():
		http.Error(w, "request canceled", http.StatusRequestTimeout)
		return
	}
	if err := view.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render patient prescriptions", http.StatusInternalServerError)
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
