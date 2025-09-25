package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	request "pharmacy-modernization-project-model/domain/prescription/contracts/request"
	response "pharmacy-modernization-project-model/domain/prescription/contracts/response"
	"pharmacy-modernization-project-model/domain/prescription/service"
	helper "pharmacy-modernization-project-model/internal/helper"
)

type PrescriptionController struct {
	svc service.PrescriptionService
	log *zap.Logger
}

func NewPrescriptionController(s service.PrescriptionService, log *zap.Logger) *PrescriptionController {
	return &PrescriptionController{svc: s, log: log}
}

func (c *PrescriptionController) RegisterRoutes(r chi.Router) {
	r.Get("/", c.List)
	r.Get("/{prescriptionID}", c.GetByID)
}

func (c *PrescriptionController) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := request.PrescriptionListQueryRequest{
		Status: r.URL.Query().Get("status"),
		Limit:  limit,
		Offset: offset,
	}

	items, err := c.svc.List(r.Context(), query.Status, query.Limit, query.Offset)
	if err != nil {
		c.log.Error("list prescriptions", zap.Error(err))
		if writeErr := helper.WriteError(w, http.StatusInternalServerError, helper.APIError{
			Code:    "prescription_list_error",
			Message: "failed to list prescriptions",
		}); writeErr != nil {
			c.log.Error("write response", zap.Error(writeErr))
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, response.FromModels(items), nil); err != nil {
		c.log.Error("write response", zap.Error(err))
	}
}

func (c *PrescriptionController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "prescriptionID")
	item, err := c.svc.GetByID(r.Context(), id)
	if err != nil {
		c.log.Error("get prescription", zap.Error(err), zap.String("prescriptionID", id))
		if writeErr := helper.WriteError(w, http.StatusInternalServerError, helper.APIError{
			Code:    "prescription_get_error",
			Message: "failed to fetch prescription",
		}); writeErr != nil {
			c.log.Error("write response", zap.Error(writeErr))
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, response.FromModel(item), nil); err != nil {
		c.log.Error("write response", zap.Error(err))
	}
}
