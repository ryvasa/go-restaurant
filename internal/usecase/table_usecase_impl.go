package usecase

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
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

	table, err := u.tableRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error table not found")
		return domain.Table{}, utils.NewNotFoundError("Table not found")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Table{}, utils.NewInternalError("Failed to commit transaction")
	}
	return table, nil
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

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Table{}, utils.NewValidationError(err)
	}

	table := domain.Table{
		Id:       uuid.New(),
		Number:   req.Number,
		Capacity: req.Capacity,
		Location: req.Location,
	}

	err = u.tableRepo.Create(tx, table)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create table")
		return domain.Table{}, utils.NewInternalError("Failed to create table")
	}

	createdTable, err := u.tableRepo.GetOneById(tx, table.Id.String())
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get table")
		return domain.Table{}, utils.NewInternalError("Failed to get table")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Table{}, utils.NewInternalError("Failed to commit transaction")
	}
	return createdTable, nil
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

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return domain.Table{}, utils.NewValidationError("Invalid ID format")
	}

	existingTable, err := u.tableRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error table not found")
		return domain.Table{}, utils.NewNotFoundError("Table not found")
	}

	table := domain.Table{
		Number:   req.Number,
		Capacity: req.Capacity,
		Location: req.Location,
		Status:   req.Status,
	}

	if table.Number == "" {
		table.Number = existingTable.Number
	}
	if table.Capacity == 0 {
		table.Capacity = existingTable.Capacity
	}
	if table.Location == "" {
		table.Location = existingTable.Location
	}
	if table.Status == "" {
		table.Status = existingTable.Status
	}
	log.Println(table)

	err = u.tableRepo.Update(tx, id, table)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update table")
		return domain.Table{}, utils.NewInternalError("Failed to update table")
	}

	updatedTable, err := u.tableRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get table")
		return domain.Table{}, utils.NewInternalError("Failed to get table")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Table{}, utils.NewInternalError("Failed to commit transaction")
	}
	return updatedTable, nil
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

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return utils.NewValidationError("Invalid ID format")
	}

	if _, err := u.tableRepo.GetOneById(tx, id); err != nil {
		logger.Log.WithError(err).Error("Error table not found")
		return utils.NewNotFoundError("Table not found")
	}

	err = u.tableRepo.Delete(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete table")
		return utils.NewInternalError("Failed to delete table")
	}

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

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return domain.Table{}, utils.NewValidationError("Invalid ID format")
	}

	_, err = u.tableRepo.GetDeleted(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error table not found to restore")
		return domain.Table{}, utils.NewNotFoundError("Table not found to restore")
	}

	err = u.tableRepo.Restore(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore table")
		return domain.Table{}, utils.NewInternalError("Failed to restore table")
	}

	restoredTable, err := u.tableRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get table")
		return domain.Table{}, utils.NewInternalError("Failed to get table")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Table{}, utils.NewInternalError("Failed to commit transaction")
	}
	return restoredTable, nil
}
