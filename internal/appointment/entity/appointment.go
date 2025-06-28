package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	DateTime        time.Time          `bson:"dateTime"`
	DoctorId        primitive.ObjectID `bson:"doctorId"`
	ClientId        primitive.ObjectID `bson:"clientId"`
	DurationMinutes int                `bson:"durationMinutes"`
	Comment         string             `bson:"comment"`
}
