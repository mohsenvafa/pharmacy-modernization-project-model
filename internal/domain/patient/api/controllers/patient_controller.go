package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	service "pharmacy-modernization-project-model/internal/domain/patient/service"
)

type PatientController struct {
	patientService service.PatientService
	log            *zap.Logger
}

func NewPatientController(patients service.PatientService, log *zap.Logger) *PatientController {
	return &PatientController{patientService: patients, log: log}
}

func (c *PatientController) RegisterRoutes(r chi.Router) {
	r.Get("/", c.List)
	r.Get("/{patientID}", c.GetByID)
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
		http.Error(w, "failed to list patients", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(items)
}

func (c *PatientController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "patientID")
	item, err := c.patientService.GetByID(r.Context(), id)
	if err != nil {
		c.log.Error("get patient", zap.Error(err), zap.String("patientID", id))
		http.Error(w, "failed to fetch patient", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(item)
}
