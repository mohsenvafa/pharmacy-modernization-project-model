package patient_search

import (
	"net/http"

	"go.uber.org/zap"

	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
)

type PatientSearchComponent struct {
	log *zap.Logger
}

func NewPatientSearchPageComponent(deps *contracts.UiDependencies) *PatientSearchComponent {
	return &PatientSearchComponent{log: deps.Log}
}

func (c *PatientSearchComponent) Handler(w http.ResponseWriter, r *http.Request) {
	// Set cache-busting headers to prevent caching issues
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	view := PatientSearchPageComponentView()

	if err := view.Render(r.Context(), w); err != nil {
		if c.log != nil {
			c.log.Error("failed to render patient search", zap.Error(err))
		}
		http.Error(w, "failed to render patient search", http.StatusInternalServerError)
		return
	}
}
