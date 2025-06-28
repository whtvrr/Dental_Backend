package delivery

import (
	"time"
)

type CreateStaffDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

type CreateClientDTO struct {
	Role        string      `json:"role"`
	FullName    string      `json:"fullName"`
	PhoneNumber string      `json:"phoneNumber"`
	Address     string      `json:"address"`
	Gender      bool        `json:"gender"`
	BirthDate   time.Time   `json:"birthDate"`
	Formula     *FormulaDTO `json:"formula,omitempty"`
	Email       string      `json:"email"` //optional
}

type UpdateStaffDTO struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

type UpdateClientDTO struct {
	Id          string      `json:"id"`
	Role        string      `json:"role"`
	FullName    string      `json:"fullName"`
	PhoneNumber string      `json:"phoneNumber"`
	Address     string      `json:"address"`
	Gender      bool        `json:"gender"`
	BirthDate   time.Time   `json:"birthDate"`
	Formula     *FormulaDTO `json:"formula,omitempty"`
	Email       string      `json:"email,omitempty"` //optional
}

type ToothStatusDTO struct {
	StatusID string `json:"statusId"`
	VisitID  string `json:"visitId"`
}

type ToothPartsDTO struct {
	Number   int                        `json:"number" `
	Whole    *ToothStatusDTO            `json:"whole,omitempty"`
	Gum      *ToothStatusDTO            `json:"gum,omitempty"`
	Roots    *ToothStatusDTO            `json:"roots,omitempty"`
	Segments map[string]*ToothStatusDTO `json:"segments,omitempty"`
}

type FormulaDTO struct {
	Teeth []ToothPartsDTO `json:"teeth"`
}
