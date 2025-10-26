package patient_detail

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/domain/patient/contracts/request"
	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	addresscomponents "pharmacy-modernization-project-model/domain/patient/ui/components/address_list"
	patientinvoices "pharmacy-modernization-project-model/domain/patient/ui/components/patient_invoices"
	patientprescriptions "pharmacy-modernization-project-model/domain/patient/ui/components/patient_prescriptions"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/domain/patient/ui/paths"
	"pharmacy-modernization-project-model/internal/bind"
	helper "pharmacy-modernization-project-model/internal/helper"
)

type PatientDetailComponent struct {
	patientsService           patSvc.PatientService
	addressListComponent      *addresscomponents.AddressListComponent
	prescriptionListComponent *patientprescriptions.PrescriptionListComponent
	invoiceListComponent      *patientinvoices.InvoiceListComponent
	log                       *zap.Logger
}

func NewPatientDetailComponent(
	deps *contracts.UiDependencies,
	addressListComponent *addresscomponents.AddressListComponent,
	prescriptionListComponent *patientprescriptions.PrescriptionListComponent,
	invoiceListComponent *patientinvoices.InvoiceListComponent,
) *PatientDetailComponent {
	return &PatientDetailComponent{
		patientsService:           deps.PatientSvc,
		addressListComponent:      addressListComponent,
		prescriptionListComponent: prescriptionListComponent,
		invoiceListComponent:      invoiceListComponent,
		log:                       deps.Log,
	}
}

func (h *PatientDetailComponent) Handler(w http.ResponseWriter, r *http.Request) {
	// Bind and validate path parameters
	pathVars, _, err := bind.ChiPath[request.PatientPathVars](r, chi.URLParam)
	if err != nil {
		h.log.Error("failed to bind path parameters", zap.Error(err))
		helper.WriteUIError(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	patient, err := h.patientsService.GetByID(r.Context(), pathVars.PatientID)
	if err != nil {
		h.log.Error("failed to fetch patient", zap.Error(err))
		helper.WriteUIInternalError(w, "Failed to load patient")
		return
	}
	if patient.ID == "" {
		helper.WriteUINotFound(w, "Patient not found")
		return
	}

	var addressComponent, prescriptionComponent, invoiceComponent templ.Component

	if h.addressListComponent != nil {
		component, err := h.addressListComponent.View(r.Context(), pathVars.PatientID)
		if err != nil {
			helper.WriteUIInternalError(w, "Failed to load patient addresses")
			return
		}
		addressComponent = component
	}

	if h.prescriptionListComponent != nil {
		prescriptionComponent = patientprescriptions.PlaceHolder(pathVars.PatientID)
	}

	if h.invoiceListComponent != nil {
		invoiceComponent = patientinvoices.PlaceHolder(pathVars.PatientID)
	}

	view := PatientDetailPageComponentView(r.Context(), PatientDetailPageParam{
		Patient:       patient,
		Age:           helper.CalculateAge(patient.DOB),
		AddressList:   addressComponent,
		Prescriptions: prescriptionComponent,
		Invoices:      invoiceComponent,
		BackPath:      paths.PatientListURL(),
		EditPath:      paths.PatientEditURL(pathVars.PatientID),
	})

	if err := view.Render(r.Context(), w); err != nil {
		h.log.Error("failed to render patient detail", zap.Error(err))
		helper.WriteUIInternalError(w, "Failed to render patient detail")
		return
	}
}
