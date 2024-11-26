package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type ReviewRepository interface {
	GetAllByMenuId(tx *sql.Tx, id string) ([]domain.Review, error)
	Create(tx *sql.Tx, review domain.Review) (domain.Review, error)
	GetOneById(tx *sql.Tx, id string) (domain.Review, error)
	Update(tx *sql.Tx, review domain.Review) (domain.Review, error)
}
