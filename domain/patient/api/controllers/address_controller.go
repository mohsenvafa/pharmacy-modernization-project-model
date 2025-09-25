package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	addressRequest "pharmacy-modernization-project-model/domain/patient/contracts/request"
	service "pharmacy-modernization-project-model/domain/patient/service"
	helper "pharmacy-modernization-project-model/internal/helper"
)

type AddressController struct {
	addressService service.AddressService
	log            *zap.Logger
}

func NewAddressController(addresses service.AddressService, log *zap.Logger) *AddressController {
	return &AddressController{addressService: addresses, log: log}
}

func (c *AddressController) RegisterRoutes(r chi.Router) {
	r.Get("/", c.ListByPatient)
	r.Get("/{addressID}", c.GetByID)
	r.Post("/", c.Create)
}

func (c *AddressController) ListByPatient(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	addr, err := c.addressService.GetByPatientID(r.Context(), patientID)
	if err != nil {
		c.log.Error("list addresses", zap.Error(err), zap.String("patientID", patientID))
		if writeErr := helper.WriteError(w, http.StatusInternalServerError, helper.APIError{
			Code:    "address_list_error",
			Message: "failed to load addresses",
		}); writeErr != nil {
			c.log.Error("write response", zap.Error(writeErr))
		}
		return
	}
	if err := helper.WriteJSON(w, http.StatusOK, addr, nil); err != nil {
		c.log.Error("write response", zap.Error(err))
	}
}

func (c *AddressController) GetByID(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	addressID := chi.URLParam(r, "addressID")

	address, err := c.addressService.GetByID(r.Context(), patientID, addressID)
	if err != nil {
		c.log.Error("get address", zap.Error(err), zap.String("patientID", patientID), zap.String("addressID", addressID))
		if writeErr := helper.WriteError(w, http.StatusInternalServerError, helper.APIError{
			Code:    "address_get_error",
			Message: "failed to load address",
		}); writeErr != nil {
			c.log.Error("write response", zap.Error(writeErr))
		}
		return
	}
	if err := helper.WriteJSON(w, http.StatusOK, address, nil); err != nil {
		c.log.Error("write response", zap.Error(err))
	}
}

func (c *AddressController) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req addressRequest.AddressCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.Warn("invalid address payload", zap.Error(err))
		if writeErr := helper.WriteError(w, http.StatusBadRequest, helper.APIError{
			Code:    "invalid_request",
			Message: "invalid address payload",
		}); writeErr != nil {
			c.log.Error("write response", zap.Error(writeErr))
		}
		return
	}

	patientID := chi.URLParam(r, "patientID")
	created, err := c.addressService.Create(r.Context(), patientID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAddress) {
			if writeErr := helper.WriteError(w, http.StatusBadRequest, helper.APIError{
				Code:    "invalid_request",
				Message: err.Error(),
			}); writeErr != nil {
				c.log.Error("write response", zap.Error(writeErr))
			}
			return
		}
		c.log.Error("create address", zap.Error(err), zap.String("patientID", patientID))
		if writeErr := helper.WriteError(w, http.StatusInternalServerError, helper.APIError{
			Code:    "address_create_error",
			Message: "failed to create address",
		}); writeErr != nil {
			c.log.Error("write response", zap.Error(writeErr))
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusCreated, created, nil); err != nil {
		c.log.Error("write response", zap.Error(err))
	}
}
