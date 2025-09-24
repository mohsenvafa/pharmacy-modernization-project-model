package patient_list

import (
	"net/http"
	"strconv"

	patientmodel "pharmacy-modernization-project-model/internal/domain/patient/contracts/model"
	patSvc "pharmacy-modernization-project-model/internal/domain/patient/service"

	"go.uber.org/zap"
)

type PatientListHandler struct {
	patientsService patSvc.PatientService
	log             *zap.Logger
}

func NewPatientListHandler(patients patSvc.PatientService, log *zap.Logger) *PatientListHandler {
	return &PatientListHandler{patientsService: patients, log: log}
}

const pageSize = 5

func paginatePatients(pats []patientmodel.Patient, page int) ([]patientmodel.Patient, int, int) {
	if page < 1 {
		page = 1
	}
	if len(pats) == 0 {
		return []patientmodel.Patient{}, 1, 1
	}
	totalPages := (len(pats) + pageSize - 1) / pageSize
	if page > totalPages {
		page = totalPages
	}
	start := (page - 1) * pageSize
	if start >= len(pats) {
		return []patientmodel.Patient{}, totalPages, page
	}
	end := start + pageSize
	if end > len(pats) {
		end = len(pats)
	}
	return pats[start:end], totalPages, page
}

func (u *PatientListHandler) Handler(w http.ResponseWriter, r *http.Request) {
	patients, err := u.patientsService.List(r.Context(), "", 1000, 0)
	if err != nil {
		http.Error(w, "failed to load patients", http.StatusInternalServerError)
		return
	}

	pageParam := r.URL.Query().Get("page")
	pageNum, err := strconv.Atoi(pageParam)
	if err != nil {
		pageNum = 1
	}
	currentPage := pageNum
	patientsPage, totalPages, currentPage := paginatePatients(patients, currentPage)
	page := PatientListPage(PatientListPageParam{
		Patients:    patientsPage,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	})
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render patient list", http.StatusInternalServerError)
		return
	}
}
