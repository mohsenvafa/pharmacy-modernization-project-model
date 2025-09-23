package service

import (
	"context"

	addressModel "github.com/pharmacy-modernization-project-model/internal/domain/patient/contracts/model"
	addressrepo "github.com/pharmacy-modernization-project-model/internal/domain/patient/repository"
)

type AddressService interface {
	GetByPatientID(ctx context.Context, patientID string) (addressModel.Address, error)
	Update(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error)
}

type addressSvc struct {
	repo addressrepo.AddressRepository
}

func NewAddressService(r addressrepo.AddressRepository) AddressService {
	return &addressSvc{repo: r}
}

func (s *addressSvc) GetByPatientID(ctx context.Context, patientID string) (addressModel.Address, error) {
	return s.repo.GetByPatientID(ctx, patientID)
}

func (s *addressSvc) Update(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error) {
	return s.repo.Update(ctx, patientID, address)
}
