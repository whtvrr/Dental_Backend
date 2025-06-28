package delivery

import "time"

type CreateAppointmentDTO struct {
	DateTime        time.Time `json:"dateTime"` //TODO:parce time in handler and get as a string
	DoctorId        string    `json:"doctorId"`
	ClientId        string    `json:"clientId"`
	DurationMinutes int       `json:"durationMinutes"`
	Comment         string    `json:"comment"`
}

type UpdateAppointmentDTO struct {
	Id              string    `json:"id"`
	DateTime        time.Time `json:"dateTime"`
	DoctorId        string    `json:"doctorId"`
	ClientId        string    `json:"clientId"`
	DurationMinutes int       `json:"durationMinutes"`
	Comment         string    `json:"comment"`
}

type DeleteAppointmentDTO struct {
	Id string `json:"id"`
}
