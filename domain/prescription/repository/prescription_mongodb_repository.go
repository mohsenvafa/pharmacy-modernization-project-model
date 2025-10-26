package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	m "pharmacy-modernization-project-model/domain/prescription/contracts/model"
	platformErrors "pharmacy-modernization-project-model/internal/platform/errors"
	"pharmacy-modernization-project-model/internal/validators/validation_logic"
)

// PrescriptionMongoRepository implements PrescriptionRepository interface using MongoDB
type PrescriptionMongoRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewPrescriptionMongoRepository creates a new MongoDB prescription repository
func NewPrescriptionMongoRepository(collection *mongo.Collection, logger *zap.Logger) PrescriptionRepository {
	return &PrescriptionMongoRepository{
		collection: collection,
		logger:     logger,
	}
}

// handleError processes errors and converts them to appropriate repository errors
func (r *PrescriptionMongoRepository) handleError(operation string, err error) error {
	if err == nil {
		return nil
	}

	r.logger.Error("MongoDB operation failed",
		zap.String("operation", operation),
		zap.Error(err))

	// Use the shared error handling
	return platformErrors.HandleMongoError(operation, err)
}

// List retrieves prescriptions with pagination and optional status filter
func (r *PrescriptionMongoRepository) List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB List operation completed",
			zap.String("status", status),
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Duration("duration", time.Since(start)))
	}()

	// Build filter with input validation
	filter := bson.M{}
	if status != "" {
		// Validate status input to prevent NoSQL injection
		if err := validation_logic.ValidateOneOf("status", status, "Draft", "Active", "Paused", "Completed"); err != nil {
			r.logger.Warn("Invalid status provided",
				zap.String("status", status),
				zap.Error(err))
			return nil, platformErrors.NewValidationError("status", status, "Invalid status value")
		}
		filter["status"] = status
	}

	// Configure options
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, r.handleError("List", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var prescriptions []m.Prescription
	if err := cursor.All(ctx, &prescriptions); err != nil {
		return nil, r.handleError("List", err)
	}

	r.logger.Debug("Successfully retrieved prescriptions from MongoDB",
		zap.Int("count", len(prescriptions)))

	return prescriptions, nil
}

// GetByID retrieves a prescription by ID
func (r *PrescriptionMongoRepository) GetByID(ctx context.Context, id string) (m.Prescription, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB GetByID operation completed",
			zap.String("id", id),
			zap.Duration("duration", time.Since(start)))
	}()

	// Validate input to prevent NoSQL injection
	if err := validation_logic.ValidateID("id", id); err != nil {
		r.logger.Warn("Invalid prescription ID provided",
			zap.String("id", id),
			zap.Error(err))
		return m.Prescription{}, platformErrors.NewValidationError("id", id, "Invalid prescription ID format")
	}

	filter := bson.M{"_id": id}

	var prescription m.Prescription
	err := r.collection.FindOne(ctx, filter).Decode(&prescription)
	if err != nil {
		return m.Prescription{}, r.handleError("GetByID", err)
	}

	r.logger.Debug("Successfully retrieved prescription from MongoDB",
		zap.String("id", id))

	return prescription, nil
}

// Create creates a new prescription
func (r *PrescriptionMongoRepository) Create(ctx context.Context, p m.Prescription) (m.Prescription, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB Create operation completed",
			zap.String("id", p.ID),
			zap.Duration("duration", time.Since(start)))
	}()

	// Set creation timestamp if not set
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	// Insert document
	_, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return m.Prescription{}, r.handleError("Create", err)
	}

	r.logger.Info("Successfully created prescription in MongoDB",
		zap.String("id", p.ID),
		zap.String("patient_id", p.PatientID))

	return p, nil
}

