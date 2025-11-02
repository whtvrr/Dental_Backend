package repositories

import (
	"context"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnamnesisRepository interface {
	Create(ctx context.Context, anamnesis *entities.Anamnesis) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Anamnesis, error)
	Update(ctx context.Context, anamnesis *entities.Anamnesis) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, offset, limit int) ([]*entities.Anamnesis, error)
}
