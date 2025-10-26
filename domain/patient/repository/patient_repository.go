package repository

import (
	"context"

	m "pharmacy-modernization-project-model/domain/patient/contracts/model"
	"pharmacy-modernization-project-model/domain/patient/contracts/request"
)

type PatientRepository interface {
	List(ctx context.Context, req request.PatientListQueryRequest) ([]m.Patient, error)
	GetByID(ctx context.Context, id string) (m.Patient, error)
	Create(ctx context.Context, p m.Patient) (m.Patient, error)
	Update(ctx context.Context, id string, p m.Patient) (m.Patient, error)
	Count(ctx context.Context, req request.PatientListQueryRequest) (int, error)
}
