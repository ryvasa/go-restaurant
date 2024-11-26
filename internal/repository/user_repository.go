package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type UserRepository interface {
	GetAll(tx *sql.Tx) ([]domain.User, error)
	Create(tx *sql.Tx, user domain.User) (domain.User, error)
	Get(tx *sql.Tx, id string) (domain.User, error)
	Update(tx *sql.Tx, user domain.User) (domain.User, error)
	GetByEmail(tx *sql.Tx, email string) (domain.User, error)
	Delete(tx *sql.Tx, id string) error
	Restore(tx *sql.Tx, id string) (domain.User, error)
	GetDeletedUserById(tx *sql.Tx, id string) ([]domain.User, error)
}
