package patient_detail

import (
	"net/http"
	"time"

	patSvc "pharmacy-modernization-project-model/domain/patient/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type PatientDetailHandler struct {
	patientsService patSvc.PatientService
	log             *zap.Logger
}

func NewPatientDetailHandler(patients patSvc.PatientService, log *zap.Logger) *PatientDetailHandler {
	return &PatientDetailHandler{patientsService: patients, log: log}
}

func (h *PatientDetailHandler) Handler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "patientID")
	if id == "" {
		http.Error(w, "missing patient id", http.StatusBadRequest)
		return
	}

	patient, err := h.patientsService.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("failed to fetch patient", zap.Error(err), zap.String("id", id))
		http.Error(w, "failed to load patient", http.StatusInternalServerError)
		return
	}
	if patient.ID == "" {
		http.NotFound(w, r)
		return
	}

	page := PatientDetailPage(PatientDetailPageParam{
		Patient: patient,
		Age:     calculateAge(patient.DOB),
	})

	if err := page.Render(r.Context(), w); err != nil {
		h.log.Error("failed to render patient detail", zap.Error(err), zap.String("id", id))
		http.Error(w, "failed to render patient detail", http.StatusInternalServerError)
		return
	}
}

func calculateAge(dob time.Time) int {
	if dob.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - dob.Year()
	anniversary := time.Date(now.Year(), dob.Month(), dob.Day(), dob.Hour(), dob.Minute(), dob.Second(), dob.Nanosecond(), dob.Location())
	if now.Before(anniversary) {
		age--
	}
	if age < 0 {
		age = 0
	}
	return age
}
