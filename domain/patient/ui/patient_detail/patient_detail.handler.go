package patient_detail

import (
	"net/http"

	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	addressListcomponents "pharmacy-modernization-project-model/domain/patient/ui/components/addresslist_server_side"
	tools "pharmacy-modernization-project-model/internal/helper"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type PatientDetailHandler struct {
	patientsService    patSvc.PatientService
	addressListHandler *addressListcomponents.AddressListComponentHandler
	log                *zap.Logger
}

func NewPatientDetailHandler(patients patSvc.PatientService, addressListHandler *addressListcomponents.AddressListComponentHandler, log *zap.Logger) *PatientDetailHandler {
	return &PatientDetailHandler{
		patientsService:    patients,
		addressListHandler: addressListHandler,
		log:                log,
	}
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

	addressComponent, err := h.addressListHandler.Handler(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to load patient addresses", http.StatusInternalServerError)
		return
	}

	page := PatientDetailPage(PatientDetailPageParam{
		Patient:              patient,
		Age:                  tools.CalculateAge(patient.DOB),
		AddressListComponent: addressComponent,
	})

	if err := page.Render(r.Context(), w); err != nil {
		h.log.Error("failed to render patient detail", zap.Error(err), zap.String("id", id))
		http.Error(w, "failed to render patient detail", http.StatusInternalServerError)
		return
	}
}
