package builder

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	patientrepo "pharmacy-modernization-project-model/domain/patient/repository"
)

// CreatePatientRepository creates the appropriate patient repository based on dependencies
func CreatePatientRepository(logger *zap.Logger, mongoCollection *mongo.Collection) patientrepo.PatientRepository {
	// Use MongoDB repository if collection is provided, otherwise fallback to memory
	if mongoCollection != nil {
		return patientrepo.NewPatientMongoRepository(mongoCollection, logger)
	}

	return patientrepo.NewPatientMemoryRepository()
}
