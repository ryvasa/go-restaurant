package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type UserUsecaseImpl struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{userRepo}
}

func (u *UserUsecaseImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all users")
		return nil, utils.NewInternalError("Failed to get all users")
	}
	return users, nil
}

func (u *UserUsecaseImpl) Create(ctx context.Context, req dto.CreateUserRequest) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.User{}, utils.NewValidationError(err)
	}

	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if user.Email == req.Email {
		logger.Log.WithError(err).Error("Error email already exists")
		return domain.User{}, utils.NewConflictError("Email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to hash password")
		return domain.User{}, utils.NewInternalError("Failed to hash password")
	}

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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	user, err := u.userRepo.Get(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.User{}, utils.NewNotFoundError("User not found")
	}
	return user, nil
}

func (u *UserUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateUserRequest) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.User{}, utils.NewValidationError(err)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		return domain.User{}, utils.NewValidationError("Invalid id format")
	}

	_, err = u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		logger.Log.WithError(err).Error("Error not found")
		return domain.User{}, utils.NewNotFoundError("User not found")
	}

	user := domain.User{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
		Phone: req.Phone,
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to hash password")
			return domain.User{}, utils.NewInternalError("Failed to hash password")
		}
		user.Password = hashedPassword
	}

	return u.userRepo.Update(ctx, user)
}
