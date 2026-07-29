package database

import (
	"fmt"
	"ielts-vocab/internal/config"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectToPostgres opens a PostgreSQL connection using the loaded database
// configuration.
func ConnectToPostgres(c *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseDSN(c)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	if err = validateDatabase(c); err != nil {
		return nil, err
	}

	return db, nil
}

func validateDatabase(c *config.Config) error {
	for key, value := range map[string]string{
		"DB_HOST":     c.DBHost,
		"DB_USER":     c.DBUser,
		"DB_PASSWORD": c.DBPassword,
		"DB_NAME":     c.DBName,
		"DB_PORT":     c.DBPort,
		"DB_SSLMODE":  c.DBSSLMode,
		"DB_TIMEZONE": c.DBTimeZone,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("required environment variable %s is not set", key)
		}
	}
	return nil
}

func databaseDSN(c *config.Config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		quoteDSNValue(c.DBHost),
		quoteDSNValue(c.DBUser),
		quoteDSNValue(c.DBPassword),
		quoteDSNValue(c.DBName),
		quoteDSNValue(c.DBPort),
		quoteDSNValue(c.DBSSLMode),
		c.DBTimeZone,
	)
}

// quoteDSNValue escapes a value for PostgreSQL's key=value connection-string
// format. This keeps passwords containing spaces, quotes, or backslashes valid.
func quoteDSNValue(value string) string {
	escaped := strings.NewReplacer(
		"\\\\", "\\\\\\\\",
		"'", "\\'",
	).Replace(value)
	return "'" + escaped + "'"
}
