package patient_list

import (
	"net/http"

	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	"go.uber.org/zap"
)

type PatientListHandler struct {
	patientsService patSvc.Service
	log             *zap.Logger
}

func NewPatientListHandler(patients patSvc.Service, log *zap.Logger) *PatientListHandler {
	return &PatientListHandler{patientsService: patients, log: log}
}

func (u *PatientListHandler) Handler(w http.ResponseWriter, r *http.Request) {
	pats, err := u.patientsService.List(r.Context(), "", 1000, 0)
	if err != nil {
		http.Error(w, "failed to load patients", http.StatusInternalServerError)
		return
	}

	page := PatientListPage(PatientListPageParam{
		NumberOfPatients: len(pats),
	})
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render patient list", http.StatusInternalServerError)
		return
	}
}
