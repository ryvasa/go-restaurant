package repository

import (
	"context"
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type ReviewRepositoryImpl struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) ReviewRepository {
	return &ReviewRepositoryImpl{db}
}

func (r *ReviewRepositoryImpl) GetAllByMenuId(ctx context.Context, id string) ([]domain.Review, error) {
	reviews := []domain.Review{}
	rows, err := r.db.QueryContext(ctx, "SELECT id,rating,comment,user_id,menu_id,created_at,updated_at FROM review WHERE menu_id = ? AND deleted_at IS NULL AND deleted = false", id)
	if err != nil {
		return nil, err
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

func (r *ReviewRepositoryImpl) Create(ctx context.Context, review domain.Review) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO review (rating, comment, user_id, menu_id) VALUES (?, ?, ?, ?)", review.Rating, review.Comment, review.UserId, review.MenuId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewRepositoryImpl) GetOneById(ctx context.Context, id string) (domain.Review, error) {
	var review domain.Review
	err := r.db.QueryRowContext(ctx, "SELECT * FROM review WHERE id = ? AND deleted_at IS NULL AND deleted = false", id).Scan(&review.Id, &review.Rating, &review.Comment, &review.UserId, &review.MenuId, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}

func (r *ReviewRepositoryImpl) Update(ctx context.Context, review domain.Review) error {
	_, err := r.db.ExecContext(ctx, "UPDATE review SET rating = $1, comment = $2, updated_at = $3 WHERE id = $4", review.Rating, review.Comment, review.UpdatedAt, review.Id)
	if err != nil {
		return err
	}
	return nil
}
