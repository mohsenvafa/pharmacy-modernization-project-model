package service

import (
	"context"
	repo "github.com/pharmacy-modernization-project-model/internal/domain/prescription/repository"
	m "github.com/pharmacy-modernization-project-model/internal/domain/prescription/model"
	"go.uber.org/zap"
)

type Service interface {
	List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error)
	GetByID(ctx context.Context, id string) (m.Prescription, error)
}

type svc struct { repo repo.Repository; log *zap.Logger }

func New(r repo.Repository, l *zap.Logger) Service { return &svc{repo:r, log:l} }

func (s *svc) List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error) {
	return s.repo.List(ctx,status,limit,offset)
}
func (s *svc) GetByID(ctx context.Context, id string) (m.Prescription, error) {
	return s.repo.GetByID(ctx,id)
}
