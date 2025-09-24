package integrations

import (
	"time"

	"go.uber.org/zap"

	irisbilling "github.com/pharmacy-modernization-project-model/internal/integrations/iris_billing"
	irispharmacy "github.com/pharmacy-modernization-project-model/internal/integrations/iris_pharmacy"
	"github.com/pharmacy-modernization-project-model/internal/platform/config"
)

type Dependencies struct {
	Config *config.Config
	Logger *zap.Logger
}

type Export struct {
	Pharmacy irispharmacy.Client
	Billing  irisbilling.Client
}

func New(deps Dependencies) Export {
	if deps.Config == nil {
		return Export{}
	}

	pharmacy := irispharmacy.Module(irispharmacy.ModuleDependencies{
		Config: irispharmacy.Config{
			BaseURL: deps.Config.External.Pharmacy.BaseURL,
			Path:    deps.Config.External.Pharmacy.Path,
		},
		Logger:  deps.Logger,
		UseMock: deps.Config.External.Pharmacy.UseMock,
		Timeout: parseDuration(deps.Config.External.Pharmacy.Timeout, 10*time.Second),
	}).Client

	billing := irisbilling.Module(irisbilling.ModuleDependencies{
		Config: irisbilling.Config{
			BaseURL: deps.Config.External.Billing.BaseURL,
			Path:    deps.Config.External.Billing.Path,
		},
		Logger:  deps.Logger,
		UseMock: deps.Config.External.Billing.UseMock,
		Timeout: parseDuration(deps.Config.External.Billing.Timeout, 10*time.Second),
	}).Client

	return Export{Pharmacy: pharmacy, Billing: billing}
}

func parseDuration(value string, fallback time.Duration) time.Duration {
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	return fallback
}
