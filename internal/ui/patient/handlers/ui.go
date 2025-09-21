package handlers

import (
	"net/http"
	"go.uber.org/zap"
	service "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
)

type UI struct { svc service.Service; log *zap.Logger }
func New(s service.Service, l *zap.Logger) *UI { return &UI{svc:s, log:l} }

func (u *UI) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Patients list page (templ placeholder)"))
}
