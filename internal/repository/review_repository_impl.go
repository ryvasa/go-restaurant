package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type ReviewRepositoryImpl struct {
}

func NewReviewRepository() ReviewRepository {
	return &ReviewRepositoryImpl{}
}

func (r *ReviewRepositoryImpl) GetAllByMenuId(tx *sql.Tx, id string) ([]domain.Review, error) {
	reviews := []domain.Review{}
	rows, err := tx.Query("SELECT id,rating,comment,user_id,menu_id,created_at,updated_at FROM review WHERE menu_id = ? AND deleted_at IS NULL AND deleted = false", id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all reviews menu")
		return nil, utils.NewInternalError("Failed to get all reviews menu")
	}
	defer rows.Close()
	for rows.Next() {
		var review domain.Review
		err := rows.Scan(&review.Id, &review.Comment, &review.Rating, &review.UserId, &review.MenuId, &review.UpdatedAt, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *ReviewRepositoryImpl) Create(tx *sql.Tx, review domain.Review) (domain.Review, error) {
	_, err := tx.Exec("INSERT INTO review (rating, comment, user_id, menu_id) VALUES (?, ?, ?, ?)", review.Rating, review.Comment, review.UserId, review.MenuId)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create review")
		return domain.Review{}, utils.NewInternalError("Failed to create review")
	}
	createdReview, err := r.GetOneById(tx, review.Id.String())
	return createdReview, nil
}

func (r *ReviewRepositoryImpl) GetOneById(tx *sql.Tx, id string) (domain.Review, error) {
	var review domain.Review
	err := tx.QueryRow("SELECT * FROM review WHERE id = ? AND deleted_at IS NULL AND deleted = false", id).Scan(&review.Id, &review.Rating, &review.Comment, &review.UserId, &review.MenuId, &review.CreatedAt, &review.UpdatedAt)

	if err != nil {
		logger.Log.WithError(err).Error("Error review not found")
		return domain.Review{}, utils.NewNotFoundError("Review not found")
	}

	return review, nil
}

func (r *ReviewRepositoryImpl) Update(tx *sql.Tx, review domain.Review) (domain.Review, error) {
	_, err := tx.Exec("UPDATE review SET rating = $1, comment = $2, updated_at = $3 WHERE id = $4", review.Rating, review.Comment, review.UpdatedAt, review.Id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update review")
		return domain.Review{}, utils.NewNotFoundError("Failed to update review")
	}
	updatedReview, err := r.GetOneById(tx, review.Id.String())
	return updatedReview, nil
}
