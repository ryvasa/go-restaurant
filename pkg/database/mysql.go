package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ryvasa/go-restaurant/pkg/config"
)

func ProvideDSN(cfg *config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=True&charset=utf8mb4&loc=Local", // Tambahkan charset dan loc
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Name,
	)
}

func NewMySQLConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(10)                 // Set max open connections
	db.SetMaxIdleConns(5)                  // Set max idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Set connection max lifetime

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
