package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusType string

const (
	StatusTypeDiagnosis StatusType = "diagnosis"
	StatusTypeTreatment StatusType = "treatment"
	StatusTypeTooth     StatusType = "tooth"
)

type Status struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Type        StatusType         `bson:"type" json:"type"`
	Code        string             `bson:"code" json:"code"` // ICD-10 codes like K04, D05
	Description *string            `bson:"description,omitempty" json:"description,omitempty"`
	Color       *string            `bson:"color,omitempty" json:"color,omitempty"` // For tooth formula visualization
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}