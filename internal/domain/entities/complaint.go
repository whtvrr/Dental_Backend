package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Complaint struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description *string            `bson:"description,omitempty" json:"description,omitempty"`
	Category    *string            `bson:"category,omitempty" json:"category,omitempty"` // pain, cosmetic, etc.
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}