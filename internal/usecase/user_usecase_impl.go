package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/dto"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/utils"
	"github.com/ryvasa/go-restaurant/internal/domain"
	"github.com/ryvasa/go-restaurant/internal/repository"
)

// TODO: update move bisnis logic to user usecase
type UserUsecaseImpl struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{userRepo}
}

func (u *UserUsecaseImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *UserUsecaseImpl) Create(ctx context.Context, req dto.CreateUserRequest) (domain.User, error) {
	// Validasi input
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return domain.User{}, utils.NewValidationError(err)
	}

	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if user.Email == req.Email {
		return domain.User{}, utils.NewConflictError("Email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return domain.User{}, utils.NewInternalError("Failed to hash password")
	}

	// Buat user baru
	user = domain.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Password: hashedPassword,
		Email:    req.Email,
		Role:     "customer",
	}

	return u.userRepo.Create(ctx, user)
}

func (u *UserUsecaseImpl) Get(ctx context.Context, id string) (domain.User, error) {
	return u.userRepo.Get(ctx, id)
}

func (u *UserUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateUserRequest) (domain.User, error) {
	// Validasi input
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return domain.User{}, fmt.Errorf("validation error: %v", err)
	}

	// Parse UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		return domain.User{}, err
	}

	// Siapkan data update
	user := domain.User{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
		Phone: req.Phone,
	}

	// Hash password jika ada
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return domain.User{}, err
		}
		user.Password = hashedPassword
	}

	return u.userRepo.Update(ctx, user)
}
