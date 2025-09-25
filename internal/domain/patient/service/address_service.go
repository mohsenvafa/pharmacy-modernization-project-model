package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	addressModel "pharmacy-modernization-project-model/internal/domain/patient/contracts/model"
	addressRequest "pharmacy-modernization-project-model/internal/domain/patient/contracts/request"
	addressrepo "pharmacy-modernization-project-model/internal/domain/patient/repository"
)

var ErrInvalidAddress = errors.New("missing required address fields")

type AddressService interface {
	GetByPatientID(ctx context.Context, patientID string) ([]addressModel.Address, error)
	GetByID(ctx context.Context, patientID, addressID string) (addressModel.Address, error)
	Create(ctx context.Context, patientID string, req addressRequest.AddressCreateRequest) (addressModel.Address, error)
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

func (s *addressSvc) Create(ctx context.Context, patientID string, req addressRequest.AddressCreateRequest) (addressModel.Address, error) {
	if strings.TrimSpace(req.Line1) == "" || strings.TrimSpace(req.City) == "" || strings.TrimSpace(req.State) == "" || strings.TrimSpace(req.Zip) == "" {
		return addressModel.Address{}, ErrInvalidAddress
	}

	address := addressModel.Address{
		ID:        uuid.NewString(),
		PatientID: patientID,
		Line1:     req.Line1,
		Line2:     req.Line2,
		City:      req.City,
		State:     req.State,
		Zip:       req.Zip,
	}

	return s.repo.Upsert(ctx, patientID, address)
}

func (s *addressSvc) Upsert(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error) {
	return s.repo.Upsert(ctx, patientID, address)
}
