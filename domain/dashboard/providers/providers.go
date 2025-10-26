package providers

import (
	"context"

	"pharmacy-modernization-project-model/domain/patient/contracts/request"
)

type PatientStatsProvider interface {
	Count(ctx context.Context, req request.PatientListQueryRequest) (int, error)
}

type PrescriptionStatsProvider interface {
	CountByStatus(ctx context.Context, status string) (int, error)
}
