package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	addressModel "pharmacy-modernization-project-model/domain/patient/contracts/model"
	platformErrors "pharmacy-modernization-project-model/internal/platform/errors"
	"pharmacy-modernization-project-model/internal/validators/validation_logic"
)

// AddressMongoRepository implements AddressRepository interface using MongoDB
type AddressMongoRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewAddressMongoRepository creates a new MongoDB address repository
func NewAddressMongoRepository(collection *mongo.Collection, logger *zap.Logger) AddressRepository {
	return &AddressMongoRepository{
		collection: collection,
		logger:     logger,
	}
}

// handleError processes MongoDB errors and converts them to appropriate repository errors
func (r *AddressMongoRepository) handleError(operation string, err error) error {
	if err == nil {
		return nil
	}

	r.logger.Error("MongoDB operation failed",
		zap.String("operation", operation),
		zap.Error(err))

	return platformErrors.HandleMongoError(operation, err)
}

// ListByPatientID retrieves all addresses for a specific patient
func (r *AddressMongoRepository) ListByPatientID(ctx context.Context, patientID string) ([]addressModel.Address, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB ListByPatientID operation completed",
			zap.String("patient_id", patientID),
			zap.Duration("duration", time.Since(start)))
	}()

	// Validate input to prevent NoSQL injection
	if err := validation_logic.ValidateID("patient_id", patientID); err != nil {
		r.logger.Warn("Invalid patient_id provided",
			zap.Error(err))
		return nil, platformErrors.NewValidationError("patient_id", patientID, "Invalid patient ID format")
	}

	filter := bson.M{"patient_id": patientID}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, r.handleError("ListByPatientID", err)
	}
	defer cursor.Close(ctx)

	var addresses []addressModel.Address
	if err := cursor.All(ctx, &addresses); err != nil {
		return nil, r.handleError("ListByPatientID", err)
	}

	r.logger.Debug("Successfully retrieved addresses from MongoDB",
		zap.String("patient_id", patientID),
		zap.Int("count", len(addresses)))

	return addresses, nil
}

// GetByID retrieves a specific address by ID and patient ID
func (r *AddressMongoRepository) GetByID(ctx context.Context, patientID, addressID string) (addressModel.Address, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB GetByID operation completed",
			zap.String("patient_id", patientID),
			zap.String("address_id", addressID),
			zap.Duration("duration", time.Since(start)))
	}()

	// Validate input to prevent NoSQL injection
	if err := validation_logic.ValidateID("patient_id", patientID); err != nil {
		r.logger.Warn("Invalid patient_id provided",
			zap.Error(err))
		return addressModel.Address{}, platformErrors.NewValidationError("patient_id", patientID, "Invalid patient ID format")
	}
	if err := validation_logic.ValidateID("address_id", addressID); err != nil {
		r.logger.Warn("Invalid address_id provided",
			zap.Error(err))
		return addressModel.Address{}, platformErrors.NewValidationError("address_id", addressID, "Invalid address ID format")
	}

	filter := bson.M{
		"_id":        addressID,
		"patient_id": patientID,
	}

	var address addressModel.Address
	err := r.collection.FindOne(ctx, filter).Decode(&address)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return addressModel.Address{}, nil // Return empty address if not found (matches memory repo behavior)
		}
		return addressModel.Address{}, r.handleError("GetByID", err)
	}

	r.logger.Debug("Successfully retrieved address from MongoDB",
		zap.String("patient_id", patientID),
		zap.String("address_id", addressID))

	return address, nil
}

// Upsert creates or updates an address
func (r *AddressMongoRepository) Upsert(ctx context.Context, patientID string, address addressModel.Address) (addressModel.Address, error) {
	start := time.Now()
	defer func() {
		r.logger.Debug("MongoDB Upsert operation completed",
			zap.String("patient_id", patientID),
			zap.String("address_id", address.ID),
			zap.Duration("duration", time.Since(start)))
	}()

	// Validate input to prevent NoSQL injection
	if err := validation_logic.ValidateID("patient_id", patientID); err != nil {
		r.logger.Warn("Invalid patient_id provided",
			zap.Error(err))
		return addressModel.Address{}, platformErrors.NewValidationError("patient_id", patientID, "Invalid patient ID format")
	}
	if address.ID != "" {
		if err := validation_logic.ValidateID("address_id", address.ID); err != nil {
			r.logger.Warn("Invalid address_id provided",
				zap.Error(err))
			return addressModel.Address{}, platformErrors.NewValidationError("address_id", address.ID, "Invalid address ID format")
		}
	}

	// Generate ID if not provided
	if address.ID == "" {
		address.ID = fmt.Sprintf("%s-addr-%d", patientID, time.Now().Unix())
	}

	// Ensure patient ID is set
	address.PatientID = patientID

	filter := bson.M{
		"_id":        address.ID,
		"patient_id": patientID,
	}

	update := bson.M{
		"$set": bson.M{
			"_id":        address.ID,
			"patient_id": address.PatientID,
			"line1":      address.Line1,
			"line2":      address.Line2,
			"city":       address.City,
			"state":      address.State,
			"zip":        address.Zip,
			"updated_at": time.Now(),
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return addressModel.Address{}, r.handleError("Upsert", err)
	}

	r.logger.Info("Successfully upserted address in MongoDB",
		zap.String("patient_id", patientID),
		zap.String("address_id", address.ID))

	return address, nil
}

// CreateIndexes creates recommended indexes for optimal performance
func (r *AddressMongoRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "patient_id", Value: 1}},
			Options: options.Index().
				SetName("patient_id_1").
				SetBackground(true),
		},
		{
			Keys: bson.D{{Key: "_id", Value: 1}, {Key: "patient_id", Value: 1}},
			Options: options.Index().
				SetName("_id_1_patient_id_1").
				SetUnique(true).
				SetBackground(true),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return r.handleError("CreateIndexes", err)
	}

	r.logger.Info("Successfully created MongoDB indexes for addresses collection")
	return nil
}
