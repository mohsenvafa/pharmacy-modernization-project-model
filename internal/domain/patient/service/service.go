package service

import (
	"context"
	repo "github.com/pharmacy-modernization-project-model/internal/domain/patient/repository"
	m "github.com/pharmacy-modernization-project-model/internal/domain/patient/model"
	"go.uber.org/zap"
)

type Service interface {
	List(ctx context.Context, q string, limit, offset int) ([]m.Patient, error)
	GetByID(ctx context.Context, id string) (m.Patient, error)
}

type svc struct { repo repo.Repository; log *zap.Logger }

func New(r repo.Repository, l *zap.Logger) Service { return &svc{repo:r, log:l} }

func (s *svc) List(ctx context.Context, q string, limit, offset int) ([]m.Patient, error) {
	return s.repo.List(ctx,q,limit,offset)
}
func (s *svc) GetByID(ctx context.Context, id string) (m.Patient, error) {
	return s.repo.GetByID(ctx,id)
}
