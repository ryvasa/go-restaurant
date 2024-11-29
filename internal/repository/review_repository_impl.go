package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type ReviewRepositoryImpl struct {
	db DB
}

func NewReviewRepository(db DB) ReviewRepository {
	return &ReviewRepositoryImpl{db}
}

func (r *ReviewRepositoryImpl) GetAllByMenuId(ctx context.Context, menuId uuid.UUID) ([]domain.Review, error) {
	reviews := []domain.Review{}
	rows, err := r.db.QueryContext(ctx, "SELECT id,rating,comment,user_id,menu_id,order_id,created_at,updated_at FROM review WHERE menu_id = ? AND deleted_at IS NULL AND deleted = false", menuId)
	if err != nil {
		return []domain.Review{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var review domain.Review
		err := rows.Scan(&review.Id, &review.Rating, &review.Comment, &review.UserId, &review.MenuId, &review.OrderId, &review.UpdatedAt, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *ReviewRepositoryImpl) Create(ctx context.Context, review domain.Review) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO review (id,rating, comment, user_id, menu_id, order_id) VALUES (?, ?, ?, ?, ?, ?)", review.Id, review.Rating, review.Comment, review.UserId, review.MenuId, review.OrderId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewRepositoryImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Review, error) {
	var review domain.Review
	err := r.db.QueryRowContext(ctx, "SELECT id, rating, comment, user_id, menu_id, order_id, created_at, updated_at FROM review WHERE id = ? AND deleted_at IS NULL AND deleted = false", id).Scan(&review.Id, &review.Rating, &review.Comment, &review.UserId, &review.MenuId, &review.OrderId, &review.CreatedAt, &review.UpdatedAt)

	if err != nil {
		return domain.Review{}, err
	}

	return review, nil
}

func (r *ReviewRepositoryImpl) Update(ctx context.Context, id uuid.UUID, review domain.Review) error {
	_, err := r.db.ExecContext(ctx, "UPDATE review SET rating = ?, comment = ? WHERE id = ?", review.Rating, review.Comment, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewRepositoryImpl) CheckReviewedItem(ctx context.Context, userId, menuId, orderId uuid.UUID) bool {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM review WHERE user_id = ? AND menu_id = ? AND order_id = ?", userId, menuId, orderId).Scan(&count)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to check reviewed item")
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func (r *ReviewRepositoryImpl) CountReviewByMenuId(ctx context.Context, menuId uuid.UUID) (int, float64, error) {
	var totalReviews int
	var totalRating sql.NullFloat64

	query := "SELECT COUNT(*) AS total_reviews, COALESCE(SUM(rating), 0) AS total_rating FROM review WHERE menu_id = ?"
	err := r.db.QueryRowContext(ctx, query, menuId).Scan(&totalReviews, &totalRating)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get review stats by menu id")
		return 0, 0, nil
	}

	return totalReviews, totalRating.Float64, nil
}
