package repository

import (
	"context"
	m "github.com/pharmacy-modernization-project-model/internal/domain/prescription/model"
)

type PrescriptionRepository interface {
	List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error)
	GetByID(ctx context.Context, id string) (m.Prescription, error)
	Create(ctx context.Context, p m.Prescription) (m.Prescription, error)
	Update(ctx context.Context, id string, p m.Prescription) (m.Prescription, error)
}
