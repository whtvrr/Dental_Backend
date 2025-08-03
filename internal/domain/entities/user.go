package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleDoctor       UserRole = "doctor"
	RoleReceptionist UserRole = "receptionist"
	RoleClient       UserRole = "client"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        *string            `bson:"email,omitempty" json:"email,omitempty"`
	PasswordHash *string            `bson:"password_hash,omitempty" json:"-"`
	Role         UserRole           `bson:"role" json:"role"`
	FullName     string             `bson:"full_name" json:"full_name"`
	PhoneNumber  *string            `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Address      *string            `bson:"address,omitempty" json:"address,omitempty"`
	Gender       *string            `bson:"gender,omitempty" json:"gender,omitempty"`
	BirthDate    *time.Time         `bson:"birth_date,omitempty" json:"birth_date,omitempty"`
	FormulaID    *primitive.ObjectID `bson:"formula_id,omitempty" json:"formula_id,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

func (u *User) IsClient() bool {
	return u.Role == RoleClient
}

func (u *User) IsDoctor() bool {
	return u.Role == RoleDoctor
}

func (u *User) IsReceptionist() bool {
	return u.Role == RoleReceptionist
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) CanAuthenticate() bool {
	return u.Role == RoleAdmin || u.Role == RoleDoctor || u.Role == RoleReceptionist
}