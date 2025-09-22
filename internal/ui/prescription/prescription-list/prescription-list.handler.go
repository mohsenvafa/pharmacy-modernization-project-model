package prescription_list

import (
	"net/http"

	presSvc "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
	"go.uber.org/zap"
)

type PrescriptionListHandler struct {
	prescriptionsService presSvc.Service
	log                  *zap.Logger
}

func NewPrescriptionListHandler(prescriptions presSvc.Service, log *zap.Logger) *PrescriptionListHandler {
	return &PrescriptionListHandler{prescriptionsService: prescriptions, log: log}
}

func (h *PrescriptionListHandler) Handler(w http.ResponseWriter, r *http.Request) {
	prescriptions, err := h.prescriptionsService.List(r.Context(), "", 1000, 0)
	if err != nil {
		http.Error(w, "failed to load prescriptions", http.StatusInternalServerError)
		return
	}

	page := PrescriptionListPage(PrescriptionListPageParam{
		NumberOfPrescriptions: len(prescriptions),
	})
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render prescription list", http.StatusInternalServerError)
		return
	}
}
