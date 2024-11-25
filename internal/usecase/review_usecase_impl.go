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

type ReviewUsecaseImpl struct {
	reviewRepo repository.ReviewRepository
	userRepo   repository.UserRepository
	menuRepo   repository.MenuRepository
}

func NewReviewUsecase(reviewRepo repository.ReviewRepository, userRepo repository.UserRepository, menuRepo repository.MenuRepository) ReviewUsecase {
	return &ReviewUsecaseImpl{reviewRepo, userRepo, menuRepo}
}

func (u *ReviewUsecaseImpl) GetAllByMenuId(ctx context.Context, id string) ([]domain.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	reviews, err := u.reviewRepo.GetAllByMenuId(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all reviews menu")
		return nil, utils.NewInternalError("Failed to get all reviews menu")
	}
	return reviews, nil
}

func (u *ReviewUsecaseImpl) Create(ctx context.Context, req dto.CreateReviewRequest, userId string) (domain.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Review{}, utils.NewValidationError(err)
	}
	user, err := u.userRepo.Get(ctx, userId)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.Review{}, utils.NewNotFoundError("User not found")
	}

	menuId := req.MenuId
	menu, err := u.menuRepo.Get(ctx, menuId)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.Review{}, utils.NewNotFoundError("Menu not found")
	}

	review := domain.Review{
		Id:      uuid.New(),
		Rating:  req.Rating,
		Comment: req.Comment,
		UserId:  user.Id,
		MenuId:  menu.Id,
	}
	err = u.reviewRepo.Create(ctx, review)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}

func (u *ReviewUsecaseImpl) GetOneById(ctx context.Context, id string) (domain.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	review, err := u.reviewRepo.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error review not found")
		return domain.Review{}, utils.NewNotFoundError("Review not found")
	}
	return review, err
}

func (u *ReviewUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateReviewRequest) (domain.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	reviewId, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		return domain.Review{}, utils.NewValidationError("Invalid id format")
	}

	_, err = u.reviewRepo.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error review not found")
		return domain.Review{}, utils.NewNotFoundError("Review not found")
	}

	review := domain.Review{
		Id:      reviewId,
		Rating:  req.Rating,
		Comment: req.Comment,
	}
	err = u.reviewRepo.Update(ctx, review)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}
