CREATE TABLE recipes (
    id CHAR(36) PRIMARY KEY,
    menu_id CHAR(36) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP DEFAULT NULL
);
