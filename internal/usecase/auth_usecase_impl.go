package usecase

import (
	"context"
	"time"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type AuthUsecaseImpl struct {
	userRepo  repository.UserRepository
	tokenUtil *utils.TokenUtil
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	tokenUtil *utils.TokenUtil,
) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:  userRepo,
		tokenUtil: tokenUtil,
	}
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, req dto.LoginDto) (domain.Auth, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Auth{}, utils.NewValidationError(err)
	}

	// Cari user berdasarkan email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.Auth{}, utils.NewNotFoundError("Invalid email or password")
	}

	// Verifikasi password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		logger.Log.Error("Error invalid password")
		return domain.Auth{}, utils.NewUnauthorizedError("Invalid email or password")
	}

	// Generate JWT token
	token, err := u.tokenUtil.GenerateToken(user.Id.String(), user.Role)
	if err != nil {
		logger.Log.WithError(err).Error("Error generating token")
		return domain.Auth{}, utils.NewInternalError("Failed to generate token")
	}

	// Buat response
	auth := domain.Auth{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Token:     token,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return auth, nil
}
