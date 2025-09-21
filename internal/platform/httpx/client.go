package httpx

import (
	"net/http"
	"time"

	"github.com/pharmacy-modernization-project-model/internal/platform/config"
)

func NewClient(cfg *config.Config) *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns: 10,
			IdleConnTimeout: 30 * time.Second,
		},
	}
}
