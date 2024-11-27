CREATE TABLE IF NOT EXISTS tables (
    id CHAR(36) PRIMARY KEY,
    number VARCHAR(255) NOT NULL,
    capacity INTEGER NOT NULL,
    location ENUM("indoor","outdoor") NOT NULL,
    status ENUM("avaible","reserved","out of service") NOT NULL DEFAULT 'avaible',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP DEFAULT NULL
);
