package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type ReservationRepositoryImpl struct {
}

func NewReservationRepository() ReservationRepository {
	return &ReservationRepositoryImpl{}
}

func (r *ReservationRepositoryImpl) GetAll(tx *sql.Tx) ([]domain.Reservation, error) {
	reservations := []domain.Reservation{}
	rows, err := tx.Query("SELECT id,table_id,user_id,reservation_date,reservation_time,number_of_guests,status,created_at,updated_at FROM reservations WHERE deleted = false AND deleted_at IS NULL")
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

func (r *ReservationRepositoryImpl) GetOneById(tx *sql.Tx, id string) (domain.Reservation, error) {
	reservation := domain.Reservation{}
	err := tx.QueryRow("SELECT id,table_id,user_id,reservation_date,reservation_time,number_of_guests,status,created_at,updated_at FROM reservations WHERE id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&reservation.Id, &reservation.TableId, &reservation.UserId, &reservation.ReservationDate, &reservation.ReservationTime, &reservation.NumberOfGuests, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
	if err != nil {
		return domain.Reservation{}, err
	}
	return reservation, nil
}

func (r *ReservationRepositoryImpl) GetOneByTableId(tx *sql.Tx, id string) (domain.Reservation, error) {
	reservation := domain.Reservation{}
	err := tx.QueryRow("SELECT id,table_id,user_id,reservation_date,reservation_time,number_of_guests,status,created_at,updated_at FROM reservations WHERE table_id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&reservation.Id, &reservation.TableId, &reservation.UserId, &reservation.ReservationDate, &reservation.ReservationTime, &reservation.NumberOfGuests, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
	if err != nil {
		return domain.Reservation{}, err
	}
	return reservation, nil
}

func (r *ReservationRepositoryImpl) Create(tx *sql.Tx, reservation domain.Reservation) error {
	log.Println(reservation)
	_, err := tx.Exec("INSERT INTO reservations (id,table_id,user_id,reservation_date,reservation_time,number_of_guests) VALUES (?, ?, ?, ?, ?, ?)",
		reservation.Id, reservation.TableId, reservation.UserId, reservation.ReservationDate, reservation.ReservationTime, reservation.NumberOfGuests)
	if err != nil {
		logger.Log.WithError(err)
		return err
	}
	return nil
}

func (r *ReservationRepositoryImpl) Update(tx *sql.Tx, id string, reservation domain.Reservation) error {
	_, err := tx.Exec("UPDATE reservations SET reservation_date = ?, reservation_time = ?, status = ?, number_of_guests = ? WHERE id = ?",
		reservation.ReservationDate, reservation.ReservationTime, reservation.Status, reservation.NumberOfGuests, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepositoryImpl) Delete(tx *sql.Tx, id string) error {
	_, err := tx.Exec("UPDATE reservations SET deleted = true, deleted_at = ? WHERE id = ?", time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepositoryImpl) GetDeleted(tx *sql.Tx, id string) (domain.Reservation, error) {
	reservation := domain.Reservation{}
	err := tx.QueryRow("SELECT id,table_id,user_id,reservation_date,reservation_time,number_of_guests,status,created_at,updated_at FROM reservations WHERE id = ? AND deleted = true AND deleted_at IS NOT NULL", id).Scan(&reservation.Id, &reservation.TableId, &reservation.UserId, &reservation.ReservationDate, &reservation.ReservationTime, &reservation.NumberOfGuests, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
	if err != nil {
		return domain.Reservation{}, err
	}
	return reservation, nil
}

func (r *ReservationRepositoryImpl) Restore(tx *sql.Tx, id string) error {
	_, err := tx.Exec("UPDATE reservations SET deleted = ?, deleted_at = ? WHERE id = ?", false, nil, id)
	if err != nil {
		return err
	}
	return nil
}
