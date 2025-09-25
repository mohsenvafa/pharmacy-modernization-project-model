package repository

import (
	"context"
	m "pharmacy-modernization-project-model/domain/prescription/model"
)

type PrescriptionRepository interface {
	List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error)
	GetByID(ctx context.Context, id string) (m.Prescription, error)
	Create(ctx context.Context, p m.Prescription) (m.Prescription, error)
	Update(ctx context.Context, id string, p m.Prescription) (m.Prescription, error)
	CountByStatus(ctx context.Context, status string) (int, error)
}
