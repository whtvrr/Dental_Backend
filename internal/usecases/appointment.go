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

	if medicalData.TeethNumbers != nil {
		if err := validateTeethNumbers(medicalData.TeethNumbers); err != nil {
			return err
		}
	}

	// Update appointment with medical data
	appointment.ComplaintID = medicalData.ComplaintID
	appointment.CustomComplaint = medicalData.CustomComplaint
	appointment.Anamnesis = medicalData.Anamnesis
	appointment.DiagnosisID = medicalData.DiagnosisID
	appointment.TreatmentID = medicalData.TreatmentID
	appointment.Comment = medicalData.Comment
	appointment.TeethNumbers = medicalData.TeethNumbers
	appointment.Status = entities.AppointmentStatusCompleted

	// Completion timestamp for all tooth statuses
	completionTime := time.Now()

	// Handle formula logic based on user's FormulaID field
	if medicalData.Formula != nil {
		// Add appointment_id and timestamp to all tooth statuses
		for i := range medicalData.Formula.Teeth {
			tooth := &medicalData.Formula.Teeth[i]

			// Set appointment_id and timestamp for Whole (crown)
			if tooth.Whole != nil {
				tooth.Whole.AppointmentID = id
				tooth.Whole.Timestamp = completionTime
			}

			// Set appointment_id and timestamp for Gum (jaw)
			if tooth.Gum != nil {
				tooth.Gum.AppointmentID = id
				tooth.Gum.Timestamp = completionTime
			}

			// Set appointment_id and timestamp for each Root
			if tooth.Roots != nil {
				for j := range tooth.Roots {
					tooth.Roots[j].AppointmentID = id
					tooth.Roots[j].Timestamp = completionTime
				}
			}

			// Set appointment_id and timestamp for each Segment
			if tooth.Segments != nil {
				for key := range tooth.Segments {
					tooth.Segments[key].AppointmentID = id
					tooth.Segments[key].Timestamp = completionTime
				}
			}
		}

		// Store appointment's own formula
		appointmentFormula := medicalData.Formula

		client, err := uc.userRepo.GetByID(ctx, medicalData.ClientID)
		if err != nil {
			return err
		}

		if client.FormulaID != nil {
			// Get existing formula
			existingFormula, err := uc.formulaRepo.GetByID(ctx, *client.FormulaID)
			if err != nil {
				return err
			}

			// Merge appointment formula with existing formula
			mergedFormula := uc.mergeFormulas(existingFormula, medicalData.Formula)
			mergedFormula.ID = *client.FormulaID
			mergedFormula.UserID = medicalData.ClientID
			mergedFormula.UpdatedAt = completionTime

			err = uc.formulaRepo.Update(ctx, mergedFormula)
			if err != nil {
				return err
			}
		} else {
			medicalData.Formula.UserID = medicalData.ClientID
			medicalData.Formula.CreatedAt = completionTime
			medicalData.Formula.UpdatedAt = completionTime

			err := uc.formulaRepo.Create(ctx, medicalData.Formula)
			if err != nil {
				return err
			}

			// Update user's FormulaID
			client.FormulaID = &medicalData.Formula.ID
			client.UpdatedAt = completionTime
			err = uc.userRepo.Update(ctx, client)
			if err != nil {
				return err
			}
		}

		// Set only the appointment's own formula in the appointment record
		appointment.Formula = appointmentFormula
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
	TeethNumbers    []int               `json:"teeth_numbers,omitempty"`
}

// mergeFormulas merges appointment formula into existing user formula
// Teeth numbers are the main identifiers. If the same tooth exists in both:
// - For Whole, Gum: overwrite with appointment data if present
// - For Roots: always overwrite with appointment data
// - For Segments: merge by segment key (mid, rt, lt, rb, lb), overwriting conflicts
func (uc *AppointmentUseCase) mergeFormulas(existing, appointment *entities.Formula) *entities.Formula {
	// Create a map of existing teeth for quick lookup
	existingTeethMap := make(map[int]*entities.Tooth)
	for i := range existing.Teeth {
		existingTeethMap[existing.Teeth[i].Number] = &existing.Teeth[i]
	}

	// Process appointment teeth
	for i := range appointment.Teeth {
		appointmentTooth := &appointment.Teeth[i]

		if existingTooth, exists := existingTeethMap[appointmentTooth.Number]; exists {
			// Merge tooth parts
			uc.mergeToothParts(existingTooth, appointmentTooth)
		} else {
			// Add new tooth from appointment
			existing.Teeth = append(existing.Teeth, *appointmentTooth)
			existingTeethMap[appointmentTooth.Number] = appointmentTooth
		}
	}

	return existing
}

// mergeToothParts merges individual tooth parts according to the rules:
// - Whole/Gum: overwrite if appointment has data
// - Roots: always overwrite with appointment data
// - Segments: merge by key, overwriting conflicts
func (uc *AppointmentUseCase) mergeToothParts(existing, appointment *entities.Tooth) {
	// Merge Whole (crown) - overwrite if appointment has data
	if appointment.Whole != nil {
		existing.Whole = appointment.Whole
	}

	// Merge Gum (jaw) - overwrite if appointment has data
	if appointment.Gum != nil {
		existing.Gum = appointment.Gum
	}

	// Merge Roots - always overwrite with appointment data
	if appointment.Roots != nil {
		existing.Roots = appointment.Roots
	}

	// Merge Segments - merge by key, overwriting conflicts
	if appointment.Segments != nil {
		if existing.Segments == nil {
			existing.Segments = make(map[string]*entities.ToothStatus)
		}

		// Valid segment keys: mid, rt, lt, rb, lb
		for key, status := range appointment.Segments {
			existing.Segments[key] = status
		}
	}
}

// validateTeethNumbers validates that all teeth numbers are between 1-32 and unique
func validateTeethNumbers(teethNumbers []int) error {
	if len(teethNumbers) == 0 {
		return nil
	}

	seen := make(map[int]bool)
	for _, num := range teethNumbers {
		if num < 1 || num > 32 {
			return errors.New("teeth numbers must be between 1 and 32")
		}
		if seen[num] {
			return errors.New("duplicate teeth numbers are not allowed")
		}
		seen[num] = true
	}
	return nil
}
