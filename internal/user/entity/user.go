package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email"`        //optional for clients
	PasswordHash string             `bson:"passwordHash"` //not for clients
	Role         string             `bson:"role"`
	FullName     string             `bson:"fullName"` //this and following for clinets only
	PhoneNumber  string             `bson:"phoneNumber"`
	Address      string             `bson:"address"`
	Gender       bool               `bson:"gender"`
	BirthDate    time.Time          `bson:"birthDate"`
	Formula      *Formula           `bson:"formula"`
}

type ToothStatus struct {
	StatusID primitive.ObjectID `bson:"statusId"`
	VisitID  primitive.ObjectID `bson:"visitId"`
}

type ToothParts struct {
	Number   int                     `bson:"number"`
	Whole    *ToothStatus            `bson:"whole,omitempty"`
	Gum      *ToothStatus            `bson:"gum,omitempty"`
	Roots    *ToothStatus            `bson:"roots,omitempty"`
	Segments map[string]*ToothStatus `bson:"segments,omitempty"`
}

type Formula struct {
	Teeth []ToothParts `bson:"teeth"`
}
