package patientprescriptions

import (
	"context"
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"go.uber.org/zap"

	patientproviders "pharmacy-modernization-project-model/domain/patient/providers"
)

type PrescriptionListDependencies struct {
	Provider patientproviders.PatientPrescriptionProvider
	Log      *zap.Logger
}

type PrescriptionListComponentHandler struct {
	provider patientproviders.PatientPrescriptionProvider
	log      *zap.Logger
}

func NewPrescriptionListComponentHandler(deps PrescriptionListDependencies) *PrescriptionListComponentHandler {
	return &PrescriptionListComponentHandler{provider: deps.Provider, log: deps.Log}
}

func (h *PrescriptionListComponentHandler) component(ctx context.Context, patientID string) (templ.Component, error) {
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

	return PrescriptionListComponent(params), nil
}

func (h *PrescriptionListComponentHandler) Handler(w http.ResponseWriter, r *http.Request) {
	patientID := r.URL.Query().Get("patientId")
	component, err := h.component(r.Context(), patientID)
	if err != nil {
		http.Error(w, "failed to load patient prescriptions", http.StatusInternalServerError)
		return
	}
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render patient prescriptions", http.StatusInternalServerError)
	}
}
