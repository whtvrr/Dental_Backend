package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusInProgress AppointmentStatus = "in_progress"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCanceled  AppointmentStatus = "canceled"
)

type Appointment struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DateTime         time.Time          `bson:"date_time" json:"date_time"`
	DoctorID         primitive.ObjectID `bson:"doctor_id" json:"doctor_id"`
	ClientID         primitive.ObjectID `bson:"client_id" json:"client_id"`
	DurationMinutes  int                `bson:"duration_minutes" json:"duration_minutes"`
	Status           AppointmentStatus  `bson:"status" json:"status"`
	
	// Medical fields - filled during appointment
	ComplaintID      *primitive.ObjectID `bson:"complaint_id,omitempty" json:"complaint_id,omitempty"`
	CustomComplaint  *string            `bson:"custom_complaint,omitempty" json:"custom_complaint,omitempty"`
	Anamnesis        *string            `bson:"anamnesis,omitempty" json:"anamnesis,omitempty"`
	DiagnosisID      *primitive.ObjectID `bson:"diagnosis_id,omitempty" json:"diagnosis_id,omitempty"`
	TreatmentID      *primitive.ObjectID `bson:"treatment_id,omitempty" json:"treatment_id,omitempty"`
	Comment          *string            `bson:"comment,omitempty" json:"comment,omitempty"`
	
	// Formula state after this appointment
	Formula          *Formula           `bson:"formula,omitempty" json:"formula,omitempty"`

	// Teeth numbers treated in this appointment (1-32)
	TeethNumbers     []int              `bson:"teeth_numbers,omitempty" json:"teeth_numbers,omitempty"`

	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

