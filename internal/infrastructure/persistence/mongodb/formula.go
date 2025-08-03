package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type formulaRepository struct {
	collection *mongo.Collection
}

func NewFormulaRepository(db *mongo.Database) repositories.FormulaRepository {
	return &formulaRepository{
		collection: db.Collection("formulas"),
	}
}

func (r *formulaRepository) Create(ctx context.Context, formula *entities.Formula) error {
	formula.CreatedAt = time.Now()
	formula.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, formula)
	if err != nil {
		return err
	}
	
	formula.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *formulaRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Formula, error) {
	var formula entities.Formula
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&formula)
	if err != nil {
		return nil, err
	}
	return &formula, nil
}

func (r *formulaRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.Formula, error) {
	var formula entities.Formula
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&formula)
	if err != nil {
		return nil, err
	}
	return &formula, nil
}

func (r *formulaRepository) Update(ctx context.Context, formula *entities.Formula) error {
	formula.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": formula.ID},
		bson.M{"$set": formula},
	)
	return err
}

func (r *formulaRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *formulaRepository) UpdateToothStatus(ctx context.Context, formulaID primitive.ObjectID, toothNumber int, part string, status *entities.ToothStatus) error {
	var updateField string
	switch part {
	case "whole":
		updateField = fmt.Sprintf("teeth.%d.whole", toothNumber-1)
	case "gum":
		updateField = fmt.Sprintf("teeth.%d.gum", toothNumber-1)
	case "roots":
		updateField = fmt.Sprintf("teeth.%d.roots", toothNumber-1)
	default:
		updateField = fmt.Sprintf("teeth.%d.segments.%s", toothNumber-1, part)
	}
	
	update := bson.M{
		"$set": bson.M{
			updateField:   status,
			"updated_at": time.Now(),
		},
	}
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": formulaID},
		update,
	)
	return err
}