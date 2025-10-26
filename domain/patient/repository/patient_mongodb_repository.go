package repository

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	m "pharmacy-modernization-project-model/domain/patient/contracts/model"
	"pharmacy-modernization-project-model/domain/patient/contracts/request"
	patientErrors "pharmacy-modernization-project-model/domain/patient/errors"
	platformErrors "pharmacy-modernization-project-model/internal/platform/errors"
	"pharmacy-modernization-project-model/internal/validators/validation_logic"
)

// PatientMongoRepository implements PatientRepository interface using MongoDB
type PatientMongoRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

type patientID string
type patientPhone string
type patientState string

// escapeRegexChars escapes special regex characters in a string to prevent regex injection
func escapeRegexChars(input string) string {
	return regexp.QuoteMeta(input)
}

// NewPatientMongoRepository creates a new MongoDB patient repository
func NewPatientMongoRepository(collection *mongo.Collection, logger *zap.Logger) PatientRepository {
	return &PatientMongoRepository{
		collection: collection,
		logger:     logger,
	}
}

// handleError processes MongoDB errors and converts them to appropriate repository errors
func (r *PatientMongoRepository) handleError(operation string, err error) error {
	if err == nil {
		return nil
	}

	r.logger.Error("MongoDB operation failed",
		zap.String("operation", operation),
		zap.Error(err))

	// Use the shared error handling
	return platformErrors.HandleMongoError(operation, err)
}

// List retrieves patients with pagination and optional search
func (r *PatientMongoRepository) List(ctx context.Context, req request.PatientListQueryRequest) ([]m.Patient, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB List operation completed",
			zap.Int("limit", req.Limit),
			zap.Int("offset", req.Offset),
			zap.Duration("duration", time.Since(start)))
	}()

	// Build filter with input sanitization to prevent regex injection
	filter := bson.M{}
	if req.PatientName != "" {
		// Escape special regex characters to prevent regex injection
		escapedQuery := escapeRegexChars(req.PatientName)
		filter = bson.M{
			"name": bson.M{"$regex": escapedQuery, "$options": "i"},
		}
	}

	// Configure options
	opts := options.Find().
		SetLimit(int64(req.Limit)).
		SetSkip(int64(req.Offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by creation date descending

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, r.handleError("List", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var patients []m.Patient
	if err := cursor.All(ctx, &patients); err != nil {
		return nil, r.handleError("List", err)
	}

	r.logger.Debug("Successfully retrieved patients from MongoDB",
		zap.Int("count", len(patients)))

	return patients, nil
}

// GetByID retrieves a patient by ID
func (r *PatientMongoRepository) GetByID(ctx context.Context, id string) (m.Patient, error) {
	// Input validation using bind package
	if err := validation_logic.ValidateID("id", id); err != nil {
		return m.Patient{}, err
	}

	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB GetByID operation completed",
			zap.Duration("duration", time.Since(start)))
	}()

	filter := bson.M{"_id": patientID(id)}

	var patient m.Patient
	err := r.collection.FindOne(ctx, filter).Decode(&patient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return m.Patient{}, patientErrors.NewRecordNotFoundError("Patient", id)
		}
		return m.Patient{}, fmt.Errorf("failed to get patient %s: %w", id, err)
	}

	r.logger.Debug("Successfully retrieved patient from MongoDB")

	return patient, nil
}

// Create creates a new patient
func (r *PatientMongoRepository) Create(ctx context.Context, p m.Patient) (m.Patient, error) {
	// Input validation using bind package
	if err := validation_logic.ValidateID("id", p.ID); err != nil {
		return m.Patient{}, err
	}
	if err := validation_logic.ValidateRequired("name", p.Name); err != nil {
		return m.Patient{}, err
	}
	if err := validation_logic.ValidatePhone("phone", p.Phone); err != nil {
		return m.Patient{}, err
	}

	// Validate name length
	if err := validation_logic.ValidateLength("name", p.Name, 2, 100); err != nil {
		return m.Patient{}, err
	}

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
		if mongo.IsDuplicateKeyError(err) {
			return m.Patient{}, patientErrors.NewDuplicateRecordError("Patient", p.ID)
		}
		return m.Patient{}, fmt.Errorf("failed to create patient %s: %w", p.ID, err)
	}

	r.logger.Info("Successfully created patient in MongoDB")

	return p, nil
}

