ALTER TABLE reservations
    DROP CONSTRAINT fk_reservation_user,
    DROP CONSTRAINT fk_reservation_table;
