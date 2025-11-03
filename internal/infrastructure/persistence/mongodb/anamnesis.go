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

type anamnesisRepository struct {
	collection *mongo.Collection
}

func NewAnamnesisRepository(db *mongo.Database) repositories.AnamnesisRepository {
	return &anamnesisRepository{
		collection: db.Collection("anamnesis"),
	}
}

func (r *anamnesisRepository) Create(ctx context.Context, anamnesis *entities.Anamnesis) error {
	anamnesis.CreatedAt = time.Now()
	anamnesis.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, anamnesis)
	if err != nil {
		return err
	}

	anamnesis.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *anamnesisRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Anamnesis, error) {
	var anamnesis entities.Anamnesis
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&anamnesis)
	if err != nil {
		return nil, err
	}
	return &anamnesis, nil
}

func (r *anamnesisRepository) Update(ctx context.Context, anamnesis *entities.Anamnesis) error {
	anamnesis.UpdatedAt = time.Now()

	// Only update the fields that should change, preserve created_at
	update := bson.M{
		"$set": bson.M{
			"text":       anamnesis.Text,
			"updated_at": anamnesis.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": anamnesis.ID},
		update,
	)
	return err
}

func (r *anamnesisRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *anamnesisRepository) List(ctx context.Context, offset, limit int) ([]*entities.Anamnesis, error) {
	opts := options.Find()
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var anamnesises []*entities.Anamnesis
	for cursor.Next(ctx) {
		var anamnesis entities.Anamnesis
		if err := cursor.Decode(&anamnesis); err != nil {
			return nil, err
		}
		anamnesises = append(anamnesises, &anamnesis)
	}

	return anamnesises, cursor.Err()
}
