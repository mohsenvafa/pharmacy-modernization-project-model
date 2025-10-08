package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patientsecurity "pharmacy-modernization-project-model/domain/patient/security"
	service "pharmacy-modernization-project-model/domain/patient/service"
	helper "pharmacy-modernization-project-model/internal/helper"
	"pharmacy-modernization-project-model/internal/platform/auth"
	"pharmacy-modernization-project-model/internal/platform/httpx"
)

type PatientController struct {
	patientService service.PatientService
	log            *zap.Logger
}

func NewPatientController(patients service.PatientService, log *zap.Logger) *PatientController {
	return &PatientController{patientService: patients, log: log}
}

func (c *PatientController) RegisterRoutes(r chi.Router) {
	// All patient API routes require authentication (header-based for API)
	r.Use(auth.RequireAuthFromHeader())

	// Read operations - requires patient:read or admin:all
	r.With(auth.RequirePermissionsMatchAny(patientsecurity.ReadAccess)).Get("/", c.List)
	r.With(auth.RequirePermissionsMatchAny(patientsecurity.ReadAccess)).Get("/{patientID}", c.GetByID)
}

func (c *PatientController) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 20
	}

	items, err := c.patientService.List(r.Context(), q, limit, offset)
	if err != nil {
		c.log.Error("list patients", zap.Error(err))
		if writeErr := helper.WriteError(w, http.StatusInternalServerError, helper.APIError{
			Code:    "patient_list_error",
			Message: "failed to list patients",
		}); writeErr != nil {
			c.log.Error("write response", zap.Error(writeErr))
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, items, nil); err != nil {
		c.log.Error("write response", zap.Error(err))
	}
}

func (c *PatientController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "patientID")

	item, err := c.patientService.GetByID(r.Context(), id)
	if err != nil {
		c.log.Error("get patient", zap.Error(err), zap.String("patientID", id))
		c.handleError(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, item, nil); err != nil {
		c.log.Error("write response", zap.Error(err))
	}
}

// handleError handles different types of errors and returns appropriate HTTP responses
func (c *PatientController) handleError(w http.ResponseWriter, r *http.Request, err error) {
	// Use the shared error handler
	httpx.WriteError(w, r, err)
}
