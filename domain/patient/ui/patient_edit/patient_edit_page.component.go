package patient_edit

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/domain/patient/ui/paths"
	"pharmacy-modernization-project-model/domain/patient/validation"
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

// EditFormData represents the form data for editing a patient
type EditFormData struct {
	validation.PatientFormData
	Success    bool
	SuccessMsg string
}

func (h *PatientEditComponent) Handler(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	if patientID == "" {
		http.Error(w, "missing patient id", http.StatusBadRequest)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		h.handleFormSubmission(w, r, patientID)
		return
	}

	// Handle GET request - show edit form
	h.showEditForm(w, r, patientID, EditFormData{})
}

func (h *PatientEditComponent) showEditForm(w http.ResponseWriter, r *http.Request, patientID string, formData EditFormData) {
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
		formData = EditFormData{
			PatientFormData: validation.PatientFormData{
				ID:    patient.ID,
				Name:  patient.Name,
				Phone: patient.Phone,
				DOB:   patient.DOB.Format("2006-01-02"),
				State: patient.State,
			},
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
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	formData := EditFormData{
		PatientFormData: validation.PatientFormData{
			ID:     patientID,
			Name:   r.FormValue("name"),
			Phone:  r.FormValue("phone"),
			DOB:    r.FormValue("dob"),
			State:  r.FormValue("state"),
			Errors: make(map[string]string),
		},
	}

	// Validate form data
	validationErrors := validation.ValidatePatientForm(formData.PatientFormData)
	if len(validationErrors) > 0 {
		formData.Errors = validationErrors
		h.showEditForm(w, r, patientID, formData)
		return
	}

	// Parse DOB
	dob, err := time.Parse("2006-01-02", formData.DOB)
	if err != nil {
		formData.Errors["dob"] = "Invalid date format"
		h.showEditForm(w, r, patientID, formData)
		return
	}

	// Get existing patient to preserve other fields
	existingPatient, err := h.patientsService.GetByID(r.Context(), patientID)
	if err != nil {
		h.log.Error("failed to fetch existing patient", zap.Error(err), zap.String("id", patientID))
		http.Error(w, "failed to load patient", http.StatusInternalServerError)
		return
	}

	// Update patient fields
	existingPatient.Name = formData.Name
	existingPatient.Phone = formData.Phone
	existingPatient.DOB = dob
	existingPatient.State = formData.State

	// Save updated patient
	err = h.patientsService.Update(r.Context(), existingPatient)
	if err != nil {
		h.log.Error("failed to update patient", zap.Error(err), zap.String("id", patientID))
		formData.Errors["general"] = "Failed to update patient. Please try again."
		h.showEditForm(w, r, patientID, formData)
		return
	}

	// Success - redirect to patient detail page
	http.Redirect(w, r, paths.PatientDetailURL(patientID), http.StatusSeeOther)
}
