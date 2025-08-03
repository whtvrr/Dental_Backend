package repositories

import (
	"context"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FormulaRepository interface {
	Create(ctx context.Context, formula *entities.Formula) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Formula, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.Formula, error)
	Update(ctx context.Context, formula *entities.Formula) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	UpdateToothStatus(ctx context.Context, formulaID primitive.ObjectID, toothNumber int, part string, status *entities.ToothStatus) error
}