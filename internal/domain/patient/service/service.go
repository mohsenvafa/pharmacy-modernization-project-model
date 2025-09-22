package service

import (
	"context"

	m "github.com/pharmacy-modernization-project-model/internal/domain/patient/model"
	repo "github.com/pharmacy-modernization-project-model/internal/domain/patient/repository"
	"go.uber.org/zap"
)

type Service interface {
	List(ctx context.Context, q string, limit, offset int) ([]m.Patient, error)
	GetByID(ctx context.Context, id string) (m.Patient, error)
}

type svc struct {
	repo repo.PatientRepository
	log  *zap.Logger
}

func New(r repo.PatientRepository, l *zap.Logger) Service { return &svc{repo: r, log: l} }

func (s *svc) List(ctx context.Context, q string, limit, offset int) ([]m.Patient, error) {
	return s.repo.List(ctx, q, limit, offset)
}
func (s *svc) GetByID(ctx context.Context, id string) (m.Patient, error) {
	return s.repo.GetByID(ctx, id)
}
