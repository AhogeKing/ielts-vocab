package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	DBTimeZone string

	JWTSecret string
	JWTIssuer string
	JWTTTL    time.Duration
}

func LoadConfig() (*Config, error) {
	// A .env file is convenient for local development. Deployment environments
	// can provide these variables directly, so a missing .env file is allowed.
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	jwtTTL, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		return nil, fmt.Errorf("parse JWT_TTL: %w", err)
	}

	config := &Config{
		Port:       os.Getenv("PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		DBTimeZone: os.Getenv("DB_TIMEZONE"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		JWTIssuer:  os.Getenv("JWT_ISSUER"),
		JWTTTL:     jwtTTL,
	}
	if err := config.validateJWT(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validateJWT() error {
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}
	if strings.TrimSpace(c.JWTIssuer) == "" {
		return fmt.Errorf("JWT_ISSUER must be set")
	}
	if c.JWTTTL <= 0 {
		return fmt.Errorf("JWT_TTL must be positive")
	}
	return nil
}
