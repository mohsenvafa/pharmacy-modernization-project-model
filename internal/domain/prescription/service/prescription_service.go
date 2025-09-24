package service

import (
	"context"

	m "pharmacy-modernization-project-model/internal/domain/prescription/model"
	repo "pharmacy-modernization-project-model/internal/domain/prescription/repository"
	irisbilling "pharmacy-modernization-project-model/internal/integrations/iris_billing"
	irispharmacy "pharmacy-modernization-project-model/internal/integrations/iris_pharmacy"

	"go.uber.org/zap"
)

type PrescriptionService interface {
	List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error)
	GetByID(ctx context.Context, id string) (m.Prescription, error)
	CountByStatus(ctx context.Context, status string) (int, error)
}

type svc struct {
	repo     repo.PrescriptionRepository
	log      *zap.Logger
	pharmacy irispharmacy.Client
	billing  irisbilling.Client
}

func New(r repo.PrescriptionRepository, l *zap.Logger, pharmacy irispharmacy.Client, billing irisbilling.Client) PrescriptionService {
	return &svc{repo: r, log: l, pharmacy: pharmacy, billing: billing}
}

func (s *svc) List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error) {
	return s.repo.List(ctx, status, limit, offset)
}
func (s *svc) GetByID(ctx context.Context, id string) (m.Prescription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *svc) CountByStatus(ctx context.Context, status string) (int, error) {
	return s.repo.CountByStatus(ctx, status)
}
