package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type TableUsecaseImpl struct {
	tableRepo repository.TableRepository
	txRepo    repository.TransactionRepository
}

func NewTableUsecase(tableRepo repository.TableRepository, txRepo repository.TransactionRepository) TableUsecase {
	return &TableUsecaseImpl{
		tableRepo,
		txRepo,
	}
}
func (u *TableUsecaseImpl) GetAll(ctx context.Context) ([]domain.Table, error) {
	tables, err := u.tableRepo.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all tables")
		return []domain.Table{}, utils.NewInternalError("Failed to get all tables")
	}
	return tables, nil
}

func (u *TableUsecaseImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Table, error) {

	table, err := u.tableRepo.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error table not found")
		return domain.Table{}, utils.NewNotFoundError("Table not found")
	}

	return table, nil
}

func (u *TableUsecaseImpl) Create(ctx context.Context, req dto.CreateTableRequest) (domain.Table, error) {
	result := domain.Table{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {

		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		table := domain.Table{
			Id:       uuid.New(),
			Number:   req.Number,
			Capacity: req.Capacity,
			Location: req.Location,
		}

		err := adapters.TableRepository.Create(ctx, table)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to create table")
			return utils.NewInternalError("Failed to create table")
		}

		createdTable, err := adapters.TableRepository.GetOneById(ctx, table.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get table")
			return utils.NewInternalError("Failed to get table")
		}
		result = createdTable
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *TableUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTableRequest) (domain.Table, error) {
	result := domain.Table{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {

		existingTable, err := adapters.TableRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error table not found")
			return utils.NewNotFoundError("Table not found")
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

		err = adapters.TableRepository.Update(ctx, id, table)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update table")
			return utils.NewInternalError("Failed to update table")
		}

		updatedTable, err := adapters.TableRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get table")
			return utils.NewInternalError("Failed to get table")
		}

		result = updatedTable
		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *TableUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if _, err := adapters.TableRepository.GetOneById(ctx, id); err != nil {
			logger.Log.WithError(err).Error("Error table not found")
			return utils.NewNotFoundError("Table not found")
		}

		err := adapters.TableRepository.Delete(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to delete table")
			return utils.NewInternalError("Failed to delete table")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *TableUsecaseImpl) Restore(ctx context.Context, id uuid.UUID) (domain.Table, error) {
	result := domain.Table{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {

		_, err := adapters.TableRepository.GetDeleted(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error table not found to restore")
			return utils.NewNotFoundError("Table not found to restore")
		}

		err = adapters.TableRepository.Restore(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to restore table")
			return utils.NewInternalError("Failed to restore table")
		}

		restoredTable, err := adapters.TableRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get table")
			return utils.NewInternalError("Failed to get table")
		}

		result = restoredTable
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
