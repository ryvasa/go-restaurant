package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type ReviewRepositoryImpl struct {
}

func NewReviewRepository() ReviewRepository {
	return &ReviewRepositoryImpl{}
}

func (r *ReviewRepositoryImpl) GetAllByMenuId(tx *sql.Tx, menuId string) ([]domain.Review, error) {
	reviews := []domain.Review{}
	rows, err := tx.Query("SELECT id,rating,comment,user_id,menu_id,order_id,created_at,updated_at FROM review WHERE menu_id = ? AND deleted_at IS NULL AND deleted = false", menuId)
	if err != nil {
		return []domain.Review{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var review domain.Review
		err := rows.Scan(&review.Id, &review.Comment, &review.Rating, &review.UserId, &review.MenuId, &review.OrderId, &review.UpdatedAt, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *ReviewRepositoryImpl) Create(tx *sql.Tx, review domain.Review) (domain.Review, error) {
	_, err := tx.Exec("INSERT INTO review (id,rating, comment, user_id, menu_id, order_id) VALUES (?, ?, ?, ?, ?, ?)", review.Id, review.Rating, review.Comment, review.UserId, review.MenuId, review.OrderId)
	if err != nil {
		return domain.Review{}, err
	}
	createdReview, err := r.GetOneById(tx, review.Id.String())
	return createdReview, nil
}

func (r *ReviewRepositoryImpl) GetOneById(tx *sql.Tx, id string) (domain.Review, error) {
	var review domain.Review
	err := tx.QueryRow("SELECT id, rating, comment, user_id, menu_id, order_id, created_at, updated_at FROM review WHERE id = ? AND deleted_at IS NULL AND deleted = false", id).Scan(&review.Id, &review.Rating, &review.Comment, &review.UserId, &review.MenuId, &review.OrderId, &review.CreatedAt, &review.UpdatedAt)

	if err != nil {
		return domain.Review{}, err
	}

	return review, nil
}

func (r *ReviewRepositoryImpl) Update(tx *sql.Tx, review domain.Review) (domain.Review, error) {
	_, err := tx.Exec("UPDATE review SET rating = ?, comment = ?, updated_at = ? WHERE id = ?", review.Rating, review.Comment, review.UpdatedAt, review.Id)
	if err != nil {
		return domain.Review{}, err
	}
	updatedReview, err := r.GetOneById(tx, review.Id.String())
	return updatedReview, nil
}

func (r *ReviewRepositoryImpl) CheckReviewedItem(tx *sql.Tx, userId, menuId, orderId string) bool {
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM review WHERE user_id = ? AND menu_id = ? AND order_id = ?", userId, menuId, orderId).Scan(&count)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to check reviewed item")
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func (r *ReviewRepositoryImpl) CountReviewByMenuId(tx *sql.Tx, menuId string) (int, float64, error) {
	var totalReviews int
	var totalRating sql.NullFloat64

	query := "SELECT COUNT(*) AS total_reviews, COALESCE(SUM(rating), 0) AS total_rating FROM review WHERE menu_id = ?"
	err := tx.QueryRow(query, menuId).Scan(&totalReviews, &totalRating)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get review stats by menu id")
		return 0, 0, nil
	}

	return totalReviews, totalRating.Float64, nil
}
