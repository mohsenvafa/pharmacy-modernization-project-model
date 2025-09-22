package patient_list

import (
	"net/http"

	patientmodel "github.com/pharmacy-modernization-project-model/internal/domain/patient/model"
	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	"go.uber.org/zap"
)

type PatientListHandler struct {
	patientsService patSvc.PatientService
	log             *zap.Logger
}

func NewPatientListHandler(patients patSvc.PatientService, log *zap.Logger) *PatientListHandler {
	return &PatientListHandler{patientsService: patients, log: log}
}

func firstNPatients(pats []patientmodel.Patient, n int) []patientmodel.Patient {
	if len(pats) <= n {
		return pats
	}
	return pats[:n]
}

func (u *PatientListHandler) Handler(w http.ResponseWriter, r *http.Request) {
	patients, err := u.patientsService.List(r.Context(), "", 1000, 0)
	if err != nil {
		http.Error(w, "failed to load patients", http.StatusInternalServerError)
		return
	}

	page := PatientListPage(PatientListPageParam{
		Patients: firstNPatients(patients, 5),
	})
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render patient list", http.StatusInternalServerError)
		return
	}
}
