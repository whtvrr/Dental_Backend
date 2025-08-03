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

type appointmentRepository struct {
	collection *mongo.Collection
}

func NewAppointmentRepository(db *mongo.Database) repositories.AppointmentRepository {
	return &appointmentRepository{
		collection: db.Collection("appointments"),
	}
}

func (r *appointmentRepository) Create(ctx context.Context, appointment *entities.Appointment) error {
	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, appointment)
	if err != nil {
		return err
	}
	
	appointment.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *appointmentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Appointment, error) {
	var appointment entities.Appointment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&appointment)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *appointmentRepository) Update(ctx context.Context, appointment *entities.Appointment) error {
	appointment.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": appointment.ID},
		bson.M{"$set": appointment},
	)
	return err
}

func (r *appointmentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *appointmentRepository) List(ctx context.Context, offset, limit int) ([]*entities.Appointment, error) {
	opts := options.Find()
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "date_time", Value: 1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var appointments []*entities.Appointment
	for cursor.Next(ctx) {
		var appointment entities.Appointment
		if err := cursor.Decode(&appointment); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	
	return appointments, cursor.Err()
}

func (r *appointmentRepository) GetByDoctorID(ctx context.Context, doctorID primitive.ObjectID, from, to time.Time) ([]*entities.Appointment, error) {
	filter := bson.M{
		"doctor_id": doctorID,
		"date_time": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	
	opts := options.Find().SetSort(bson.D{{Key: "date_time", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var appointments []*entities.Appointment
	for cursor.Next(ctx) {
		var appointment entities.Appointment
		if err := cursor.Decode(&appointment); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	
	return appointments, cursor.Err()
}

func (r *appointmentRepository) GetByClientID(ctx context.Context, clientID primitive.ObjectID) ([]*entities.Appointment, error) {
	opts := options.Find().SetSort(bson.D{{Key: "date_time", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"client_id": clientID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var appointments []*entities.Appointment
	for cursor.Next(ctx) {
		var appointment entities.Appointment
		if err := cursor.Decode(&appointment); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	
	return appointments, cursor.Err()
}

func (r *appointmentRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]*entities.Appointment, error) {
	filter := bson.M{
		"date_time": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	
	opts := options.Find().SetSort(bson.D{{Key: "date_time", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var appointments []*entities.Appointment
	for cursor.Next(ctx) {
		var appointment entities.Appointment
		if err := cursor.Decode(&appointment); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	
	return appointments, cursor.Err()
}

func (r *appointmentRepository) GetByStatus(ctx context.Context, status entities.AppointmentStatus) ([]*entities.Appointment, error) {
	opts := options.Find().SetSort(bson.D{{Key: "date_time", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var appointments []*entities.Appointment
	for cursor.Next(ctx) {
		var appointment entities.Appointment
		if err := cursor.Decode(&appointment); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	
	return appointments, cursor.Err()
}