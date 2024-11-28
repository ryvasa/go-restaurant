package usecase

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
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

	reservations, err := u.reservationRepo.GetAll(tx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all reservations")
		return []domain.Reservation{}, utils.NewInternalError("Failed to get all reservations")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return []domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return reservations, nil

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

	reservation, err := u.reservationRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error reservation not found")
		return domain.Reservation{}, utils.NewNotFoundError("Reservation not found")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return reservation, nil

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

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Reservation{}, utils.NewValidationError(err)
	}

	tableId, err := uuid.Parse(req.TableId)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid table id format")
		return domain.Reservation{}, utils.NewValidationError("Invalid table id format")
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid user id format")
		return domain.Reservation{}, utils.NewValidationError("Invalid user id format")
	}
	reservationDate, err := time.Parse("2006-01-02", req.ReservationDate)
	if err != nil {
		logger.Log.WithField("input_date", req.ReservationDate).WithError(err).Error("Invalid reservation date format")
		return domain.Reservation{}, utils.NewValidationError("Invalid reservation date format")
	}

	reservationTime, err := time.Parse("15:04:05", req.ReservationTime)
	if err != nil {
		logger.Log.WithField("input_time", req.ReservationTime).WithError(err).Error("Invalid reservation time format")
		return domain.Reservation{}, utils.NewValidationError("Invalid reservation time format")
	}

	if reservationDate.IsZero() || reservationTime.IsZero() {
		logger.Log.Error("Reservation date or time has default zero value")
		return domain.Reservation{}, utils.NewValidationError("Invalid reservation date or time")
	}

	existingReservation, err := u.reservationRepo.GetOneByTableId(tx, tableId.String())
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get reservation")
		return domain.Reservation{}, utils.NewInternalError("Failed to get reservation")
	}

	if existingReservation.Status == "confirmed" && existingReservation.ReservationDate.Equal(reservationDate) {
		// Hitung selisih waktu antara reservasi yang ada dan reservasi baru
		existingTime, _ := time.Parse("15:04:05", existingReservation.ReservationTime)
		newTime, _ := time.Parse("15:04:05", reservationTime.Format("15:04:05"))
		log.Println(existingTime, newTime)
		// Cari selisih waktu dalam jam
		timeDiff := newTime.Sub(existingTime).Hours()

		// Cek apakah selisih waktu kurang dari 2 jam
		if timeDiff < 2 && timeDiff > -2 {
			logger.Log.Error("Reservation time conflicts with existing reservation")
			return domain.Reservation{}, utils.NewValidationError("Cannot make a reservation within 2 hours of an existing confirmed reservation")
		}
	}

	reservation := domain.Reservation{
		Id:              uuid.New(),
		UserId:          userId,
		TableId:         tableId,
		ReservationDate: reservationDate,
		ReservationTime: reservationTime.Format("15:04:05"),
		NumberOfGuests:  req.NumberOfGuests,
	}

	err = u.reservationRepo.Create(tx, reservation)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create reservation")
		return domain.Reservation{}, utils.NewInternalError("Failed to create reservation")
	}

	createdreservation, err := u.reservationRepo.GetOneById(tx, reservation.Id.String())
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get reservation")
		return domain.Reservation{}, utils.NewInternalError("Failed to get reservation")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return createdreservation, nil

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

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return domain.Reservation{}, utils.NewValidationError("Invalid ID format")
	}

	existingReservation, err := u.reservationRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error reservation not found")
		return domain.Reservation{}, utils.NewNotFoundError("Reservation not found")
	}

	reservation := domain.Reservation{
		NumberOfGuests: req.NumberOfGuests,
		Status:         req.Status,
	}
	if req.Status == "" {
		reservation.Status = existingReservation.Status
	}
	if req.NumberOfGuests == 0 {
		reservation.NumberOfGuests = existingReservation.NumberOfGuests
	}

	if req.ReservationDate == "" {
		reservation.ReservationDate = existingReservation.ReservationDate
	} else {
		reservationDate, err := time.Parse("2006-01-02", req.ReservationDate)
		if err != nil {
			logger.Log.WithError(err).Error("Invalid reservation date format")
			return domain.Reservation{}, utils.NewValidationError("Invalid reservation date format")
		}
		reservation.ReservationDate = reservationDate

	}

	if req.ReservationTime == "" {
		reservation.ReservationTime = existingReservation.ReservationTime
	} else {
		reservationTime, err := time.Parse("15:04:05", req.ReservationTime)
		if err != nil {
			logger.Log.WithError(err).Error("Invalid reservation time format")
			return domain.Reservation{}, utils.NewValidationError("Invalid reservation time format")
		}
		reservation.ReservationTime =
			reservationTime.Format("15:04:05")
	}

	err = u.reservationRepo.Update(tx, id, reservation)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update reservation")
		return domain.Reservation{}, utils.NewInternalError("Failed to update reservation")
	}

	updatedReservation, err := u.reservationRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get reservation")
		return domain.Reservation{}, utils.NewInternalError("Failed to get reservation")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return updatedReservation, nil

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

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return utils.NewValidationError("Invalid ID format")
	}

	if _, err := u.reservationRepo.GetOneById(tx, id); err != nil {
		logger.Log.WithError(err).Error("Error reservation not found")
		return utils.NewNotFoundError("Reservation not found")
	}

	err = u.reservationRepo.Delete(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete reservation")
		return utils.NewInternalError("Failed to delete reservation")
	}

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

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return domain.Reservation{}, utils.NewValidationError("Invalid ID format")
	}

	_, err = u.reservationRepo.GetDeleted(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error reservation not found to restore")
		return domain.Reservation{}, utils.NewNotFoundError("Reservation not found to restore")
	}

	err = u.reservationRepo.Restore(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore reservation")
		return domain.Reservation{}, utils.NewInternalError("Failed to restore reservation")
	}

	restoredReservation, err := u.reservationRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get reservation")
		return domain.Reservation{}, utils.NewInternalError("Failed to get reservation")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Reservation{}, utils.NewInternalError("Failed to commit transaction")
	}
	return restoredReservation, nil

}
