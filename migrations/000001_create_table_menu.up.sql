CREATE TABLE IF NOT EXISTS menu (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price INT NOT NULL,
    category ENUM(
            'main',          -- Makanan Utama
            'appetizer',     -- Pembuka
            'dessert',       -- Penutup
            'drink',         -- Minuman
            'snack',         -- Camilan
            'vegetarian',    -- Vegetarian/Vegan
            'kids',          -- Hidangan Anak
            'local',         -- Hidangan Lokal
            'special',       -- Spesial Chef
            'combo',         -- Paket
            'breakfast',     -- Sarapan
            'healthy',       -- Makanan Sehat
            'international', -- Hidangan Internasional
            'seafood',       -- Hidangan Laut
            'spicy'          -- Hidangan Pedas
        ) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    rating INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP DEFAULT NULL
);
