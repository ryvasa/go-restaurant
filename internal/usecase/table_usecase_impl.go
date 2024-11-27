package usecase

import (
	"context"
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type TableUsecaseImpl struct {
	db        *sql.DB
	tableRepo repository.TableRepository
}

func NewTableUsecase(db *sql.DB, tableRepo repository.TableRepository) TableUsecase {
	return &TableUsecaseImpl{
		db,
		tableRepo,
	}
}
func (u *TableUsecaseImpl) GetAll(ctx context.Context) ([]domain.Table, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return []domain.Table{}, err
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	tables, err := u.tableRepo.GetAll(tx)
	if err != nil {
		logger.Log.WithError(err).Error("Error when getting all tables")
		return []domain.Table{}, utils.NewInternalError("Failed to get all tables")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return []domain.Table{}, utils.NewInternalError("Failed to commit transaction")
	}

	return tables, nil
}

func (u *TableUsecaseImpl) GetOneById(ctx context.Context, id string) (domain.Table, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Table{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	// TODO
	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Table{}, utils.NewInternalError("Failed to commit transaction")
	}
	return domain.Table{}, nil
}

func (u *TableUsecaseImpl) Create(ctx context.Context, req dto.CreateTableRequest) (domain.Table, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Table{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	// TODO

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Table{}, utils.NewInternalError("Failed to commit transaction")
	}
	return domain.Table{}, nil
}

func (u *TableUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateTableRequest) (domain.Table, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Table{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	// TODO

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Table{}, utils.NewInternalError("Failed to commit transaction")
	}
	return domain.Table{}, nil
}

func (u *TableUsecaseImpl) Delete(ctx context.Context, id string) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	// TODO
	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return utils.NewInternalError("Failed to commit transaction")
	}
	return nil
}

func (u *TableUsecaseImpl) Restore(ctx context.Context, id string) (domain.Table, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Table{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	// TODO
	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Table{}, utils.NewInternalError("Failed to commit transaction")
	}
	return domain.Table{}, nil
}
