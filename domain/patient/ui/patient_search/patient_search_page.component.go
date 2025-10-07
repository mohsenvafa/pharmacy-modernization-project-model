package patient_search

import (
	"net/http"

	"go.uber.org/zap"

	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
)

type PatientListComponent struct {
	log *zap.Logger
}

func NewPatientPageComponent(deps *contracts.UiDependencies) *PatientListComponent {
	return &PatientListComponent{log: deps.Log}
}

func (c *PatientListComponent) Handler(w http.ResponseWriter, r *http.Request) {

	view := PatientSearchPageComponentView()

	if err := view.Render(r.Context(), w); err != nil {
		if c.log != nil {
			c.log.Error("failed to render patient list", zap.Error(err))
		}
		http.Error(w, "failed to render patient list", http.StatusInternalServerError)
		return
	}
}
