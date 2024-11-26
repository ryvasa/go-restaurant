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

type ReviewUsecaseImpl struct {
	db         *sql.DB
	reviewRepo repository.ReviewRepository
	userRepo   repository.UserRepository
	menuRepo   repository.MenuRepository
}

func NewReviewUsecase(db *sql.DB, reviewRepo repository.ReviewRepository, userRepo repository.UserRepository, menuRepo repository.MenuRepository) ReviewUsecase {
	return &ReviewUsecaseImpl{db, reviewRepo, userRepo, menuRepo}
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
		logger.Log.WithError(err).Error("Error failed to get all reviews menu")
		return nil, utils.NewInternalError("Failed to get all reviews menu")
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
	user, err := u.userRepo.Get(tx, userId)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.Review{}, utils.NewNotFoundError("User not found")
	}

	menuId := req.MenuId
	menu, err := u.menuRepo.Get(tx, menuId)
	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		return domain.Review{}, utils.NewNotFoundError("Menu not found")
	}

	review := domain.Review{
		Id:      uuid.New(),
		Rating:  req.Rating,
		Comment: req.Comment,
		UserId:  user.Id,
		MenuId:  menu.Id,
	}
	err = u.reviewRepo.Create(tx, review)
	if err != nil {
		return domain.Review{}, err
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

func (u *ReviewUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateReviewRequest) (domain.Review, error) {
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

	_, err = u.reviewRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error review not found")
		return domain.Review{}, utils.NewNotFoundError("Review not found")
	}

	review := domain.Review{
		Id:      reviewId,
		Rating:  req.Rating,
		Comment: req.Comment,
	}
	err = u.reviewRepo.Update(tx, review)
	if err != nil {
		return domain.Review{}, err
	}
	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Review{}, utils.NewInternalError("Failed to commit transaction")
	}
	return review, nil
}
