package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentUseCase struct {
	appointmentRepo repositories.AppointmentRepository
	userRepo        repositories.UserRepository
	formulaRepo     repositories.FormulaRepository
}

func NewAppointmentUseCase(
	appointmentRepo repositories.AppointmentRepository,
	userRepo repositories.UserRepository,
	formulaRepo repositories.FormulaRepository,
) *AppointmentUseCase {
	return &AppointmentUseCase{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
		formulaRepo:     formulaRepo,
	}
}

func (uc *AppointmentUseCase) CreateAppointment(ctx context.Context, appointment *entities.Appointment) error {
	// Validate doctor exists
	doctor, err := uc.userRepo.GetByID(ctx, appointment.DoctorID)
	if err != nil {
		return err
	}
	if !doctor.IsDoctor() {
		return errors.New("user is not a doctor")
	}

	// Validate client exists
	client, err := uc.userRepo.GetByID(ctx, appointment.ClientID)
	if err != nil {
		return err
	}
	if !client.IsClient() {
		return errors.New("user is not a client")
	}

	appointment.Status = entities.AppointmentStatusScheduled
	return uc.appointmentRepo.Create(ctx, appointment)
}

func (uc *AppointmentUseCase) GetAppointment(ctx context.Context, id primitive.ObjectID) (*entities.Appointment, error) {
	return uc.appointmentRepo.GetByID(ctx, id)
}

func (uc *AppointmentUseCase) UpdateAppointment(ctx context.Context, appointment *entities.Appointment) error {
	return uc.appointmentRepo.Update(ctx, appointment)
}

func (uc *AppointmentUseCase) CompleteAppointment(ctx context.Context, id primitive.ObjectID, medicalData *AppointmentMedicalData) error {
	appointment, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update appointment with medical data
	appointment.ComplaintID = medicalData.ComplaintID
	appointment.CustomComplaint = medicalData.CustomComplaint
	appointment.Anamnesis = medicalData.Anamnesis
	appointment.DiagnosisID = medicalData.DiagnosisID
	appointment.TreatmentID = medicalData.TreatmentID
	appointment.Comment = medicalData.Comment
	appointment.Status = entities.AppointmentStatusCompleted

	// Handle formula logic based on user's FormulaID field
	if medicalData.Formula != nil {
		client, err := uc.userRepo.GetByID(ctx, medicalData.ClientID)
		if err != nil {
			return err
		}

		if client.FormulaID != nil {
			// Update existing formula
			medicalData.Formula.ID = *client.FormulaID
			medicalData.Formula.UserID = medicalData.ClientID
			medicalData.Formula.UpdatedAt = time.Now()

			err := uc.formulaRepo.Update(ctx, medicalData.Formula)
			if err != nil {
				return err
			}
		} else {
			// Create new formula for the user
			medicalData.Formula.UserID = medicalData.ClientID
			medicalData.Formula.CreatedAt = time.Now()
			medicalData.Formula.UpdatedAt = time.Now()

			err := uc.formulaRepo.Create(ctx, medicalData.Formula)
			if err != nil {
				return err
			}

			// Update user's FormulaID
			client.FormulaID = &medicalData.Formula.ID
			client.UpdatedAt = time.Now()
			err = uc.userRepo.Update(ctx, client)
			if err != nil {
				return err
			}
		}

		// Set the formula in the appointment record
		appointment.Formula = medicalData.Formula
	}

	return uc.appointmentRepo.Update(ctx, appointment)
}

func (uc *AppointmentUseCase) CancelAppointment(ctx context.Context, id primitive.ObjectID) error {
	appointment, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	appointment.Status = entities.AppointmentStatusCanceled
	return uc.appointmentRepo.Update(ctx, appointment)
}

func (uc *AppointmentUseCase) DeleteAppointment(ctx context.Context, id primitive.ObjectID) error {
	return uc.appointmentRepo.Delete(ctx, id)
}

func (uc *AppointmentUseCase) ListAppointments(ctx context.Context, offset, limit int) ([]*entities.Appointment, error) {
	return uc.appointmentRepo.List(ctx, offset, limit)
}

func (uc *AppointmentUseCase) GetDoctorAppointments(ctx context.Context, doctorID primitive.ObjectID, from, to time.Time) ([]*entities.Appointment, error) {
	return uc.appointmentRepo.GetByDoctorID(ctx, doctorID, from, to)
}

func (uc *AppointmentUseCase) GetClientAppointments(ctx context.Context, clientID primitive.ObjectID) ([]*entities.Appointment, error) {
	return uc.appointmentRepo.GetByClientID(ctx, clientID)
}

func (uc *AppointmentUseCase) GetAppointmentsByDateRange(ctx context.Context, from, to time.Time) ([]*entities.Appointment, error) {
	return uc.appointmentRepo.GetByDateRange(ctx, from, to)
}

type AppointmentMedicalData struct {
	ComplaintID     *primitive.ObjectID `json:"complaint_id,omitempty"`
	CustomComplaint *string             `json:"custom_complaint,omitempty"`
	Anamnesis       *string             `json:"anamnesis,omitempty"`
	DiagnosisID     *primitive.ObjectID `json:"diagnosis_id,omitempty"`
	TreatmentID     *primitive.ObjectID `json:"treatment_id,omitempty"`
	Comment         *string             `json:"comment,omitempty"`
	ClientID        primitive.ObjectID  `json:"client_id" binding:"required"`
	Formula         *entities.Formula   `json:"formula,omitempty"`
}
