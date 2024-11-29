package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type ReservationRepositoryImpl struct {
	db DB
}

func NewReservationRepository(db DB) ReservationRepository {
	return &ReservationRepositoryImpl{db}
}

func (r *ReservationRepositoryImpl) GetAll(ctx context.Context) ([]domain.Reservation, error) {
	reservations := []domain.Reservation{}
	rows, err := r.db.QueryContext(ctx, "SELECT id,table_id,user_id,reservation_date,reservation_time,number_of_guests,status,created_at,updated_at FROM reservations WHERE deleted = false AND deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reservation domain.Reservation
		err := rows.Scan(&reservation.Id, &reservation.TableId, &reservation.UserId, &reservation.ReservationDate, &reservation.ReservationTime, &reservation.NumberOfGuests, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (r *ReservationRepositoryImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Reservation, error) {

	reservation := domain.Reservation{}
	err := r.db.QueryRowContext(ctx, "SELECT id,table_id,user_id,reservation_date,reservation_time,number_of_guests,status,created_at,updated_at FROM reservations WHERE id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&reservation.Id, &reservation.TableId, &reservation.UserId, &reservation.ReservationDate, &reservation.ReservationTime, &reservation.NumberOfGuests, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
	if err != nil {
		return domain.Reservation{}, err
	}
	return reservation, nil
}

func (r *ReservationRepositoryImpl) GetOneByTableId(ctx context.Context, id uuid.UUID) (domain.Reservation, error) {
	reservation := domain.Reservation{}
	err := r.db.QueryRowContext(ctx, "SELECT id,table_id,user_id,reservation_date,reservation_time,number_of_guests,status,created_at,updated_at FROM reservations WHERE table_id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&reservation.Id, &reservation.TableId, &reservation.UserId, &reservation.ReservationDate, &reservation.ReservationTime, &reservation.NumberOfGuests, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
	if err != nil {
		return domain.Reservation{}, err
	}
	return reservation, nil
}

func (r *ReservationRepositoryImpl) Create(ctx context.Context, reservation domain.Reservation) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO reservations (id,table_id,user_id,reservation_date,reservation_time,number_of_guests) VALUES (?, ?, ?, ?, ?, ?)",
		reservation.Id, reservation.TableId, reservation.UserId, reservation.ReservationDate, reservation.ReservationTime, reservation.NumberOfGuests)
	if err != nil {
		logger.Log.WithError(err)
		return err
	}
	return nil
}

func (r *ReservationRepositoryImpl) Update(ctx context.Context, id uuid.UUID, reservation domain.Reservation) error {
	_, err := r.db.ExecContext(ctx, "UPDATE reservations SET reservation_date = ?, reservation_time = ?, status = ?, number_of_guests = ? WHERE id = ?",
		reservation.ReservationDate, reservation.ReservationTime, reservation.Status, reservation.NumberOfGuests, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE reservations SET deleted = true, deleted_at = ? WHERE id = ?", time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepositoryImpl) GetDeleted(ctx context.Context, id uuid.UUID) (domain.Reservation, error) {
	reservation := domain.Reservation{}
	err := r.db.QueryRowContext(ctx, "SELECT id,table_id,user_id,reservation_date,reservation_time,number_of_guests,status,created_at,updated_at FROM reservations WHERE id = ? AND deleted = true AND deleted_at IS NOT NULL", id).Scan(&reservation.Id, &reservation.TableId, &reservation.UserId, &reservation.ReservationDate, &reservation.ReservationTime, &reservation.NumberOfGuests, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
	if err != nil {
		return domain.Reservation{}, err
	}
	return reservation, nil
}

func (r *ReservationRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE reservations SET deleted = ?, deleted_at = ? WHERE id = ?", false, nil, id)
	if err != nil {
		return err
	}
	return nil
}
