ALTER TABLE reservations
    ADD CONSTRAINT fk_reservation_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    ADD CONSTRAINT fk_reservation_table
    FOREIGN KEY (table_id)
    REFERENCES tables(id)
    ON DELETE CASCADE;
