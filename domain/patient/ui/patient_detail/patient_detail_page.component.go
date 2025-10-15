package patient_detail

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	addresscomponents "pharmacy-modernization-project-model/domain/patient/ui/components/address_list"
	patientprescriptions "pharmacy-modernization-project-model/domain/patient/ui/components/patient_prescriptions"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/domain/patient/ui/paths"
	tools "pharmacy-modernization-project-model/internal/helper"
)

type PatientDetailComponent struct {
	patientsService           patSvc.PatientService
	addressListComponent      *addresscomponents.AddressListComponent
	prescriptionListComponent *patientprescriptions.PrescriptionListComponent
	log                       *zap.Logger
}

func NewPatientDetailComponent(
	deps *contracts.UiDependencies,
	addressListComponent *addresscomponents.AddressListComponent,
	prescriptionListComponent *patientprescriptions.PrescriptionListComponent,
) *PatientDetailComponent {
	return &PatientDetailComponent{
		patientsService:           deps.PatientSvc,
		addressListComponent:      addressListComponent,
		prescriptionListComponent: prescriptionListComponent,
		log:                       deps.Log,
	}
}

func (h *PatientDetailComponent) Handler(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	if patientID == "" {
		http.Error(w, "missing patient id", http.StatusBadRequest)
		return
	}

	patient, err := h.patientsService.GetByID(r.Context(), patientID)
	if err != nil {
		h.log.Error("failed to fetch patient", zap.Error(err), zap.String("id", patientID))
		http.Error(w, "failed to load patient", http.StatusInternalServerError)
		return
	}
	if patient.ID == "" {
		http.NotFound(w, r)
		return
	}

	var addressComponent, prescriptionComponent templ.Component

	if h.addressListComponent != nil {
		component, err := h.addressListComponent.View(r.Context(), patientID)
		if err != nil {
			http.Error(w, "failed to load patient addresses", http.StatusInternalServerError)
			return
		}
		addressComponent = component
	}

	if h.prescriptionListComponent != nil {
		prescriptionComponent = patientprescriptions.PlaceHolder(patientID)
	}

	view := PatientDetailPageComponentView(r.Context(), PatientDetailPageParam{
		Patient:       patient,
		Age:           tools.CalculateAge(patient.DOB),
		AddressList:   addressComponent,
		Prescriptions: prescriptionComponent,
		BackPath:      paths.PatientListURL(),
		EditPath:      paths.PatientEditURL(patientID),
	})

	if err := view.Render(r.Context(), w); err != nil {
		h.log.Error("failed to render patient detail", zap.Error(err), zap.String("id", patientID))
		http.Error(w, "failed to render patient detail", http.StatusInternalServerError)
		return
	}
}
