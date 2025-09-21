package repository

import (
	"context"
	m "github.com/pharmacy-modernization-project-model/internal/domain/patient/model"
)

type Repository interface {
	List(ctx context.Context, query string, limit, offset int) ([]m.Patient, error)
	GetByID(ctx context.Context, id string) (m.Patient, error)
	Create(ctx context.Context, p m.Patient) (m.Patient, error)
	Update(ctx context.Context, id string, p m.Patient) (m.Patient, error)
}
