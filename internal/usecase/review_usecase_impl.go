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

type ReviewUsecaseImpl struct {
	reviewRepo repository.ReviewRepository
	userRepo   repository.UserRepository
	menuRepo   repository.MenuRepository
	orderRepo  repository.OrderRepository
	txRepo     repository.TransactionRepository
}

func NewReviewUsecase(reviewRepo repository.ReviewRepository, userRepo repository.UserRepository, menuRepo repository.MenuRepository, orderRepo repository.OrderRepository, txRepo repository.TransactionRepository) ReviewUsecase {
	return &ReviewUsecaseImpl{reviewRepo, userRepo, menuRepo, orderRepo, txRepo}
}

func (u *ReviewUsecaseImpl) GetAllByMenuId(ctx context.Context, id uuid.UUID) ([]domain.Review, error) {
	reviews, err := u.reviewRepo.GetAllByMenuId(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get reviews")
		return []domain.Review{}, utils.NewInternalError("Failed to get reviews")
	}
	return reviews, nil
}

func (u *ReviewUsecaseImpl) Create(ctx context.Context, req dto.CreateReviewRequest, userId uuid.UUID) (domain.Review, error) {
	result := domain.Review{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}
		if reviewed := adapters.ReviewRepository.CheckReviewedItem(ctx, userId, req.MenuId, req.OrderId); reviewed {
			logger.Log.WithField("user_id", userId).WithField("menu_id", req.MenuId).WithField("order_id", req.OrderId).Error("You cannot leave a review twice")
			return utils.NewBadRequestError("You cannot leave a review twice for the same order")
		}

		user, err := adapters.UserRepository.Get(ctx, userId)
		if err != nil {
			logger.Log.WithError(err).Error("Error user not found")
			return utils.NewInternalError("User not found")
		}
		order, err := adapters.OrderRepository.GetOneById(ctx, req.OrderId)
		if err != nil {
			logger.Log.WithError(err).Error("Error order not found")
			return utils.NewInternalError("Order not found")
		}
		if order.PaymentStatus != "paid" && order.Status != "success" {
			logger.Log.WithField("paymeny_status", order.PaymentStatus).WithField("status", order.Status).Error("Order not finished yet")
			return utils.NewUnauthorizedError("Order not finished yet")
		}
		if order.UserId != user.Id {
			logger.Log.WithField("user_id", user.Id).WithField("order_id.user_id", order.UserId).Error("You cannot leave a review")
			return utils.NewUnauthorizedError("You cannot leave a review")
		}
		menu, err := adapters.MenuRepository.Get(ctx, req.MenuId)
		if err != nil {
			logger.Log.WithError(err).Error("Error menu not found")
			return utils.NewNotFoundError("Menu not found")
		}
		review := domain.Review{
			Id:      uuid.New(),
			Rating:  req.Rating,
			Comment: req.Comment,
			OrderId: order.Id,
			UserId:  user.Id,
			MenuId:  menu.Id,
		}
		err = adapters.ReviewRepository.Create(ctx, review)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to create review")
			return utils.NewInternalError("Failed to create review")
		}

		count, rating, err := adapters.ReviewRepository.CountReviewByMenuId(ctx, req.MenuId)
		avgRating := rating / float64(count)

		err = adapters.MenuRepository.UpdateRating(ctx, req.MenuId, avgRating)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update rating")
			return utils.NewInternalError("Failed to update rating")
		}

		createdReview, err := adapters.ReviewRepository.GetOneById(ctx, review.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get created review")
			return utils.NewInternalError("Failed to get created review")
		}
		result = createdReview

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *ReviewUsecaseImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Review, error) {
	review, err := u.reviewRepo.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error review not found")
		return domain.Review{}, utils.NewNotFoundError("Review not found")
	}
	return review, err
}

func (u *ReviewUsecaseImpl) Update(ctx context.Context, id, userId uuid.UUID, req dto.UpdateReviewRequest) (domain.Review, error) {
	result := domain.Review{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {

		existingReview, err := adapters.ReviewRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error review not found")
			return utils.NewNotFoundError("Review not found")
		}

		if existingReview.UserId != userId {
			logger.Log.WithField("user_id", userId).WithField("review.user_id", existingReview.UserId).Error("You cannot update a review")
			return utils.NewUnauthorizedError("You cannot update a review")
		}

		if req.Rating != 0 {
			existingReview.Rating = req.Rating
		}
		if req.Comment != "" {
			existingReview.Comment = req.Comment
		}

		err = adapters.ReviewRepository.Update(ctx, id, existingReview)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update review")
			return utils.NewInternalError("Failed to update review")
		}
		menu, err := adapters.MenuRepository.Get(ctx, existingReview.MenuId)
		if err != nil {
			logger.Log.WithError(err).Error("Error menu not found")
			return utils.NewNotFoundError("menu not found")
		}

		count, rating, err := adapters.ReviewRepository.CountReviewByMenuId(ctx, menu.Id)
		avgRating := rating / float64(count)
		err = adapters.MenuRepository.UpdateRating(ctx, menu.Id, avgRating)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update rating")
			return utils.NewInternalError("Failed to update rating")
		}
		updatedReview, err := adapters.ReviewRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get updated review")
			return utils.NewInternalError("Failed to get updated review")
		}
		result = updatedReview

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
