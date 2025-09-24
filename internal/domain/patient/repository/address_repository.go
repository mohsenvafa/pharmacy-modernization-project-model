package repository

import (
	"context"

	addressModel "pharmacy-modernization-project-model/internal/domain/patient/contracts/model"
)

type AddressRepository interface {
	ListByPatientID(ctx context.Context, patientID string) ([]addressModel.Address, error)
	GetByID(ctx context.Context, patientID, addressID string) (addressModel.Address, error)
	Upsert(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error)
}