// Update updates an existing prescription
func (r *PrescriptionMongoRepository) Update(ctx context.Context, id string, p m.Prescription) (m.Prescription, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB Update operation completed",
			zap.String("id", id),
			zap.Duration("duration", time.Since(start)))
	}()

	// Validate input to prevent NoSQL injection
	if err := validation_logic.ValidateID("id", id); err != nil {
		r.logger.Warn("Invalid prescription ID provided for update",
			zap.String("id", id),
			zap.Error(err))
		return m.Prescription{}, platformErrors.NewValidationError("id", id, "Invalid prescription ID format")
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"patient_id": p.PatientID,
			"drug":       p.Drug,
			"dose":       p.Dose,
			"status":     p.Status,
			"updated_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(false)
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return m.Prescription{}, r.handleError("Update", err)
	}

	if result.MatchedCount == 0 {
		return m.Prescription{}, platformErrors.NewRepositoryError(
			platformErrors.ErrorTypeNotFound,
			"Prescription not found",
			mongo.ErrNoDocuments,
		)
	}

	// Retrieve updated document
	updatedPrescription, err := r.GetByID(ctx, id)
	if err != nil {
		return m.Prescription{}, r.handleError("Update", err)
	}

	r.logger.Info("Successfully updated prescription in MongoDB",
		zap.String("id", id),
		zap.String("patient_id", p.PatientID))

	return updatedPrescription, nil
}

// ListByPatientID retrieves prescriptions for a specific patient
func (r *PrescriptionMongoRepository) ListByPatientID(ctx context.Context, patientID string) ([]m.Prescription, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB ListByPatientID operation completed",
			zap.String("patient_id", patientID),
			zap.Duration("duration", time.Since(start)))
	}()

	// Validate input to prevent NoSQL injection
	if err := validation_logic.ValidateID("patient_id", patientID); err != nil {
		r.logger.Warn("Invalid patient_id provided",
			zap.String("patient_id", patientID),
			zap.Error(err))
		return nil, platformErrors.NewValidationError("patient_id", patientID, "Invalid patient ID format")
	}

	filter := bson.M{"patient_id": patientID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, r.handleError("ListByPatientID", err)
	}
	defer cursor.Close(ctx)

	var prescriptions []m.Prescription
	if err := cursor.All(ctx, &prescriptions); err != nil {
		return nil, r.handleError("ListByPatientID", err)
	}

	r.logger.Debug("Successfully retrieved prescriptions by patient ID from MongoDB",
		zap.String("patient_id", patientID),
		zap.Int("count", len(prescriptions)))

	return prescriptions, nil
}

// CountByStatus returns the total number of prescriptions matching the status
func (r *PrescriptionMongoRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB CountByStatus operation completed",
			zap.String("status", status),
			zap.Duration("duration", time.Since(start)))
	}()

	// Build filter with input validation
	filter := bson.M{}
	if status != "" {
		// Validate status input to prevent NoSQL injection
		if err := validation_logic.ValidateOneOf("status", status, "Draft", "Active", "Paused", "Completed"); err != nil {
			r.logger.Warn("Invalid status provided for count",
				zap.String("status", status),
				zap.Error(err))
			return 0, platformErrors.NewValidationError("status", status, "Invalid status value")
		}
		filter["status"] = status
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, r.handleError("CountByStatus", err)
	}

	r.logger.Debug("Successfully counted prescriptions in MongoDB",
		zap.String("status", status),
		zap.Int64("count", count))

	return int(count), nil
}

// HealthCheck performs a health check on the repository
func (r *PrescriptionMongoRepository) HealthCheck(ctx context.Context) error {
	// Try to count documents as a simple health check
	_, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return r.handleError("HealthCheck", err)
	}
	return nil
}

// CreateIndexes creates recommended indexes for optimal performance
func (r *PrescriptionMongoRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "patient_id", Value: 1}},
			Options: options.Index().
				SetName("patient_id_1").
				SetBackground(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
			Options: options.Index().
				SetName("status_1").
				SetBackground(true),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().
				SetName("created_at_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().
				SetName("_id_1").
				SetUnique(true).
				SetBackground(true),
		},
		{
			Keys: bson.D{{Key: "patient_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().
				SetName("patient_id_1_status_1").
				SetBackground(true),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return r.handleError("CreateIndexes", err)
	}

	r.logger.Info("Successfully created MongoDB indexes for prescriptions collection")
	return nil
}
