package repository

import (
	"context"

	addressModel "github.com/pharmacy-modernization-project-model/internal/domain/patient/contracts/model"
)

type addressMemoryRepository struct {
	items map[string]map[string]addressModel.Address
}

func NewAddressMemoryRepository() AddressRepository {
	r := &addressMemoryRepository{items: make(map[string]map[string]addressModel.Address)}

	sample := map[string][]addressModel.Address{
		"P001": {
			{ID: "A001", PatientID: "P001", Line1: "123 Main St", City: "Seattle", State: "WA", Zip: "98101"},
			{ID: "A002", PatientID: "P001", Line1: "456 Market Ave", City: "Seattle", State: "WA", Zip: "98102"},
		},
		"P002": {
			{ID: "A003", PatientID: "P002", Line1: "789 Sunset Blvd", City: "San Francisco", State: "CA", Zip: "94102"},
		},
	}

	for patientID, addresses := range sample {
		for _, addr := range addresses {
			r.Upsert(context.Background(), patientID, addr)
		}
	}

	return r
}

func (r *addressMemoryRepository) ListByPatientID(ctx context.Context, patientID string) ([]addressModel.Address, error) {
	addressesMap, ok := r.items[patientID]
	if !ok {
		return []addressModel.Address{}, nil
	}
	addresses := make([]addressModel.Address, 0, len(addressesMap))
	for _, addr := range addressesMap {
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

func (r *addressMemoryRepository) GetByID(ctx context.Context, patientID, addressID string) (addressModel.Address, error) {
	if addressesMap, ok := r.items[patientID]; ok {
		if addr, ok := addressesMap[addressID]; ok {
			return addr, nil
		}
	}
	return addressModel.Address{}, nil
}

func (r *addressMemoryRepository) Upsert(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error) {
	if address.ID == "" {
		address.ID = patientID + "-addr"
	}
	if _, ok := r.items[patientID]; !ok {
		r.items[patientID] = make(map[string]addressModel.Address)
	}
	r.items[patientID][address.ID] = address
	return address, nil
}
