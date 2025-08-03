package usecases

import (
	"context"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComplaintUseCase struct {
	complaintRepo repositories.ComplaintRepository
}

func NewComplaintUseCase(complaintRepo repositories.ComplaintRepository) *ComplaintUseCase {
	return &ComplaintUseCase{
		complaintRepo: complaintRepo,
	}
}

func (uc *ComplaintUseCase) CreateComplaint(ctx context.Context, complaint *entities.Complaint) error {
	complaint.IsActive = true
	return uc.complaintRepo.Create(ctx, complaint)
}

func (uc *ComplaintUseCase) GetComplaint(ctx context.Context, id primitive.ObjectID) (*entities.Complaint, error) {
	return uc.complaintRepo.GetByID(ctx, id)
}

func (uc *ComplaintUseCase) UpdateComplaint(ctx context.Context, complaint *entities.Complaint) error {
	return uc.complaintRepo.Update(ctx, complaint)
}

func (uc *ComplaintUseCase) DeleteComplaint(ctx context.Context, id primitive.ObjectID) error {
	return uc.complaintRepo.Delete(ctx, id)
}

func (uc *ComplaintUseCase) ListComplaints(ctx context.Context, offset, limit int) ([]*entities.Complaint, error) {
	return uc.complaintRepo.List(ctx, offset, limit)
}

func (uc *ComplaintUseCase) GetActiveComplaints(ctx context.Context) ([]*entities.Complaint, error) {
	return uc.complaintRepo.GetActive(ctx)
}

func (uc *ComplaintUseCase) GetComplaintsByCategory(ctx context.Context, category string) ([]*entities.Complaint, error) {
	return uc.complaintRepo.GetByCategory(ctx, category)
}