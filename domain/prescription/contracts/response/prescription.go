package response

import (
	"time"

	model "pharmacy-modernization-project-model/domain/prescription/contracts/model"
)

// PrescriptionResponse is the transport representation returned by the API.
type PrescriptionResponse struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patient_id"`
	Drug      string    `json:"drug"`
	Dose      string    `json:"dose"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func FromModel(m model.Prescription) PrescriptionResponse {
	return PrescriptionResponse{
		ID:        m.ID,
		PatientID: m.PatientID,
		Drug:      m.Drug,
		Dose:      m.Dose,
		Status:    string(m.Status),
		CreatedAt: m.CreatedAt,
	}
}

func FromModels(items []model.Prescription) []PrescriptionResponse {
	out := make([]PrescriptionResponse, 0, len(items))
	for _, item := range items {
		out = append(out, FromModel(item))
	}
	return out
}
