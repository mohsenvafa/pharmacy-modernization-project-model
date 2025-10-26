package prescription_list

import (
	"net/http"

	presSvc "pharmacy-modernization-project-model/domain/prescription/service"
	helper "pharmacy-modernization-project-model/internal/helper"

	"go.uber.org/zap"
)

type PrescriptionListHandler struct {
	prescriptionsService presSvc.PrescriptionService
	log                  *zap.Logger
}

func NewPrescriptionListHandler(prescriptions presSvc.PrescriptionService, log *zap.Logger) *PrescriptionListHandler {
	return &PrescriptionListHandler{prescriptionsService: prescriptions, log: log}
}

func (h *PrescriptionListHandler) Handler(w http.ResponseWriter, r *http.Request) {
	prescriptions, err := h.prescriptionsService.List(r.Context(), "", 1000, 0)
	if err != nil {
		h.log.Error("failed to load prescriptions", zap.Error(err))
		helper.WriteUIInternalError(w, "Failed to load prescriptions")
		return
	}

	page := PrescriptionListPageComponent(PrescriptionListPageParam{
		NumberOfPrescriptions: len(prescriptions),
	})
	if err := page.Render(r.Context(), w); err != nil {
		h.log.Error("failed to render prescription list", zap.Error(err))
		helper.WriteUIInternalError(w, "Failed to render prescription list")
		return
	}
}
