package usecases

import (
	"context"
	"errors"

	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserUseCase struct {
	userRepo    repositories.UserRepository
	formulaRepo repositories.FormulaRepository
}

func NewUserUseCase(userRepo repositories.UserRepository, formulaRepo repositories.FormulaRepository) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		formulaRepo: formulaRepo,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user *entities.User) error {
	if user.Role == entities.RoleClient {
		user.Email = nil
		user.PasswordHash = nil
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Create formula for client
	if user.Role == entities.RoleClient {
		formula := entities.NewFormula(user.ID)
		if err := uc.formulaRepo.Create(ctx, formula); err != nil {
			return err
		}
		user.FormulaID = &formula.ID
		return uc.userRepo.Update(ctx, user)
	}

	return nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return uc.userRepo.GetByEmail(ctx, email)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user *entities.User) error {
	existing, err := uc.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}

	if existing.Role == entities.RoleClient {
		user.Email = nil
		user.PasswordHash = nil
	}

	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete formula if client
	if user.Role == entities.RoleClient && user.FormulaID != nil {
		if err := uc.formulaRepo.Delete(ctx, *user.FormulaID); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
	}

	return uc.userRepo.Delete(ctx, id)
}

func (uc *UserUseCase) ListUsers(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	return uc.userRepo.List(ctx, offset, limit)
}

func (uc *UserUseCase) GetDoctors(ctx context.Context) ([]*entities.User, error) {
	return uc.userRepo.GetByRole(ctx, entities.RoleDoctor)
}

func (uc *UserUseCase) GetDoctorsWithPagination(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	return uc.userRepo.GetByRoleWithPagination(ctx, entities.RoleDoctor, offset, limit)
}

func (uc *UserUseCase) GetClients(ctx context.Context) ([]*entities.User, error) {
	return uc.userRepo.GetByRole(ctx, entities.RoleClient)
}

func (uc *UserUseCase) GetStaff(ctx context.Context) ([]*entities.User, error) {
	doctors, err := uc.userRepo.GetByRole(ctx, entities.RoleDoctor)
	if err != nil {
		return nil, err
	}
	
	receptionists, err := uc.userRepo.GetByRole(ctx, entities.RoleReceptionist)
	if err != nil {
		return nil, err
	}
	
	staff := make([]*entities.User, 0, len(doctors)+len(receptionists))
	staff = append(staff, doctors...)
	staff = append(staff, receptionists...)
	
	return staff, nil
}

func (uc *UserUseCase) GetStaffWithPagination(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	doctors, err := uc.userRepo.GetByRoleWithPagination(ctx, entities.RoleDoctor, offset, limit)
	if err != nil {
		return nil, err
	}
	
	// Adjust offset and limit for receptionists based on how many doctors we got
	remainingLimit := limit - len(doctors)
	receptionistOffset := 0
	if offset > len(doctors) {
		receptionistOffset = offset - len(doctors)
	}
	
	if remainingLimit > 0 {
		receptionists, err := uc.userRepo.GetByRoleWithPagination(ctx, entities.RoleReceptionist, receptionistOffset, remainingLimit)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, receptionists...)
	}
	
	return doctors, nil
}

func (uc *UserUseCase) GetClientsWithPagination(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	return uc.userRepo.GetByRoleWithPagination(ctx, entities.RoleClient, offset, limit)
}
