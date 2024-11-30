package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type DB interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type TransactionRepositoryImpl struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &TransactionRepositoryImpl{
		db: db,
	}
}

type Adapters struct {
	UserRepository             UserRepository
	MenuRepository             MenuRepository
	TableRepository            TableRepository
	ReservationRepository      ReservationRepository
	OrderRepository            OrderRepository
	OrderMenuRepository        OrderMenuRepository
	ReviewRepository           ReviewRepository
	RecipeRepository           RecipeRepository
	IngredientRepository       IngredientRepository
	RecipeIngredientRepository RecipeIngredientRepository
	InventoryRepository        InventoryRepository
}

func (p *TransactionRepositoryImpl) Transact(txFunc func(adapters Adapters) error) error {
	return runInTx(p.db, func(tx *sql.Tx) error {
		adapters := Adapters{
			UserRepository:             NewUserRepository(tx),
			MenuRepository:             NewMenuRepository(tx),
			TableRepository:            NewTableRepository(tx),
			ReservationRepository:      NewReservationRepository(tx),
			OrderRepository:            NewOrderRepository(tx),
			OrderMenuRepository:        NewOrderMenuRepository(tx),
			ReviewRepository:           NewReviewRepository(tx),
			RecipeRepository:           NewRecipeRepository(tx),
			IngredientRepository:       NewIngredientRepository(tx),
			RecipeIngredientRepository: NewRecipeIngredientRepository(tx),
			InventoryRepository:        NewInventoryRepository(tx),
		}

		return txFunc(adapters)
	})
}

func runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	logger.Log.Info("Starting transaction")
	if err != nil {
		logger.Log.WithError(err).Error("Error when starting transaction")
		return err
	}

	err = fn(tx)
	if err == nil {
		logger.Log.Info("Committing transaction")
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		logger.Log.Error("Error while transaction, rollback")
		return errors.Join(err, rollbackErr)
	}

	return err
}
