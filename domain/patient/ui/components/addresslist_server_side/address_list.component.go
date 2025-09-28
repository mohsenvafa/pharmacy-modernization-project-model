package addresslist

import (
	"context"
	"errors"

	"github.com/a-h/templ"
	"go.uber.org/zap"

	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
)

type AddressListComponent struct {
	service patSvc.AddressService
	log     *zap.Logger
}

func NewAddressListComponent(deps *contracts.UiDependencies) *AddressListComponent {
	return &AddressListComponent{service: deps.AddressSvc, log: deps.Log}
}

func (c *AddressListComponent) View(ctx context.Context, patientID string) (templ.Component, error) {
	if patientID == "" {
		return nil, errors.New("patient id is required")
	}

	addresses, err := c.service.GetByPatientID(ctx, patientID)
	if err != nil {
		if c.log != nil {
			c.log.Error("failed to load patient addresses", zap.Error(err), zap.String("patient_id", patientID))
		}
		return nil, err
	}

	return AddressListComponentView(AddressListParams{
		Title:        "Addresses",
		EmptyMessage: "No addresses on file for this patient.",
		Addresses:    addresses,
	}), nil
}
