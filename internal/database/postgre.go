package database

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresConfig struct {
	host     string
	user     string
	password string
	database string
	port     string
	sslMode  string
	timeZone string
}

func ConnectToPostgres() (*gorm.DB, error) {
	// Load reads .env for local development without overwriting real environment
	// variables supplied by a deployment environment.
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	config, err := loadPostgresConfig()
	if err != nil {
		return nil, err
	}

	dsn := config.dsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	return db, nil
}

func loadPostgresConfig() (postgresConfig, error) {
	config := postgresConfig{
		host:     os.Getenv("DB_HOST"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		database: os.Getenv("DB_NAME"),
		port:     os.Getenv("DB_PORT"),
		sslMode:  os.Getenv("DB_SSLMODE"),
		timeZone: os.Getenv("DB_TIMEZONE"),
	}

	for key, value := range map[string]string{
		"DB_HOST":     config.host,
		"DB_USER":     config.user,
		"DB_PASSWORD": config.password,
		"DB_NAME":     config.database,
		"DB_PORT":     config.port,
		"DB_SSLMODE":  config.sslMode,
		"DB_TIMEZONE": config.timeZone,
	} {
		if strings.TrimSpace(value) == "" {
			return postgresConfig{}, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	return config, nil
}

func (c postgresConfig) dsn() string {
	query := url.Values{}
	query.Set("sslmode", c.sslMode)
	query.Set("TimeZone", c.timeZone)

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.user, c.password),
		Host:     net.JoinHostPort(c.host, c.port),
		Path:     "/" + c.database,
		RawQuery: query.Encode(),
	}).String()
}
