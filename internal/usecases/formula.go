package usecases

import (
	"context"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FormulaUseCase struct {
	formulaRepo repositories.FormulaRepository
}

func NewFormulaUseCase(formulaRepo repositories.FormulaRepository) *FormulaUseCase {
	return &FormulaUseCase{
		formulaRepo: formulaRepo,
	}
}

func (uc *FormulaUseCase) GetFormula(ctx context.Context, id primitive.ObjectID) (*entities.Formula, error) {
	return uc.formulaRepo.GetByID(ctx, id)
}

func (uc *FormulaUseCase) GetFormulaByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.Formula, error) {
	return uc.formulaRepo.GetByUserID(ctx, userID)
}

func (uc *FormulaUseCase) UpdateFormula(ctx context.Context, formula *entities.Formula) error {
	return uc.formulaRepo.Update(ctx, formula)
}

func (uc *FormulaUseCase) UpdateToothStatus(ctx context.Context, formulaID primitive.ObjectID, toothNumber int, part string, status *entities.ToothStatus) error {
	return uc.formulaRepo.UpdateToothStatus(ctx, formulaID, toothNumber, part, status)
}