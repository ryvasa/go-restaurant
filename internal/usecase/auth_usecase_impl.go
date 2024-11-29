package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type AuthUsecaseImpl struct {
	userRepo  repository.UserRepository
	tokenUtil *utils.TokenUtil
	txRepo    repository.TransactionRepository
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	tokenUtil *utils.TokenUtil,
	txRepo repository.TransactionRepository,
) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:  userRepo,
		tokenUtil: tokenUtil,
		txRepo:    txRepo,
	}
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, req dto.LoginDto) (domain.Auth, error) {
	result := domain.Auth{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		user, err := adapters.UserRepository.GetByEmail(ctx, req.Email)
		if err != nil {
			logger.Log.WithError(err).Error("Error user not found")
			return utils.NewNotFoundError("Invalid email or password")
		}

		if !utils.CheckPasswordHash(req.Password, user.Password) {
			logger.Log.Error("Error invalid password")
			return utils.NewUnauthorizedError("Invalid email or password")
		}

		token, err := u.tokenUtil.GenerateToken(user.Id.String(), user.Role)
		if err != nil {
			logger.Log.WithError(err).Error("Error generating token")
			return utils.NewInternalError("Failed to generate token")
		}

		result = domain.Auth{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Token:     token,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
