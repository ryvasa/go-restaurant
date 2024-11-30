ALTER TABLE inventory ADD CONSTRAINT fk_inventory_ingredients_ingredient FOREIGN KEY (ingredient_id) REFERENCES ingredients (id) ON DELETE CASCADE ON UPDATE CASCADE;
