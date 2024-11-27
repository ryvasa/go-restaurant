package usecase

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type ReviewUsecaseImpl struct {
	db         *sql.DB
	reviewRepo repository.ReviewRepository
	userRepo   repository.UserRepository
	menuRepo   repository.MenuRepository
	orderRepo  repository.OrderRepository
}

func NewReviewUsecase(db *sql.DB, reviewRepo repository.ReviewRepository, userRepo repository.UserRepository, menuRepo repository.MenuRepository, orderRepo repository.OrderRepository) ReviewUsecase {
	return &ReviewUsecaseImpl{db, reviewRepo, userRepo, menuRepo, orderRepo}
}

func (u *ReviewUsecaseImpl) GetAllByMenuId(ctx context.Context, id string) ([]domain.Review, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return []domain.Review{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	reviews, err := u.reviewRepo.GetAllByMenuId(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get reviews")
		return []domain.Review{}, utils.NewInternalError("Failed to get reviews")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return []domain.Review{}, utils.NewInternalError("Failed to commit transaction")
	}
	return reviews, nil
}

func (u *ReviewUsecaseImpl) Create(ctx context.Context, req dto.CreateReviewRequest, userId string) (domain.Review, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Review{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Review{}, utils.NewValidationError(err)
	}
	if reviewed := u.reviewRepo.CheckReviewedItem(tx, userId, req.MenuId, req.OrderId); reviewed {
		logger.Log.WithField("user_id", userId).WithField("menu_id", req.MenuId).WithField("order_id", req.OrderId).Error("You cannot leave a review twice")
		return domain.Review{}, utils.NewBadRequestError("You cannot leave a review twice for the same order")
	}

	user, err := u.userRepo.Get(tx, userId)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.Review{}, utils.NewInternalError("User not found")
	}
	order, err := u.orderRepo.GetOneById(tx, req.OrderId)
	if err != nil {
		logger.Log.WithError(err).Error("Error order not found")
		return domain.Review{}, utils.NewInternalError("Order not found")
	}
	if order.PaymentStatus != "paid" && order.Status != "success" {
		logger.Log.WithField("paymeny_status", order.PaymentStatus).WithField("status", order.Status).Error("Order not finished yet")
		return domain.Review{}, utils.NewUnauthorizedError("Order not finished yet")
	}
	if order.UserId != user.Id {
		logger.Log.WithField("user_id", user.Id).WithField("order_id.user_id", order.UserId).Error("You cannot leave a review")
		return domain.Review{}, utils.NewUnauthorizedError("You cannot leave a review")
	}

	menu, err := u.menuRepo.Get(tx, req.MenuId)
	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		return domain.Review{}, utils.NewNotFoundError("Menu not found")
	}

	review := domain.Review{
		Id:      uuid.New(),
		Rating:  req.Rating,
		Comment: req.Comment,
		OrderId: order.Id,
		UserId:  user.Id,
		MenuId:  menu.Id,
	}
	review, err = u.reviewRepo.Create(tx, review)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create review")
		return domain.Review{}, utils.NewInternalError("Failed to create review")
	}

	count, rating, err := u.reviewRepo.CountReviewByMenuId(tx, req.MenuId)
	avgRating := rating / float64(count)
	err = u.menuRepo.UpdateRating(tx, req.MenuId, avgRating)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update rating")
		return domain.Review{}, utils.NewInternalError("Failed to update rating")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Review{}, utils.NewInternalError("Failed to commit transaction")
	}
	return review, nil
}

func (u *ReviewUsecaseImpl) GetOneById(ctx context.Context, id string) (domain.Review, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Review{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	review, err := u.reviewRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error review not found")
		return domain.Review{}, utils.NewNotFoundError("Review not found")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Review{}, utils.NewInternalError("Failed to commit transaction")
	}
	return review, err
}

func (u *ReviewUsecaseImpl) Update(ctx context.Context, id, userId string, req dto.UpdateReviewRequest) (domain.Review, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Review{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	reviewId, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		return domain.Review{}, utils.NewValidationError("Invalid id format")
	}

	existingReview, err := u.reviewRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error review not found")
		return domain.Review{}, utils.NewNotFoundError("Review not found")
	}
	authId, err := uuid.Parse(userId)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid user id format")
		return domain.Review{}, utils.NewValidationError("Invalid user id format")
	}

	if existingReview.UserId != authId {
		logger.Log.WithField("user_id", authId).WithField("review.user_id", existingReview.UserId).Error("You cannot update a review")
		return domain.Review{}, utils.NewUnauthorizedError("You cannot update a review")
	}

	review := domain.Review{
		Id:      reviewId,
		Rating:  req.Rating,
		Comment: req.Comment,
	}
	updatedReview, err := u.reviewRepo.Update(tx, review)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update review")
		return domain.Review{}, utils.NewInternalError("Failed to update review")
	}
	menu, err := u.menuRepo.Get(tx, existingReview.MenuId.String())
	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		return domain.Review{}, utils.NewNotFoundError("menu not found")
	}

	count, rating, err := u.reviewRepo.CountReviewByMenuId(tx, menu.Id.String())
	avgRating := rating / float64(count)
	log.Println(count, rating, avgRating)
	err = u.menuRepo.UpdateRating(tx, menu.Id.String(), avgRating)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update rating")
		return domain.Review{}, utils.NewInternalError("Failed to update rating")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Review{}, utils.NewInternalError("Failed to commit transaction")
	}
	return updatedReview, nil
}
