package repository

import (
	"context"
	"time"

	m "github.com/pharmacy-modernization-project-model/internal/domain/patient/model"
)

type memRepo struct{ items map[string]m.Patient }

func NewMemRepo() Repository {
	r := &memRepo{items: map[string]m.Patient{}}
	sample := []struct {
		id    string
		name  string
		phone string
		state string
		dob   time.Time
	}{
		{"P001", "Ava Thompson", "(206) 417-8842", "Washington", time.Date(1988, time.January, 12, 0, 0, 0, 0, time.UTC)},
		{"P002", "Liam Anderson", "(415) 736-5528", "California", time.Date(1979, time.March, 3, 0, 0, 0, 0, time.UTC)},
		{"P003", "Sophia Martinez", "(617) 980-3314", "Massachusetts", time.Date(1992, time.July, 27, 0, 0, 0, 0, time.UTC)},
		{"P004", "Noah Patel", "(972) 645-2091", "Texas", time.Date(1985, time.May, 5, 0, 0, 0, 0, time.UTC)},
		{"P005", "Mia Chen", "(312) 478-6605", "Illinois", time.Date(1996, time.September, 19, 0, 0, 0, 0, time.UTC)},
		{"P006", "Ethan Johnson", "(303) 825-1947", "Colorado", time.Date(1975, time.November, 8, 0, 0, 0, 0, time.UTC)},
		{"P007", "Olivia Rossi", "(646) 291-0743", "New York", time.Date(1990, time.February, 22, 0, 0, 0, 0, time.UTC)},
		{"P008", "Jackson Lee", "(503) 913-2286", "Oregon", time.Date(1983, time.April, 16, 0, 0, 0, 0, time.UTC)},
		{"P009", "Emma Davis", "(305) 744-1189", "Florida", time.Date(1998, time.December, 2, 0, 0, 0, 0, time.UTC)},
		{"P010", "Lucas Hernandez", "(713) 402-5378", "Texas", time.Date(1981, time.June, 14, 0, 0, 0, 0, time.UTC)},
	}

	for _, s := range sample {
		r.items[s.id] = m.Patient{
			ID:        s.id,
			Name:      s.name,
			Phone:     s.phone,
			State:     s.state,
			DOB:       s.dob,
			CreatedAt: time.Now(),
		}
	}

	return r
}

func (r *memRepo) List(ctx context.Context, query string, limit, offset int) ([]m.Patient, error) {
	res := []m.Patient{}
	for _, v := range r.items {
		res = append(res, v)
	}
	if offset >= len(res) {
		return []m.Patient{}, nil
	}
	end := offset + limit
	if end > len(res) {
		end = len(res)
	}
	return res[offset:end], nil
}
func (r *memRepo) GetByID(ctx context.Context, id string) (m.Patient, error) { return r.items[id], nil }
func (r *memRepo) Create(ctx context.Context, p m.Patient) (m.Patient, error) {
	r.items[p.ID] = p
	return p, nil
}
func (r *memRepo) Update(ctx context.Context, id string, p m.Patient) (m.Patient, error) {
	r.items[id] = p
	return p, nil
}
