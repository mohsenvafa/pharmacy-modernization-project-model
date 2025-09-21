package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	service "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	"go.uber.org/zap"
)

type API struct { svc service.Service; log *zap.Logger }

func New(s service.Service, l *zap.Logger) *API { return &API{svc:s, log:l} }

func (a *API) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limit,_ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset,_ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit==0 {limit=20}
	items, _ := a.svc.List(r.Context(), q, limit, offset)
	_ = json.NewEncoder(w).Encode(items)
}

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, _ := a.svc.GetByID(r.Context(), id)
	_ = json.NewEncoder(w).Encode(item)
}
