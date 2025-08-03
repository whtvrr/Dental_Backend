package repositories

import (
	"context"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComplaintRepository interface {
	Create(ctx context.Context, complaint *entities.Complaint) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Complaint, error)
	Update(ctx context.Context, complaint *entities.Complaint) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, offset, limit int) ([]*entities.Complaint, error)
	GetActive(ctx context.Context) ([]*entities.Complaint, error)
	GetByCategory(ctx context.Context, category string) ([]*entities.Complaint, error)
}