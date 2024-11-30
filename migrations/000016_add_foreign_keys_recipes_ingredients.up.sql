ALTER TABLE recipes_ingredients ADD CONSTRAINT fk_recipes_ingredients_recipe FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_recipes_ingredients_ingredient FOREIGN KEY (ingredient_id) REFERENCES ingredients (id) ON DELETE CASCADE ON UPDATE CASCADE;
