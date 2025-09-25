package repository

import (
	"context"
	"fmt"
	m "pharmacy-modernization-project-model/domain/prescription/contracts/model"
	"time"
)

type PrescriptionMemoryRepository struct{ items map[string]m.Prescription }

func NewPrescriptionMemoryRepository() PrescriptionRepository {
	r := &PrescriptionMemoryRepository{items: map[string]m.Prescription{}}
	statuses := []m.Status{m.Draft, m.Active, m.Paused, m.Completed}
	for i := 1; i <= 50; i++ {
		id := fmt.Sprintf("R%03d", i)
		r.items[id] = m.Prescription{ID: id, PatientID: fmt.Sprintf("P%03d", (i%30)+1), Drug: "Amoxicillin", Dose: "500mg", Status: statuses[i%len(statuses)], CreatedAt: time.Now().AddDate(0, 0, -i)}
	}
	return r
}

func (r *PrescriptionMemoryRepository) List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error) {
	res := []m.Prescription{}
	for _, v := range r.items {
		if status == "" || string(v.Status) == status {
			res = append(res, v)
		}
	}
	if offset >= len(res) {
		return []m.Prescription{}, nil
	}
	end := offset + limit
	if end > len(res) {
		end = len(res)
	}
	return res[offset:end], nil
}
func (r *PrescriptionMemoryRepository) GetByID(ctx context.Context, id string) (m.Prescription, error) {
	return r.items[id], nil
}
func (r *PrescriptionMemoryRepository) Create(ctx context.Context, p m.Prescription) (m.Prescription, error) {
	r.items[p.ID] = p
	return p, nil
}
func (r *PrescriptionMemoryRepository) Update(ctx context.Context, id string, p m.Prescription) (m.Prescription, error) {
	r.items[id] = p
	return p, nil
}

func (r *PrescriptionMemoryRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	if status == "" {
		return len(r.items), nil
	}
	count := 0
	for _, v := range r.items {
		if string(v.Status) == status {
			count++
		}
	}
	return count, nil
}
