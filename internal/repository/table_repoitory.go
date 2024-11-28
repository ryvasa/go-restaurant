package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type TableRepository interface {
	GetAll(tx *sql.Tx) ([]domain.Table, error)
	GetOneById(tx *sql.Tx, id string) (domain.Table, error)
	Create(tx *sql.Tx, table domain.Table) error
	Update(tx *sql.Tx, id string, table domain.Table) error
	Delete(tx *sql.Tx, id string) error
	GetDeleted(tx *sql.Tx, id string) (domain.Table, error)
	Restore(tx *sql.Tx, id string) error
}
