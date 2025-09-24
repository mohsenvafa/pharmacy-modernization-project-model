package service

import (
	"context"

	addressModel "pharmacy-modernization-project-model/internal/domain/patient/contracts/model"
	addressrepo "pharmacy-modernization-project-model/internal/domain/patient/repository"
)

type AddressService interface {
	GetByPatientID(ctx context.Context, patientID string) ([]addressModel.Address, error)
	GetByID(ctx context.Context, patientID, addressID string) (addressModel.Address, error)
	Upsert(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error)
}

type addressSvc struct {
	repo addressrepo.AddressRepository
}

func NewAddressService(r addressrepo.AddressRepository) AddressService {
	return &addressSvc{repo: r}
}

func (s *addressSvc) GetByPatientID(ctx context.Context, patientID string) ([]addressModel.Address, error) {
	return s.repo.ListByPatientID(ctx, patientID)
}

func (s *addressSvc) GetByID(ctx context.Context, patientID, addressID string) (addressModel.Address, error) {
	return s.repo.GetByID(ctx, patientID, addressID)
}

func (s *addressSvc) Upsert(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error) {
	return s.repo.Upsert(ctx, patientID, address)
}
