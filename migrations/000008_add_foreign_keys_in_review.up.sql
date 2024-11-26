ALTER TABLE review
    ADD CONSTRAINT fk_review_menu
    FOREIGN KEY (menu_id)
    REFERENCES menu(id)
    ON DELETE CASCADE,

    ADD CONSTRAINT fk_review_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    ADD CONSTRAINT fk_review_order
    FOREIGN KEY (order_id)
    REFERENCES orders(id)
    ON DELETE CASCADE;