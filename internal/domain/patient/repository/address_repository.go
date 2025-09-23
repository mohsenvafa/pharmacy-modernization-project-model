package repository

import (
	"context"

	addressModel "github.com/pharmacy-modernization-project-model/internal/domain/patient/contracts/model"
)

type AddressRepository interface {
	GetByPatientID(ctx context.Context, patientID string) (addressModel.Address, error)
	Update(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error)
}
