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
	
	// Formula changes made during this appointment
	FormulaChanges   []FormulaChange    `bson:"formula_changes,omitempty" json:"formula_changes,omitempty"`
	
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type FormulaChange struct {
	ToothNumber int                `bson:"tooth_number" json:"tooth_number"`
	Part        string             `bson:"part" json:"part"` // "whole", "gum", "roots", or segment name
	StatusID    primitive.ObjectID `bson:"status_id" json:"status_id"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
}