package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type ReservationUsecaseImpl struct {
	reservationRepo repository.ReservationRepository
	tableRepo       repository.TableRepository
	txRepo          repository.TransactionRepository
}

func NewReservationUsecase(reservationRepo repository.ReservationRepository, tableRepo repository.TableRepository, txRepo repository.TransactionRepository) ReservationUsecase {
	return &ReservationUsecaseImpl{
		reservationRepo,
		tableRepo,
		txRepo,
	}
}

func (u *ReservationUsecaseImpl) getReservationById(ctx context.Context, repo repository.ReservationRepository, id uuid.UUID) (domain.Reservation, error) {
	reservation, err := repo.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error getting reservation after creation/update")
		return domain.Reservation{}, utils.NewInternalError("Failed to get reservation")
	}
	return reservation, nil
}

func (u *ReservationUsecaseImpl) GetAll(ctx context.Context) ([]domain.Reservation, error) {
	reservations, err := u.reservationRepo.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all reservations")
		return []domain.Reservation{}, utils.NewInternalError("Failed to get all reservations")
	}
	return reservations, nil

}

func (u *ReservationUsecaseImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Reservation, error) {
	reservation, err := u.reservationRepo.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error reservation not found")
		return domain.Reservation{}, utils.NewNotFoundError("Reservation not found")
	}
	return reservation, nil

}

func (u *ReservationUsecaseImpl) Create(ctx context.Context, req dto.CreateReservationRequest) (domain.Reservation, error) {
	result := domain.Reservation{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		tableId, err := uuid.Parse(req.TableId)
		if err != nil {
			logger.Log.WithError(err).Error("Error invalid table id format")
			return utils.NewValidationError("Invalid table id format")
		}

		userId, err := uuid.Parse(req.UserId)
		if err != nil {
			logger.Log.WithError(err).Error("Error invalid user id format")
			return utils.NewValidationError("Invalid user id format")
		}

		_, err = adapters.TableRepository.GetOneById(ctx, tableId)
		if err != nil {
			logger.Log.WithError(err).Error("Error table not found")
			return utils.NewNotFoundError("Table not found")
		}

		reservationDate, err := time.Parse("2006-01-02", req.ReservationDate)
		if err != nil {
			logger.Log.WithField("input_date", req.ReservationDate).WithError(err).Error("Invalid reservation date format")
			return utils.NewValidationError("Invalid reservation date format")
		}

		reservationTime, err := time.Parse("15:04:05", req.ReservationTime)
		if err != nil {
			logger.Log.WithField("input_time", req.ReservationTime).WithError(err).Error("Invalid reservation time format")
			return utils.NewValidationError("Invalid reservation time format")
		}

		if reservationDate.IsZero() || reservationTime.IsZero() {
			logger.Log.Error("Reservation date or time has default zero value")
			return utils.NewValidationError("Invalid reservation date or time")
		}

		existingReservation, err := adapters.ReservationRepository.GetOneByTableId(ctx, tableId)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get reservation")
			return utils.NewInternalError("Failed to get reservation")
		}

		if existingReservation.Status == "confirmed" && existingReservation.ReservationDate.Equal(reservationDate) {
			existingTime, _ := time.Parse("15:04:05", existingReservation.ReservationTime)
			newTime, _ := time.Parse("15:04:05", reservationTime.Format("15:04:05"))
			timeDiff := newTime.Sub(existingTime).Hours()

			if timeDiff < 2 && timeDiff > -2 {
				logger.Log.Error("Reservation time conflicts with existing reservation")
				return utils.NewValidationError("Cannot make a reservation within 2 hours of an existing confirmed reservation")
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

		err = adapters.ReservationRepository.Create(ctx, reservation)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to create reservation")
			return utils.NewInternalError("Failed to create reservation")
		}

		createdreservation, err := adapters.ReservationRepository.GetOneById(ctx, reservation.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get reservation")
			return utils.NewInternalError("Failed to get reservation")
		}
		result = createdreservation
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *ReservationUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateReservationRequest) (domain.Reservation, error) {
	result := domain.Reservation{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {

		existingReservation, err := adapters.ReservationRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error reservation not found")
			return utils.NewNotFoundError("Reservation not found")
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
				return utils.NewValidationError("Invalid reservation date format")
			}
			reservation.ReservationDate = reservationDate

		}

		if req.ReservationTime == "" {
			reservation.ReservationTime = existingReservation.ReservationTime
		} else {
			reservationTime, err := time.Parse("15:04:05", req.ReservationTime)
			if err != nil {
				logger.Log.WithError(err).Error("Invalid reservation time format")
				return utils.NewValidationError("Invalid reservation time format")
			}
			reservation.ReservationTime =
				reservationTime.Format("15:04:05")
		}

		err = adapters.ReservationRepository.Update(ctx, id, reservation)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update reservation")
			return utils.NewInternalError("Failed to update reservation")
		}

		updatedReservation, err := adapters.ReservationRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get reservation")
			return utils.NewInternalError("Failed to get reservation")
		}
		result = updatedReservation
		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil

}

func (u *ReservationUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {

		if _, err := adapters.ReservationRepository.GetOneById(ctx, id); err != nil {
			logger.Log.WithError(err).Error("Error reservation not found")
			return utils.NewNotFoundError("Reservation not found")
		}

		err := adapters.ReservationRepository.Delete(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to delete reservation")
			return utils.NewInternalError("Failed to delete reservation")
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil

}

func (u *ReservationUsecaseImpl) Restore(ctx context.Context, id uuid.UUID) (domain.Reservation, error) {
	result := domain.Reservation{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {

		_, err := adapters.ReservationRepository.GetDeleted(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error reservation not found to restore")
			return utils.NewNotFoundError("Reservation not found to restore")
		}

		err = adapters.ReservationRepository.Restore(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to restore reservation")
			return utils.NewInternalError("Failed to restore reservation")
		}

		restoredReservation, err := adapters.ReservationRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get reservation")
			return utils.NewInternalError("Failed to get reservation")
		}
		result = restoredReservation
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
