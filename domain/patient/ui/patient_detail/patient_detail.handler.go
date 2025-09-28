package patient_detail

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	addresscomponents "pharmacy-modernization-project-model/domain/patient/ui/components/addresslist_server_side"
	patientprescriptions "pharmacy-modernization-project-model/domain/patient/ui/components/patient_prescriptions"
	tools "pharmacy-modernization-project-model/internal/helper"
)

type PatientDetailHandler struct {
	patientsService         patSvc.PatientService
	addressListHandler      *addresscomponents.AddressListComponentHandler
	prescriptionListHandler *patientprescriptions.PrescriptionListComponentHandler
	log                     *zap.Logger
}

func NewPatientDetailHandler(
	patients patSvc.PatientService,
	addressListHandler *addresscomponents.AddressListComponentHandler,
	prescriptionListHandler *patientprescriptions.PrescriptionListComponentHandler,
	log *zap.Logger,
) *PatientDetailHandler {
	return &PatientDetailHandler{
		patientsService:         patients,
		addressListHandler:      addressListHandler,
		prescriptionListHandler: prescriptionListHandler,
		log:                     log,
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

	var addressComponent, prescriptionComponent templ.Component

	if h.addressListHandler != nil {
		component, err := h.addressListHandler.Handler(r.Context(), id)
		if err != nil {
			http.Error(w, "failed to load patient addresses", http.StatusInternalServerError)
			return
		}
		addressComponent = component
	}

	if h.prescriptionListHandler != nil {
		prescriptionComponent = patientprescriptions.PlaceHolder(id)
	}

	page := PatientDetailPage(PatientDetailPageParam{
		Patient:       patient,
		Age:           tools.CalculateAge(patient.DOB),
		AddressList:   addressComponent,
		Prescriptions: prescriptionComponent,
	})

	if err := page.Render(r.Context(), w); err != nil {
		h.log.Error("failed to render patient detail", zap.Error(err), zap.String("id", id))
		http.Error(w, "failed to render patient detail", http.StatusInternalServerError)
		return
	}
}
