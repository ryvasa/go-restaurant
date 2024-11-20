package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		User     string
		Password string
		Name     string
	}
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	// Load Server Config
	config.Server.Port = os.Getenv("SERVER_PORT")

	// Load Database Config
	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Name = os.Getenv("DB_NAME")

	return config, nil
}
