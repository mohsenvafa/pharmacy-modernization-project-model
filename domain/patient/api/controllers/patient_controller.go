package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/domain/patient/contracts/request"
	patientsecurity "pharmacy-modernization-project-model/domain/patient/security"
	service "pharmacy-modernization-project-model/domain/patient/service"
	"pharmacy-modernization-project-model/internal/bind"
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
	// Bind and validate query parameters
	req, fieldErrors, err := bind.Query[request.PatientListQueryRequest](r)
	if err != nil {
		c.log.Error("failed to bind query parameters", zap.Error(err))
		helper.Respond400(w, fieldErrors)
		return
	}

	// Set default limit if not provided
	if req.Limit == 0 {
		req.Limit = 20
	}

	items, err := c.patientService.List(r.Context(), req)
	if err != nil {
		c.log.Error("list patients", zap.Error(err))
		helper.WriteInternalError(w, "failed to list patients")
		return
	}

	helper.WriteOK(w, items)
}

func (c *PatientController) GetByID(w http.ResponseWriter, r *http.Request) {
	// Bind and validate path parameters
	pathVars, fieldErrors, err := bind.ChiPath[request.PatientPathVars](r, chi.URLParam)
	if err != nil {
		c.log.Error("failed to bind path parameters", zap.Error(err))
		helper.Respond400(w, fieldErrors)
		return
	}

	item, err := c.patientService.GetByID(r.Context(), pathVars.PatientID)
	if err != nil {
		c.log.Error("get patient", zap.Error(err))
		c.handleError(w, r, err)
		return
	}

	helper.WriteOK(w, item)
}

// handleError handles different types of errors and returns appropriate HTTP responses
func (c *PatientController) handleError(w http.ResponseWriter, r *http.Request, err error) {
	// Use the shared error handler
	httpx.WriteError(w, r, err)
}
