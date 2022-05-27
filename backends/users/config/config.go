package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type (
	Database struct {
		Username string
		Password string
		DBname   string
	}
	Config struct {
		Database Database
	}
)

const (
	databaseUsername = "POSTGRES_USER"
	databasePassword = "POSTGRES_PASSWORD"
	databaseDBname   = "POSTGRES_DB"
)

var (
	ErrDBnotFound = errors.New("config: did not find configs for database")
)

func Load(files ...string) (Config, error) {
	if err := godotenv.Load(files...); err != nil {
		return Config{}, err
	}
	cfg := Config{
		Database: Database{
			Username: os.Getenv(databaseUsername),
			Password: os.Getenv(databasePassword),
			DBname:   os.Getenv(databaseDBname),
		},
	}
	if cfg.Database.Username == "" || cfg.Database.Password == "" ||
		cfg.Database.DBname == "" {
		return cfg, ErrDBnotFound
	}
	return cfg, nil
}
