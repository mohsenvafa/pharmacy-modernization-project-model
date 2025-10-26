package patient_list

import (
	"net/http"

	"go.uber.org/zap"

	patientmodel "pharmacy-modernization-project-model/domain/patient/contracts/model"
	"pharmacy-modernization-project-model/domain/patient/contracts/request"
	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/domain/patient/ui/paths"
	"pharmacy-modernization-project-model/internal/bind"
	helper "pharmacy-modernization-project-model/internal/helper"
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
	// Bind and validate query parameters for pagination
	pageReq, _, err := bind.Query[request.PatientListPageRequest](r)
	if err != nil {
		c.log.Error("failed to bind page parameters", zap.Error(err))
		helper.WriteUIError(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}

	// Set default page if not provided or invalid
	pageNum := pageReq.Page
	if pageNum < 1 {
		pageNum = 1
	}

	// Get all patients for pagination
	req := request.PatientListQueryRequest{
		Limit:  1000,
		Offset: 0,
	}
	patients, err := c.patientsService.List(r.Context(), req)
	if err != nil {
		if c.log != nil {
			c.log.Error("failed to load patients", zap.Error(err))
		}
		helper.WriteUIInternalError(w, "Failed to load patients")
		return
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
		helper.WriteUIInternalError(w, "Failed to render patient list")
		return
	}
}
