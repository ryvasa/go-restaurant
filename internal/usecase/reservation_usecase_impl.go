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

type ReservationUsecaseImpl struct {
	db              *sql.DB
	reservationRepo repository.ReservationRepository
	tableRepo       repository.TableRepository
}

func NewReservationUsecase(db *sql.DB, reservationRepo repository.ReservationRepository, tableRepo repository.TableRepository) ReservationUsecase {
	return &ReservationUsecaseImpl{
		db,
		reservationRepo,
		tableRepo,
	}
}

func (u *ReservationUsecaseImpl) GetAll(ctx context.Context) ([]domain.Reservation, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return []domain.Reservation{}, utils.NewInternalError("Failed to begin transaction")
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
		return []domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return []domain.Reservation{}, nil

}

func (u *ReservationUsecaseImpl) GetOneById(ctx context.Context, id string) (domain.Reservation, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Reservation{}, utils.NewInternalError("Failed to begin transaction")
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
		return domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return domain.Reservation{}, nil

}

func (u *ReservationUsecaseImpl) Create(ctx context.Context, req dto.CreateReservationRequest) (domain.Reservation, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Reservation{}, utils.NewInternalError("Failed to begin transaction")
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
		return domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return domain.Reservation{}, nil

}

func (u *ReservationUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateReservationRequest) (domain.Reservation, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Reservation{}, utils.NewInternalError("Failed to begin transaction")
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
		return domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return domain.Reservation{}, nil

}

func (u *ReservationUsecaseImpl) Delete(ctx context.Context, id string) error {
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

func (u *ReservationUsecaseImpl) Restore(ctx context.Context, id string) (domain.Reservation, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Reservation{}, utils.NewInternalError("Failed to begin transaction")
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
		return domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return domain.Reservation{}, nil

}
