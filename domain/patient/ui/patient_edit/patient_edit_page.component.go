package patient_edit

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/domain/patient/contracts/request"
	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/domain/patient/ui/contracts/form_data"
	"pharmacy-modernization-project-model/domain/patient/ui/paths"
	"pharmacy-modernization-project-model/internal/bind"
	helper "pharmacy-modernization-project-model/internal/helper"
)

type PatientEditComponent struct {
	patientsService patSvc.PatientService
	log             *zap.Logger
}

func NewPatientEditComponent(deps *contracts.UiDependencies) *PatientEditComponent {
	return &PatientEditComponent{
		patientsService: deps.PatientSvc,
		log:             deps.Log,
	}
}

// ShowEditForm handles GET requests to show the edit form
func (h *PatientEditComponent) ShowEditForm(w http.ResponseWriter, r *http.Request) {
	// Bind and validate path parameters
	pathVars, _, err := bind.ChiPath[request.PatientPathVars](r, chi.URLParam)
	if err != nil {
		h.log.Error("failed to bind path parameters", zap.Error(err))
		helper.WriteUIError(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	// Show edit form
	h.showEditForm(w, r, pathVars.PatientID, form_data.PatientFormData{})
}

// HandleFormSubmission handles POST requests to process form submission
func (h *PatientEditComponent) HandleFormSubmission(w http.ResponseWriter, r *http.Request) {
	// Bind and validate path parameters
	pathVars, _, err := bind.ChiPath[request.PatientPathVars](r, chi.URLParam)
	if err != nil {
		h.log.Error("failed to bind path parameters", zap.Error(err))
		helper.WriteUIError(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	// Process form submission
	h.handleFormSubmission(w, r, pathVars.PatientID)
}

func (h *PatientEditComponent) showEditForm(w http.ResponseWriter, r *http.Request, patientID string, formData form_data.PatientFormData) {
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

	// If no form data provided, populate with current patient data
	if formData.ID == "" {
		formData = form_data.PatientFormData{
			ID:    patient.ID,
			Name:  patient.Name,
			Phone: patient.Phone,
			DOB:   patient.DOB.Format("2006-01-02"),
			State: patient.State,
		}
	}

	view := PatientEditPageComponentView(PatientEditPageParam{
		Patient:    patient,
		FormData:   formData,
		BackPath:   paths.PatientDetailURL(patientID),
		SubmitPath: paths.PatientEditURL(patientID),
	})

	if err := view.Render(r.Context(), w); err != nil {
		h.log.Error("failed to render patient edit", zap.Error(err), zap.String("id", patientID))
		http.Error(w, "failed to render patient edit", http.StatusInternalServerError)
		return
	}
}

func (h *PatientEditComponent) handleFormSubmission(w http.ResponseWriter, r *http.Request, patientID string) {
	// Bind and validate form data
	formReq, fieldErrors, err := bind.Form[request.PatientEditFormRequest](r)
	if err != nil {
		h.log.Error("failed to bind form data", zap.Error(err))
		formData := form_data.PatientFormData{
			ID:     patientID,
			Errors: helper.ConvertFieldErrorsToUIErrors(fieldErrors),
		}
		h.showEditForm(w, r, patientID, formData)
		return
	}

	// Parse DOB (validation already handled by custom validator)
	dob, err := time.Parse("2006-01-02", formReq.DOB)
	if err != nil {
		// This should not happen due to validation, but just in case
		formData := form_data.PatientFormData{
			ID:     patientID,
			Name:   formReq.Name,
			Phone:  formReq.Phone,
			DOB:    formReq.DOB,
			State:  formReq.State,
			Errors: map[string]string{"dob": "Invalid date format"},
		}
		h.showEditForm(w, r, patientID, formData)
		return
	}

	// Get existing patient to preserve other fields
	existingPatient, err := h.patientsService.GetByID(r.Context(), patientID)
	if err != nil {
		h.log.Error("failed to fetch existing patient", zap.Error(err), zap.String("id", patientID))
		helper.WriteUIInternalError(w, "Failed to load patient")
		return
	}

	// Update patient fields
	existingPatient.Name = formReq.Name
	existingPatient.Phone = formReq.Phone
	existingPatient.DOB = dob
	existingPatient.State = formReq.State

	// Save updated patient
	err = h.patientsService.Update(r.Context(), existingPatient)
	if err != nil {
		h.log.Error("failed to update patient", zap.Error(err), zap.String("id", patientID))
		formData := form_data.PatientFormData{
			ID:     patientID,
			Name:   formReq.Name,
			Phone:  formReq.Phone,
			DOB:    formReq.DOB,
			State:  formReq.State,
			Errors: map[string]string{"general": "Failed to update patient. Please try again."},
		}
		h.showEditForm(w, r, patientID, formData)
		return
	}

	// Success - redirect to patient detail page
	http.Redirect(w, r, paths.PatientDetailURL(patientID), http.StatusSeeOther)
}
