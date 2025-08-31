package mongodb

import (
	"context"
	"time"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type statusRepository struct {
	collection *mongo.Collection
}

func NewStatusRepository(db *mongo.Database) repositories.StatusRepository {
	return &statusRepository{
		collection: db.Collection("statuses"),
	}
}

func (r *statusRepository) Create(ctx context.Context, status *entities.Status) error {
	status.CreatedAt = time.Now()
	status.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, status)
	if err != nil {
		return err
	}
	
	status.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *statusRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Status, error) {
	var status entities.Status
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *statusRepository) Update(ctx context.Context, status *entities.Status) error {
	status.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": status.ID},
		bson.M{"$set": status},
	)
	return err
}

func (r *statusRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *statusRepository) List(ctx context.Context, offset, limit int) ([]*entities.Status, error) {
	opts := options.Find()
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "type", Value: 1}, {Key: "title", Value: 1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var statuses []*entities.Status
	for cursor.Next(ctx) {
		var status entities.Status
		if err := cursor.Decode(&status); err != nil {
			return nil, err
		}
		statuses = append(statuses, &status)
	}
	
	return statuses, cursor.Err()
}

func (r *statusRepository) GetByType(ctx context.Context, statusType entities.StatusType) ([]*entities.Status, error) {
	opts := options.Find().SetSort(bson.D{{Key: "title", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"type": statusType}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var statuses []*entities.Status
	for cursor.Next(ctx) {
		var status entities.Status
		if err := cursor.Decode(&status); err != nil {
			return nil, err
		}
		statuses = append(statuses, &status)
	}
	
	return statuses, cursor.Err()
}

func (r *statusRepository) GetActiveByType(ctx context.Context, statusType entities.StatusType) ([]*entities.Status, error) {
	filter := bson.M{
		"type":      statusType,
		"is_active": true,
	}
	
	opts := options.Find().SetSort(bson.D{{Key: "title", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var statuses []*entities.Status
	for cursor.Next(ctx) {
		var status entities.Status
		if err := cursor.Decode(&status); err != nil {
			return nil, err
		}
		statuses = append(statuses, &status)
	}
	
	return statuses, cursor.Err()
}

func (r *statusRepository) GetActiveByTypeWithPagination(ctx context.Context, statusType entities.StatusType, offset, limit int) ([]*entities.Status, error) {
	filter := bson.M{
		"type":      statusType,
		"is_active": true,
	}
	
	opts := options.Find()
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "title", Value: 1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var statuses []*entities.Status
	for cursor.Next(ctx) {
		var status entities.Status
		if err := cursor.Decode(&status); err != nil {
			return nil, err
		}
		statuses = append(statuses, &status)
	}
	
	return statuses, cursor.Err()
}