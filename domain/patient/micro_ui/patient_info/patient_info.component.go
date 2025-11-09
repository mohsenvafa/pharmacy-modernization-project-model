package patientinfo

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patientservice "pharmacy-modernization-project-model/domain/patient/service"
)

type Dependencies struct {
	Service patientservice.PatientService
	Log     *zap.Logger
}

type Handler struct {
	svc patientservice.PatientService
	log *zap.Logger
}

func NewHandler(deps *Dependencies) *Handler {
	return &Handler{
		svc: deps.Service,
		log: deps.Log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	patientID := strings.TrimSpace(chi.URLParam(r, "patientID"))
	if patientID == "" {
		h.renderError(ctx, w, http.StatusBadRequest, errors.New("missing patient id"))
		return
	}

	authToken := strings.TrimSpace(r.URL.Query().Get("auth_token"))
	if authToken == "" {
		h.renderError(ctx, w, http.StatusUnauthorized, errors.New("missing auth_token"))
		return
	}

	patient, err := h.svc.GetByID(ctx, patientID)
	if err != nil {
		if h.log != nil {
			h.log.Error("failed to load patient for micro ui",
				zap.Error(err),
				zap.String("patient_id", patientID),
			)
		}

		h.renderError(ctx, w, http.StatusInternalServerError, errors.New("failed to load patient"))
		return
	}

	view := PatientInfoFragment(PatientInfoParams{
		Patient: patient,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := view.Render(ctx, w); err != nil && h.log != nil {
		h.log.Error("failed to render patient info fragment", zap.Error(err))
	}
}

func (h *Handler) renderError(ctx context.Context, w http.ResponseWriter, status int, err error) {
	addCORSHeaders(w)
	w.WriteHeader(status)

	if h.log != nil {
		h.log.Warn("micro ui patient info error", zap.Error(err))
	}

	view := PatientInfoError(err.Error())
	_ = view.Render(ctx, w)
}

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
}
