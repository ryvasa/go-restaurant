package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type UserUsecaseImpl struct {
	userRepo repository.UserRepository
	txRepo   repository.TransactionRepository
}

func NewUserUsecase(userRepo repository.UserRepository, txRepo repository.TransactionRepository) UserUsecase {
	return &UserUsecaseImpl{userRepo, txRepo}
}

func (u *UserUsecaseImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all users")
		return nil, utils.NewInternalError("Failed to get all users")
	}
	return users, nil
}

func (u *UserUsecaseImpl) Create(ctx context.Context, req dto.CreateUserRequest) (domain.User, error) {
	result := domain.User{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		user, err := adapters.UserRepository.GetByEmail(ctx, req.Email)
		if user.Email == req.Email {
			logger.Log.WithError(err).Error("Error email already exists")
			return utils.NewConflictError("Email already exists")
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to hash password")
			return utils.NewInternalError("Failed to hash password")
		}

		user = domain.User{
			Id:       uuid.New(),
			Name:     req.Name,
			Password: hashedPassword,
			Email:    req.Email,
			Role:     "customer",
		}

		err = adapters.UserRepository.Create(ctx, user)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to create user")
			return utils.NewInternalError("Failed to create user")
		}
		createdUser, err := adapters.UserRepository.Get(ctx, user.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get user")
			return utils.NewInternalError("Failed to get user")
		}
		result = createdUser
		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *UserUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := u.userRepo.Get(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.User{}, utils.NewNotFoundError("User not found")
	}
	return user, nil
}

func (u *UserUsecaseImpl) Update(ctx context.Context, id, authId uuid.UUID, req dto.UpdateUserRequest) (domain.User, error) {
	result := domain.User{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		if authId != id {
			logger.Log.WithField("user_id", authId).WithField("update_user_id", id).Error("You cannot update a user")
			return utils.NewUnauthorizedError("You cannot update a user")
		}

		existingUser, err := adapters.UserRepository.Get(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error user not found")
			return utils.NewNotFoundError("User not found")
		}

		existingUser, err = adapters.UserRepository.GetByEmail(ctx, existingUser.Email)
		if err != nil {
			logger.Log.WithError(err).Error("Error user not found")
			return utils.NewNotFoundError("User not found")
		}

		user := domain.User{
			Name:  req.Name,
			Email: req.Email,
			Role:  req.Role,
			Phone: &req.Phone,
		}

		if req.Name == "" {
			user.Name = existingUser.Name
		}
		if req.Phone == "" {
			user.Phone = existingUser.Phone
		}
		if req.Role == "" {
			user.Role = existingUser.Role
		}
		if req.Email == "" {
			user.Email = existingUser.Email
		} else {
			existingUserWithEmail, err := adapters.UserRepository.GetByEmail(ctx, req.Email)
			if err != nil {
				logger.Log.WithError(err).Error("Error user not found")
				return utils.NewNotFoundError("User not found")
			}
			if existingUserWithEmail.Email == req.Email {
				logger.Log.WithError(err).Error("Error email already exists")
				return utils.NewConflictError("Email already exists")
			}
		}

		if req.Password != "" {
			hashedPassword, err := utils.HashPassword(req.Password)
			if err != nil {
				logger.Log.WithError(err).Error("Error failed to hash password")
				return utils.NewInternalError("Failed to hash password")
			}
			user.Password = hashedPassword
		} else {
			user.Password = existingUser.Password
		}

		err = adapters.UserRepository.Update(ctx, id, user)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update user")
			return utils.NewInternalError("Failed to update user")
		}

		updatedUser, err := adapters.UserRepository.Get(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get user")
			return utils.NewInternalError("Failed to get user")
		}
		result = updatedUser
		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *UserUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if _, err := adapters.UserRepository.Get(ctx, id); err != nil {
			logger.Log.WithError(err).Error("Error user not found")
			return utils.NewNotFoundError("User not found")
		}

		err := adapters.UserRepository.Delete(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to delete user")
			return utils.NewInternalError("Failed to delete user")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUsecaseImpl) Restore(ctx context.Context, id uuid.UUID) (domain.User, error) {
	result := domain.User{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		_, err := adapters.UserRepository.GetDeletedUserById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error user not found to restore")
			return utils.NewNotFoundError("User not found to restore")
		}

		err = adapters.UserRepository.Restore(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to restore user")
			return utils.NewInternalError("Failed to restore user")
		}
		restoredUser, err := adapters.UserRepository.Get(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get user")
			return utils.NewInternalError("Failed to get user")
		}
		result = restoredUser

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
