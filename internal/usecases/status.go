package usecases

import (
	"context"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusUseCase struct {
	statusRepo repositories.StatusRepository
}

func NewStatusUseCase(statusRepo repositories.StatusRepository) *StatusUseCase {
	return &StatusUseCase{
		statusRepo: statusRepo,
	}
}

func (uc *StatusUseCase) CreateStatus(ctx context.Context, status *entities.Status) error {
	status.IsActive = true
	return uc.statusRepo.Create(ctx, status)
}

func (uc *StatusUseCase) GetStatus(ctx context.Context, id primitive.ObjectID) (*entities.Status, error) {
	return uc.statusRepo.GetByID(ctx, id)
}

func (uc *StatusUseCase) UpdateStatus(ctx context.Context, status *entities.Status) error {
	return uc.statusRepo.Update(ctx, status)
}

func (uc *StatusUseCase) DeleteStatus(ctx context.Context, id primitive.ObjectID) error {
	return uc.statusRepo.Delete(ctx, id)
}

func (uc *StatusUseCase) ListStatuses(ctx context.Context, offset, limit int) ([]*entities.Status, error) {
	return uc.statusRepo.List(ctx, offset, limit)
}

func (uc *StatusUseCase) GetStatusesByType(ctx context.Context, statusType entities.StatusType) ([]*entities.Status, error) {
	return uc.statusRepo.GetByType(ctx, statusType)
}

func (uc *StatusUseCase) GetActiveStatusesByType(ctx context.Context, statusType entities.StatusType) ([]*entities.Status, error) {
	return uc.statusRepo.GetActiveByType(ctx, statusType)
}

func (uc *StatusUseCase) GetDiagnosisStatuses(ctx context.Context) ([]*entities.Status, error) {
	return uc.statusRepo.GetActiveByType(ctx, entities.StatusTypeDiagnosis)
}

func (uc *StatusUseCase) GetTreatmentStatuses(ctx context.Context) ([]*entities.Status, error) {
	return uc.statusRepo.GetActiveByType(ctx, entities.StatusTypeTreatment)
}

func (uc *StatusUseCase) GetToothStatuses(ctx context.Context) ([]*entities.Status, error) {
	return uc.statusRepo.GetActiveByType(ctx, entities.StatusTypeTooth)
}