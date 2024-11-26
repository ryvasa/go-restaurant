package usecase

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type UserUsecaseImpl struct {
	db       *sql.DB
	userRepo repository.UserRepository
}

func NewUserUsecase(db *sql.DB, userRepo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{db, userRepo}
}

func (u *UserUsecaseImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return []domain.User{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	users, err := u.userRepo.GetAll(tx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all users")
		return nil, utils.NewInternalError("Failed to get all users")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return []domain.User{}, utils.NewInternalError("Failed to commit transaction")
	}
	return users, nil
}

func (u *UserUsecaseImpl) Create(ctx context.Context, req dto.CreateUserRequest) (domain.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.User{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.User{}, utils.NewValidationError(err)
	}

	user, err := u.userRepo.GetByEmail(tx, req.Email)
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
		Id:       uuid.New(),
		Name:     req.Name,
		Password: hashedPassword,
		Email:    req.Email,
		Role:     "customer",
	}

	createdUser, err := u.userRepo.Create(tx, user)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create user")
		return domain.User{}, utils.NewInternalError("Failed to create user")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.User{}, utils.NewInternalError("Failed to commit transaction")
	}
	return createdUser, nil
}

func (u *UserUsecaseImpl) Get(ctx context.Context, id string) (domain.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.User{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	user, err := u.userRepo.Get(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.User{}, utils.NewNotFoundError("User not found")
	}
	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.User{}, utils.NewInternalError("Failed to commit transaction")
	}
	return user, nil
}

func (u *UserUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateUserRequest) (domain.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.User{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.User{}, utils.NewValidationError(err)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		return domain.User{}, utils.NewValidationError("Invalid id format")
	}

	_, err = u.userRepo.Get(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error not found")
		return domain.User{}, utils.NewNotFoundError("User not found")
	}

	user, err := u.userRepo.GetByEmail(tx, req.Email)
	if user.Email == req.Email {
		logger.Log.WithError(err).Error("Error email already exists")
		return domain.User{}, utils.NewConflictError("Email already exists")
	}

	user = domain.User{
		Id:    userID,
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
		Phone: &req.Phone,
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to hash password")
			return domain.User{}, utils.NewInternalError("Failed to hash password")
		}
		user.Password = hashedPassword
	}

	updatedUser, err := u.userRepo.Update(tx, user)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update user")
		return domain.User{}, utils.NewInternalError("Failed to update user")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.User{}, utils.NewInternalError("Failed to commit transaction")
	}

	return updatedUser, nil
}

func (u *UserUsecaseImpl) Delete(ctx context.Context, id string) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return utils.NewValidationError("Invalid ID format")
	}

	if _, err := u.userRepo.Get(tx, id); err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return utils.NewNotFoundError("User not found")
	}

	err = u.userRepo.Delete(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete user")
		return utils.NewInternalError("Failed to delete user")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return utils.NewInternalError("Failed to commit transaction")
	}

	return nil
}

func (u *UserUsecaseImpl) Restore(ctx context.Context, id string) (domain.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.User{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return domain.User{}, utils.NewValidationError("Invalid ID format")
	}

	if _, err := u.userRepo.GetDeletedUserById(tx, id); err != nil {
		logger.Log.WithError(err).Error("Error user not found to restore")
		return domain.User{}, utils.NewNotFoundError("User not found to restore")
	}

	restoredUser, err := u.userRepo.Restore(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore user")
		return domain.User{}, utils.NewInternalError("Failed to restore user")
	}
	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.User{}, utils.NewInternalError("Failed to commit transaction")
	}
	return restoredUser, nil
}
