package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Formula struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Teeth     []Tooth            `bson:"teeth" json:"teeth"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Tooth struct {
	Number   int                     `bson:"number" json:"number"`
	Whole    *ToothStatus            `bson:"whole,omitempty" json:"whole,omitempty"`
	Gum      *ToothStatus            `bson:"gum,omitempty" json:"gum,omitempty"`
	Roots    *ToothStatus            `bson:"roots,omitempty" json:"roots,omitempty"`
	Segments map[string]*ToothStatus `bson:"segments,omitempty" json:"segments,omitempty"`
}

type ToothStatus struct {
	StatusID      primitive.ObjectID `bson:"status_id" json:"status_id"`
	AppointmentID primitive.ObjectID `bson:"appointment_id" json:"appointment_id"`
	Timestamp     time.Time          `bson:"timestamp" json:"timestamp"`
	Note          *string            `bson:"note,omitempty" json:"note,omitempty"`
}

// Initialize a new formula with 32 teeth
func NewFormula(userID primitive.ObjectID) *Formula {
	teeth := make([]Tooth, 32)
	for i := 0; i < 32; i++ {
		teeth[i] = Tooth{
			Number:   i + 1,
			Segments: make(map[string]*ToothStatus),
		}
	}

	return &Formula{
		UserID:    userID,
		Teeth:     teeth,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
