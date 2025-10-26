package controllers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	addressRequest "pharmacy-modernization-project-model/domain/patient/contracts/request"
	patientsecurity "pharmacy-modernization-project-model/domain/patient/security"
	service "pharmacy-modernization-project-model/domain/patient/service"
	"pharmacy-modernization-project-model/internal/bind"
	helper "pharmacy-modernization-project-model/internal/helper"
	"pharmacy-modernization-project-model/internal/platform/auth"
)

type AddressController struct {
	addressService service.AddressService
	log            *zap.Logger
}

func NewAddressController(addresses service.AddressService, log *zap.Logger) *AddressController {
	return &AddressController{addressService: addresses, log: log}
}

func (c *AddressController) RegisterRoutes(r chi.Router) {
	// Address routes inherit auth from parent but add specific permissions

	// Read operations - requires patient:read or admin:all
	r.With(auth.RequirePermissionsMatchAny(patientsecurity.ReadAccess)).Get("/", c.ListByPatient)
	r.With(auth.RequirePermissionsMatchAny(patientsecurity.ReadAccess)).Get("/{addressID}", c.GetByID)

	// Write operations - requires patient:write or admin:all
	r.With(auth.RequirePermissionsMatchAny(patientsecurity.WriteAccess)).Post("/", c.Create)
}

func (c *AddressController) ListByPatient(w http.ResponseWriter, r *http.Request) {
	// Bind and validate path parameters
	pathVars, fieldErrors, err := bind.ChiPath[addressRequest.PatientPathVars](r, chi.URLParam)
	if err != nil {
		c.log.Error("failed to bind path parameters", zap.Error(err))
		helper.Respond400(w, fieldErrors)
		return
	}

	addr, err := c.addressService.GetByPatientID(r.Context(), pathVars.PatientID)
	if err != nil {
		c.log.Error("list addresses", zap.Error(err))
		helper.WriteInternalError(w, "failed to load addresses")
		return
	}
	helper.WriteOK(w, addr)
}

func (c *AddressController) GetByID(w http.ResponseWriter, r *http.Request) {
	// Bind and validate path parameters
	pathVars, fieldErrors, err := bind.ChiPath[addressRequest.AddressPathVars](r, chi.URLParam)
	if err != nil {
		c.log.Error("failed to bind path parameters", zap.Error(err))
		helper.Respond400(w, fieldErrors)
		return
	}

	address, err := c.addressService.GetByID(r.Context(), pathVars.PatientID, pathVars.AddressID)
	if err != nil {
		c.log.Error("get address", zap.Error(err))
		helper.WriteInternalError(w, "failed to load address")
		return
	}
	helper.WriteOK(w, address)
}

func (c *AddressController) Create(w http.ResponseWriter, r *http.Request) {
	// Bind and validate path parameters
	pathVars, fieldErrors, err := bind.ChiPath[addressRequest.PatientPathVars](r, chi.URLParam)
	if err != nil {
		c.log.Error("failed to bind path parameters", zap.Error(err))
		helper.Respond400(w, fieldErrors)
		return
	}

	// Bind and validate JSON body
	req, fieldErrors, err := bind.JSON[addressRequest.AddressCreateRequest](r)
	if err != nil {
		c.log.Warn("invalid address payload", zap.Error(err))
		helper.Respond400(w, fieldErrors)
		return
	}

	created, err := c.addressService.Create(r.Context(), pathVars.PatientID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAddress) {
			helper.WriteError(w, http.StatusBadRequest, helper.APIError{
				Code:    "invalid_request",
				Message: err.Error(),
			})
			return
		}
		c.log.Error("create address", zap.Error(err))
		helper.WriteInternalError(w, "failed to create address")
		return
	}

	helper.WriteCreated(w, created)
}
