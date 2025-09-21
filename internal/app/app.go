package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/pharmacy-modernization-project-model/internal/platform/config"
	"github.com/pharmacy-modernization-project-model/internal/platform/logging"
)

type App struct {
	Cfg    *config.Config
	Logger *logging.LoggerBundle
	Router chi.Router
	Server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	app := &App{Cfg: cfg}
	if err := app.wire(); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) Run() error {
	a.Server = &http.Server{
		Addr:              fmt.Sprintf(":%d", a.Cfg.App.Port),
		Handler:           a.Router,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return a.Server.ListenAndServe()
}
