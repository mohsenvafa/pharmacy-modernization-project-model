package providers

import "context"

type PatientStatsProvider interface {
	Count(ctx context.Context, query string) (int, error)
}

type PrescriptionStatsProvider interface {
	CountByStatus(ctx context.Context, status string) (int, error)
}
