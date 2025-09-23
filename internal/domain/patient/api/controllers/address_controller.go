package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	service "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
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
}

func (c *AddressController) ListByPatient(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	addr, err := c.addressService.GetByPatientID(r.Context(), patientID)
	if err != nil {
		c.log.Error("list addresses", zap.Error(err), zap.String("patientID", patientID))
		http.Error(w, "failed to load addresses", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(addr)
}

func (c *AddressController) GetByID(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	addressID := chi.URLParam(r, "addressID")

	address, err := c.addressService.GetByID(r.Context(), patientID, addressID)
	if err != nil {
		c.log.Error("get address", zap.Error(err), zap.String("patientID", patientID), zap.String("addressID", addressID))
		http.Error(w, "failed to load address", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(address)
}
