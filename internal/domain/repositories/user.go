package repositories

import (
	"context"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, offset, limit int) ([]*entities.User, error)
	GetByRole(ctx context.Context, role entities.UserRole) ([]*entities.User, error)
	GetByRoleWithPagination(ctx context.Context, role entities.UserRole, offset, limit int, query string) ([]*entities.User, error)
}
