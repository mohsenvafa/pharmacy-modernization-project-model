package patient_edit

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/domain/patient/ui/paths"
	tools "pharmacy-modernization-project-model/internal/helper"
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
	ID         string
	Name       string
	Phone      string
	DOB        string // Format: YYYY-MM-DD
	State      string
	Errors     map[string]string
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
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	formData := EditFormData{
		ID:     patientID,
		Name:   r.FormValue("name"),
		Phone:  r.FormValue("phone"),
		DOB:    r.FormValue("dob"),
		State:  r.FormValue("state"),
		Errors: make(map[string]string),
	}

	// Validate form data
	validationErrors := h.validateFormData(formData)
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

func (h *PatientEditComponent) validateFormData(formData EditFormData) map[string]string {
	errors := make(map[string]string)

	// Validate name
	if formData.Name == "" {
		errors["name"] = "Name is required"
	} else if len(formData.Name) < 2 {
		errors["name"] = "Name must be at least 2 characters"
	} else if len(formData.Name) > 100 {
		errors["name"] = "Name must be less than 100 characters"
	}

	// Validate phone
	if formData.Phone == "" {
		errors["phone"] = "Phone number is required"
	} else if !h.isValidPhone(formData.Phone) {
		errors["phone"] = "Please enter a valid phone number (e.g., (555) 123-4567)"
	}

	// Validate DOB
	if formData.DOB == "" {
		errors["dob"] = "Date of birth is required"
	} else {
		dob, err := time.Parse("2006-01-02", formData.DOB)
		if err != nil {
			errors["dob"] = "Please enter a valid date"
		} else {
			// Check if DOB is not in the future
			if dob.After(time.Now()) {
				errors["dob"] = "Date of birth cannot be in the future"
			}
			// Check if patient is not too old (reasonable limit)
			age := tools.CalculateAge(dob)
			if age > 150 {
				errors["dob"] = "Please enter a valid date of birth"
			}
		}
	}

	// Validate state
	if formData.State == "" {
		errors["state"] = "State is required"
	} else if len(formData.State) > 50 {
		errors["state"] = "State must be less than 50 characters"
	}

	return errors
}

func (h *PatientEditComponent) isValidPhone(phone string) bool {
	// Simple phone validation - allows various formats
	// Remove all non-digit characters
	digits := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			digits += string(char)
		}
	}
	// Check if we have 10 digits (US phone number)
	return len(digits) == 10
}
