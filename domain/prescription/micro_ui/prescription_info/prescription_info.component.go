package prescriptioninfo

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	prescriptionservice "pharmacy-modernization-project-model/domain/prescription/service"
)

// Dependencies enumerates the dependencies required by the micro UI handler.
type Dependencies struct {
	Service prescriptionservice.PrescriptionService
	Log     *zap.Logger
}

// Handler exposes an HTTP handler for rendering prescription info fragments.
type Handler struct {
	svc prescriptionservice.PrescriptionService
	log *zap.Logger
}

// NewHandler constructs a new Handler using the provided dependencies.
func NewHandler(deps *Dependencies) *Handler {
	return &Handler{
		svc: deps.Service,
		log: deps.Log,
	}
}

// Handle renders the prescription info fragment based on request parameters.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	prescriptionID := strings.TrimSpace(chi.URLParam(r, "prescriptionID"))
	if prescriptionID == "" {
		h.renderError(ctx, w, r, http.StatusBadRequest, errors.New("missing prescription id"))
		return
	}

	authToken := strings.TrimSpace(r.URL.Query().Get("auth_token"))
	if authToken == "" {
		h.renderError(ctx, w, r, http.StatusUnauthorized, errors.New("missing auth_token"))
		return
	}

	prescription, err := h.svc.GetByID(ctx, prescriptionID)
	if err != nil {
		if h.log != nil {
			h.log.Error("failed to load prescription for micro ui",
				zap.Error(err),
				zap.String("prescription_id", prescriptionID),
			)
		}

		h.renderError(ctx, w, r, http.StatusInternalServerError, errors.New("failed to load prescription"))
		return
	}

	view := PrescriptionInfoFragment(PrescriptionInfoParams{
		Prescription: prescription,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := view.Render(ctx, w); err != nil && h.log != nil {
		h.log.Error("failed to render prescription info fragment", zap.Error(err))
	}
}

func (h *Handler) renderError(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, err error) {
	addCORSHeaders(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if h.log != nil {
		h.log.Warn("micro ui prescription info error", zap.Error(err))
	}

	view := PrescriptionInfoError(err.Error())
	_ = view.Render(ctx, w)
}

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
}
