package repositories

import (
	"context"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusRepository interface {
	Create(ctx context.Context, status *entities.Status) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Status, error)
	Update(ctx context.Context, status *entities.Status) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, offset, limit int) ([]*entities.Status, error)
	GetByType(ctx context.Context, statusType entities.StatusType) ([]*entities.Status, error)
	GetActiveByType(ctx context.Context, statusType entities.StatusType) ([]*entities.Status, error)
}