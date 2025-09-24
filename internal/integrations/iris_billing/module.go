package iris_billing

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ModuleDependencies struct {
	Config     Config
	Logger     *zap.Logger
	HTTPClient *http.Client
	UseMock    bool
	Timeout    time.Duration
}

type ModuleExport struct {
	Client Client
}

func Module(deps ModuleDependencies) ModuleExport {
	if deps.UseMock {
		return ModuleExport{Client: NewMockClient(nil, deps.Logger)}
	}
	httpClient := deps.HTTPClient
	if httpClient == nil {
		timeout := deps.Timeout
		if timeout <= 0 {
			timeout = 10 * time.Second
		}
		httpClient = &http.Client{Timeout: timeout}
	}
	client := NewHTTPClient(deps.Config, httpClient, deps.Logger)
	return ModuleExport{Client: client}
}