// Update updates an existing patient
func (r *PatientMongoRepository) Update(ctx context.Context, id string, p m.Patient) (m.Patient, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB Update operation completed",
			zap.String("id", id),
			zap.Duration("duration", time.Since(start)))
	}()

	filter := bson.M{"_id": patientID(id)}
	update := bson.M{
		"$set": bson.M{
			"name":       p.Name,
			"phone":      p.Phone,
			"state":      p.State,
			"dob":        p.DOB,
			"updated_at": time.Now(),
		},
	}

	// Add edit tracking fields if they exist
	if p.EditBy != nil {
		update["$set"].(bson.M)["edit_by"] = *p.EditBy
	}
	if p.EditTime != nil {
		update["$set"].(bson.M)["edit_time"] = *p.EditTime
	}

	opts := options.Update().SetUpsert(false)
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return m.Patient{}, r.handleError("Update", err)
	}

	if result.MatchedCount == 0 {
		return m.Patient{}, platformErrors.NewRepositoryError(
			platformErrors.ErrorTypeNotFound,
			"Patient not found",
			mongo.ErrNoDocuments,
		)
	}

	// Retrieve updated document
	updatedPatient, err := r.GetByID(ctx, id)
	if err != nil {
		return m.Patient{}, r.handleError("Update", err)
	}

	r.logger.Info("Successfully updated patient in MongoDB")

	return updatedPatient, nil
}

// Count returns the total number of patients matching the query
func (r *PatientMongoRepository) Count(ctx context.Context, req request.PatientListQueryRequest) (int, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB Count operation completed",
			zap.Duration("duration", time.Since(start)))
	}()

	// Build filter with input sanitization to prevent regex injection
	filter := bson.M{}
	if req.PatientName != "" {
		// Escape special regex characters to prevent regex injection
		escapedQuery := escapeRegexChars(req.PatientName)
		filter = bson.M{
			"name": bson.M{"$regex": escapedQuery, "$options": "i"},
		}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, r.handleError("Count", err)
	}

	r.logger.Debug("Successfully counted patients in MongoDB",
		zap.Int64("count", count))

	return int(count), nil
}

// HealthCheck performs a health check on the repository
func (r *PatientMongoRepository) HealthCheck(ctx context.Context) error {
	// Try to count documents as a simple health check
	_, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return r.handleError("HealthCheck", err)
	}
	return nil
}

// CreateIndexes creates recommended indexes for optimal performance
func (r *PatientMongoRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: "text"}},
			Options: options.Index().
				SetName("name_text"),
		},
		{
			Keys: bson.D{{Key: "state", Value: 1}},
			Options: options.Index().
				SetName("state_1"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().
				SetName("created_at_-1"),
		},
		{
			Keys: bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().
				SetName("phone_1").
				SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().
				SetName("_id_1").
				SetUnique(true),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return r.handleError("CreateIndexes", err)
	}

	r.logger.Info("Successfully created MongoDB indexes for patients collection")
	return nil
}

// BulkInsert performs bulk insert operation for better performance
func (r *PatientMongoRepository) BulkInsert(ctx context.Context, patients []m.Patient) error {
	if len(patients) == 0 {
		return nil
	}

	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB BulkInsert operation completed",
			zap.Int("count", len(patients)),
			zap.Duration("duration", time.Since(start)))
	}()

	// Convert to interface slice for bulk insert
	docs := make([]interface{}, len(patients))
	for i, patient := range patients {
		if patient.CreatedAt.IsZero() {
			patient.CreatedAt = time.Now()
		}
		docs[i] = patient
	}

	opts := options.InsertMany().SetOrdered(false)
	result, err := r.collection.InsertMany(ctx, docs, opts)
	if err != nil {
		return r.handleError("BulkInsert", err)
	}

	r.logger.Info("Successfully bulk inserted patients into MongoDB",
		zap.Int("inserted_count", len(result.InsertedIDs)))

	return nil
}

// FindByState retrieves patients by state with pagination
func (r *PatientMongoRepository) FindByState(ctx context.Context, state string, limit, offset int) ([]m.Patient, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB FindByState operation completed",
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Duration("duration", time.Since(start)))
	}()

	// Validate state input to prevent NoSQL injection
	if state != "" {
		if err := validation_logic.ValidateLength("state", state, 2, 50); err != nil {
			r.logger.Warn("Invalid state provided",
				zap.String("state", state),
				zap.Error(err))
			return nil, platformErrors.NewValidationError("state", state, "Invalid state value")
		}
	}

	filter := bson.M{"state": patientState(state)}
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, r.handleError("FindByState", err)
	}
	defer cursor.Close(ctx)

	var patients []m.Patient
	if err := cursor.All(ctx, &patients); err != nil {
		return nil, r.handleError("FindByState", err)
	}

	r.logger.Debug("Successfully retrieved patients by state from MongoDB",
		zap.Int("count", len(patients)))

	return patients, nil
}
