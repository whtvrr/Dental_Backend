package repositories

import (
	"context"
	"time"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *entities.Appointment) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Appointment, error)
	Update(ctx context.Context, appointment *entities.Appointment) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, offset, limit int) ([]*entities.Appointment, error)
	GetByDoctorID(ctx context.Context, doctorID primitive.ObjectID, from, to time.Time) ([]*entities.Appointment, error)
	GetByClientID(ctx context.Context, clientID primitive.ObjectID) ([]*entities.Appointment, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]*entities.Appointment, error)
	GetByStatus(ctx context.Context, status entities.AppointmentStatus) ([]*entities.Appointment, error)
}