package repository

import (
	"context"

	addressModel "github.com/pharmacy-modernization-project-model/internal/domain/patient/contracts/model"
)

type addressMemoryRepository struct {
	items map[string]addressModel.Address
}

func NewAddressMemoryRepository() AddressRepository {
	return &addressMemoryRepository{items: map[string]addressModel.Address{}}
}

func (r *addressMemoryRepository) GetByPatientID(ctx context.Context, patientID string) (addressModel.Address, error) {
	if addr, ok := r.items[patientID]; ok {
		return addr, nil
	}
	return addressModel.Address{}, nil
}

func (r *addressMemoryRepository) Update(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error) {
	r.items[patientID] = address
	return address, nil
}
