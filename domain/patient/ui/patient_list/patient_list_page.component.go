package patient_list

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"

	patientmodel "pharmacy-modernization-project-model/domain/patient/contracts/model"
	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/domain/patient/ui/paths"
)

type PatientListComponent struct {
	patientsService patSvc.PatientService
	log             *zap.Logger
}

func NewPatientListComponent(deps *contracts.UiDependencies) *PatientListComponent {
	return &PatientListComponent{patientsService: deps.PatientSvc, log: deps.Log}
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

func (c *PatientListComponent) Handler(w http.ResponseWriter, r *http.Request) {
	patients, err := c.patientsService.List(r.Context(), "", 1000, 0)
	if err != nil {
		if c.log != nil {
			c.log.Error("failed to load patients", zap.Error(err))
		}
		http.Error(w, "failed to load patients", http.StatusInternalServerError)
		return
	}

	pageParam := r.URL.Query().Get("page")
	pageNum, err := strconv.Atoi(pageParam)
	if err != nil {
		pageNum = 1
	}
	patientsPage, totalPages, currentPage := paginatePatients(patients, pageNum)

	view := PatientListPageComponentView(PatientListPageParam{
		Patients:    patientsPage,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		ListPath:    paths.PatientListURL(),
		DetailPath:  paths.PatientDetailURL,
	})

	if err := view.Render(r.Context(), w); err != nil {
		if c.log != nil {
			c.log.Error("failed to render patient list", zap.Error(err))
		}
		http.Error(w, "failed to render patient list", http.StatusInternalServerError)
		return
	}
}
