ALTER TABLE order_menu
    ADD CONSTRAINT fk_order_menu_order
    FOREIGN KEY (order_id)
    REFERENCES orders(id)
    ON DELETE CASCADE,

    ADD CONSTRAINT fk_order_menu_menu
    FOREIGN KEY (menu_id)
    REFERENCES menu(id)
    ON DELETE CASCADE;
