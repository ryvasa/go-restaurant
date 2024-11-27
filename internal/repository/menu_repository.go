package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type MenuRepository interface {
	GetAll(tx *sql.Tx) ([]domain.Menu, error)
	Create(tx *sql.Tx, menu domain.Menu) (domain.Menu, error)
	Get(tx *sql.Tx, id string) (domain.Menu, error)
	Update(tx *sql.Tx, menu domain.Menu) (domain.Menu, error)
	Delete(tx *sql.Tx, id string) error
	Restore(tx *sql.Tx, id string) (domain.Menu, error)
	GetDeletedMenuById(tx *sql.Tx, id string) (domain.Menu, error)
	UpdateRating(tx *sql.Tx, id string, rating float64) error
}
