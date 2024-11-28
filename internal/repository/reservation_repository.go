package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type ReservationRepository interface {
	GetAll(tx *sql.Tx) ([]domain.Reservation, error)
	GetOneById(tx *sql.Tx, id string) (domain.Reservation, error)
	GetOneByTableId(tx *sql.Tx, tableId string) (domain.Reservation, error)
	Create(tx *sql.Tx, reservation domain.Reservation) error
	Update(tx *sql.Tx, id string, reservation domain.Reservation) error
	Delete(tx *sql.Tx, id string) error
	Restore(tx *sql.Tx, id string) error
	GetDeleted(tx *sql.Tx, id string) (domain.Reservation, error)
}
