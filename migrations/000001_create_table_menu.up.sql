CREATE TABLE IF NOT EXISTS menu (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price FLOAT NOT NULL,
    category ENUM(
            'main',
            'appetizer',
            'dessert',
            'drink',
            'snack',
            'vegetarian',
            'kids',
            'local',
            'special',
            'combo',
            'breakfast',
            'healthy',
            'international',
            'seafood',
            'spicy'
            ) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    rating FLOAT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP DEFAULT NULL
);
