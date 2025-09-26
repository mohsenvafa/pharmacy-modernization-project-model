package addresslist

import (
	"context"
	"errors"

	patSvc "pharmacy-modernization-project-model/domain/patient/service"

	"github.com/a-h/templ"
	"go.uber.org/zap"
)

type AddressListComponentHandler struct {
	addresses patSvc.AddressService
	log       *zap.Logger
}

func NewAddressListComponentHandler(addresses patSvc.AddressService, log *zap.Logger) *AddressListComponentHandler {
	return &AddressListComponentHandler{addresses: addresses, log: log}
}

func (h *AddressListComponentHandler) Handler(ctx context.Context, patientID string) (templ.Component, error) {
	if patientID == "" {
		return nil, errors.New("patient id is required")
	}

	addresses, err := h.addresses.GetByPatientID(ctx, patientID)
	if err != nil {
		if h.log != nil {
			h.log.Error("failed to load patient addresses", zap.Error(err), zap.String("patient_id", patientID))
		}
		return nil, err
	}

	return AddressListComponent(AddressListParams{
		Title:        "Addresses",
		EmptyMessage: "No addresses on file for this patient.",
		Addresses:    addresses,
	}), nil
}
