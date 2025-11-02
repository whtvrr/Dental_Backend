package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Anamnesis struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Text      string             `bson:"text" json:"text"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
