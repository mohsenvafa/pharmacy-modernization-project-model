package repository

import (
	"context"
	"fmt"
	"time"
	m "github.com/pharmacy-modernization-project-model/internal/domain/patient/model"
)

type memRepo struct { items map[string]m.Patient }

func NewMemRepo() Repository {
	r := &memRepo{items: map[string]m.Patient{}}
	for i:=1;i<=30;i++ {
		id := fmt.Sprintf("P%03d", i)
		r.items[id] = m.Patient{ID:id, Name: fmt.Sprintf("Patient %d", i), Phone: "555-0101", DOB: time.Date(1980+i%20, time.Month((i%12)+1), (i%28)+1,0,0,0,0,time.UTC), CreatedAt: time.Now()}
	}
	return r
}

func (r *memRepo) List(ctx context.Context, query string, limit, offset int) ([]m.Patient, error) {
	res := []m.Patient{}
	for _, v := range r.items { res = append(res, v) }
	if offset >= len(res) { return []m.Patient{}, nil }
	end := offset+limit
	if end > len(res) { end = len(res) }
	return res[offset:end], nil
}
func (r *memRepo) GetByID(ctx context.Context, id string) (m.Patient, error) { return r.items[id], nil }
func (r *memRepo) Create(ctx context.Context, p m.Patient) (m.Patient, error) { r.items[p.ID]=p; return p, nil }
func (r *memRepo) Update(ctx context.Context, id string, p m.Patient) (m.Patient, error) { r.items[id]=p; return p, nil }
