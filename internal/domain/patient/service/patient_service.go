package service

import (
	"context"

	m "pharmacy-modernization-project-model/internal/domain/patient/contracts/model"
	repo "pharmacy-modernization-project-model/internal/domain/patient/repository"

	"go.uber.org/zap"
)

type PatientService interface {
	List(ctx context.Context, q string, limit, offset int) ([]m.Patient, error)
	GetByID(ctx context.Context, id string) (m.Patient, error)
	Count(ctx context.Context, q string) (int, error)
}

type patientSvc struct {
	repo repo.PatientRepository
	log  *zap.Logger
}

func New(r repo.PatientRepository, l *zap.Logger) PatientService { return &patientSvc{repo: r, log: l} }

func (s *patientSvc) List(ctx context.Context, q string, limit, offset int) ([]m.Patient, error) {
	return s.repo.List(ctx, q, limit, offset)
}
func (s *patientSvc) GetByID(ctx context.Context, id string) (m.Patient, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *patientSvc) Count(ctx context.Context, q string) (int, error) {
	return s.repo.Count(ctx, q)
}
