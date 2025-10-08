package builder

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	prescriptionrepo "pharmacy-modernization-project-model/domain/prescription/repository"
)

// CreatePrescriptionRepository creates the appropriate prescription repository based on dependencies
func CreatePrescriptionRepository(logger *zap.Logger, mongoCollection *mongo.Collection) prescriptionrepo.PrescriptionRepository {
	// Use MongoDB repository if collection is provided, otherwise fallback to memory
	if mongoCollection != nil {
		return prescriptionrepo.NewPrescriptionMongoRepository(mongoCollection, logger)
	}

	return prescriptionrepo.NewPrescriptionMemoryRepository()
}
