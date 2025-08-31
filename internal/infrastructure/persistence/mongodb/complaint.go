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

type complaintRepository struct {
	collection *mongo.Collection
}

func NewComplaintRepository(db *mongo.Database) repositories.ComplaintRepository {
	return &complaintRepository{
		collection: db.Collection("complaints"),
	}
}

func (r *complaintRepository) Create(ctx context.Context, complaint *entities.Complaint) error {
	complaint.CreatedAt = time.Now()
	complaint.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, complaint)
	if err != nil {
		return err
	}
	
	complaint.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *complaintRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Complaint, error) {
	var complaint entities.Complaint
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&complaint)
	if err != nil {
		return nil, err
	}
	return &complaint, nil
}

func (r *complaintRepository) Update(ctx context.Context, complaint *entities.Complaint) error {
	complaint.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": complaint.ID},
		bson.M{"$set": complaint},
	)
	return err
}

func (r *complaintRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *complaintRepository) List(ctx context.Context, offset, limit int) ([]*entities.Complaint, error) {
	opts := options.Find()
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "title", Value: 1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var complaints []*entities.Complaint
	for cursor.Next(ctx) {
		var complaint entities.Complaint
		if err := cursor.Decode(&complaint); err != nil {
			return nil, err
		}
		complaints = append(complaints, &complaint)
	}
	
	return complaints, cursor.Err()
}

func (r *complaintRepository) GetByCategory(ctx context.Context, category string) ([]*entities.Complaint, error) {
	opts := options.Find().SetSort(bson.D{{Key: "title", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"category": category}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var complaints []*entities.Complaint
	for cursor.Next(ctx) {
		var complaint entities.Complaint
		if err := cursor.Decode(&complaint); err != nil {
			return nil, err
		}
		complaints = append(complaints, &complaint)
	}
	
	return complaints, cursor.Err()
}