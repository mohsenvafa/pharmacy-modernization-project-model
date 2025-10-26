package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	request "pharmacy-modernization-project-model/domain/prescription/contracts/request"
	response "pharmacy-modernization-project-model/domain/prescription/contracts/response"
	prescriptionsecurity "pharmacy-modernization-project-model/domain/prescription/security"
	"pharmacy-modernization-project-model/domain/prescription/service"
	"pharmacy-modernization-project-model/internal/bind"
	helper "pharmacy-modernization-project-model/internal/helper"
	"pharmacy-modernization-project-model/internal/platform/auth"
)

type PrescriptionController struct {
	svc service.PrescriptionService
	log *zap.Logger
}

func NewPrescriptionController(s service.PrescriptionService, log *zap.Logger) *PrescriptionController {
	return &PrescriptionController{svc: s, log: log}
}

func (c *PrescriptionController) RegisterRoutes(r chi.Router) {
	// All prescription API routes require authentication (header-based for API)
	r.Use(auth.RequireAuthFromHeader())

	// Read operations - requires prescription:read or healthcare role or admin
	r.With(auth.RequirePermissionsMatchAny(prescriptionsecurity.ReadAccess)).Get("/", c.List)
	r.With(auth.RequirePermissionsMatchAny(prescriptionsecurity.ReadAccess)).Get("/{prescriptionID}", c.GetByID)
}

func (c *PrescriptionController) List(w http.ResponseWriter, r *http.Request) {
	// Bind and validate query parameters
	req, fieldErrors, err := bind.Query[request.PrescriptionListQueryRequest](r)
	if err != nil {
		c.log.Error("failed to bind query parameters", zap.Error(err))
		helper.Respond400(w, fieldErrors)
		return
	}

	// Set default limit if not provided
	if req.Limit == 0 {
		req.Limit = 20
	}

	items, err := c.svc.List(r.Context(), req.Status, req.Limit, req.Offset)
	if err != nil {
		c.log.Error("list prescriptions", zap.Error(err))
		helper.WriteInternalError(w, "failed to list prescriptions")
		return
	}

	helper.WriteOK(w, response.FromModels(items))
}

func (c *PrescriptionController) GetByID(w http.ResponseWriter, r *http.Request) {
	// Bind and validate path parameters
	pathVars, fieldErrors, err := bind.ChiPath[request.PrescriptionPathVars](r, chi.URLParam)
	if err != nil {
		c.log.Error("failed to bind path parameters", zap.Error(err))
		helper.Respond400(w, fieldErrors)
		return
	}

	item, err := c.svc.GetByID(r.Context(), pathVars.PrescriptionID)
	if err != nil {
		c.log.Error("get prescription", zap.Error(err), zap.String("prescriptionID", pathVars.PrescriptionID))
		helper.WriteInternalError(w, "failed to fetch prescription")
		return
	}

	helper.WriteOK(w, response.FromModel(item))
}
