package usecases

import (
	"context"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnamnesisUseCase struct {
	anamnesisRepo repositories.AnamnesisRepository
}

func NewAnamnesisUseCase(anamnesisRepo repositories.AnamnesisRepository) *AnamnesisUseCase {
	return &AnamnesisUseCase{
		anamnesisRepo: anamnesisRepo,
	}
}

func (uc *AnamnesisUseCase) CreateAnamnesis(ctx context.Context, anamnesis *entities.Anamnesis) error {
	return uc.anamnesisRepo.Create(ctx, anamnesis)
}

func (uc *AnamnesisUseCase) GetAnamnesis(ctx context.Context, id primitive.ObjectID) (*entities.Anamnesis, error) {
	return uc.anamnesisRepo.GetByID(ctx, id)
}

func (uc *AnamnesisUseCase) UpdateAnamnesis(ctx context.Context, anamnesis *entities.Anamnesis) error {
	return uc.anamnesisRepo.Update(ctx, anamnesis)
}

func (uc *AnamnesisUseCase) DeleteAnamnesis(ctx context.Context, id primitive.ObjectID) error {
	return uc.anamnesisRepo.Delete(ctx, id)
}

func (uc *AnamnesisUseCase) ListAnamnesis(ctx context.Context, offset, limit int) ([]*entities.Anamnesis, error) {
	return uc.anamnesisRepo.List(ctx, offset, limit)
}
