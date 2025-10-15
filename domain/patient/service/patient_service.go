package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	m "pharmacy-modernization-project-model/domain/patient/contracts/model"
	repo "pharmacy-modernization-project-model/domain/patient/repository"

	"go.uber.org/zap"
)

type PatientService interface {
	List(ctx context.Context, q string, limit, offset int) ([]m.Patient, error)
	GetByID(ctx context.Context, id string) (m.Patient, error)
	Count(ctx context.Context, q string) (int, error)
}

type patientSvc struct {
	repo  repo.PatientRepository
	cache Cache
	log   *zap.Logger
}

// Cache interface for patient service
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}

func New(r repo.PatientRepository, cache Cache, l *zap.Logger) PatientService {
	return &patientSvc{repo: r, cache: cache, log: l}
}

func (s *patientSvc) List(ctx context.Context, q string, limit, offset int) ([]m.Patient, error) {
	return s.repo.List(ctx, q, limit, offset)
}
func (s *patientSvc) GetByID(ctx context.Context, id string) (m.Patient, error) {
	cacheKey := fmt.Sprintf("patient:id:%s", id)

	// Try cache first
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var patient m.Patient
			if err := json.Unmarshal(cached, &patient); err == nil {
				s.log.Debug("Patient retrieved from cache", zap.String("patient_id", id))
				return patient, nil
			}
		}
	}

	s.log.Info("Getting patient from repository", zap.String("patient_id", id))

	patient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get patient",
			zap.String("patient_id", id),
			zap.Error(err))
		return m.Patient{}, err
	}

	// Cache the result
	if s.cache != nil {
		if data, err := json.Marshal(patient); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, 30*time.Minute); err != nil {
				s.log.Warn("Failed to cache patient", zap.String("patient_id", id), zap.Error(err))
			}
		}
	}

	s.log.Info("Patient retrieved successfully", zap.String("patient_id", id))
	return patient, nil
}

func (s *patientSvc) Count(ctx context.Context, q string) (int, error) {
	cacheKey := fmt.Sprintf("patient:count:%s", q)

	// Try cache first
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var count int
			if err := json.Unmarshal(cached, &count); err == nil {
				s.log.Debug("Patient count retrieved from cache", zap.String("query", q))
				return count, nil
			}
		}
	}

	count, err := s.repo.Count(ctx, q)
	if err != nil {
		return 0, err
	}

	// Cache the result with shorter TTL for counts
	if s.cache != nil {
		if data, err := json.Marshal(count); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, 5*time.Minute); err != nil {
				s.log.Warn("Failed to cache patient count", zap.String("query", q), zap.Error(err))
			}
		}
	}

	return count, nil
}
