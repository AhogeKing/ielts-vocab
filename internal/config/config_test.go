package config

import (
	"testing"
	"time"
)

func TestLoadConfigLoadsDatabaseAndJWTSettings(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_USER", "test-user")
	t.Setenv("DB_PASSWORD", "test-password")
	t.Setenv("DB_NAME", "test-db")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("DB_TIMEZONE", "Asia/Shanghai")
	t.Setenv("JWT_SECRET", "01234567890123456789012345678901")
	t.Setenv("JWT_ISSUER", "ielts-vocab-test")
	t.Setenv("JWT_TTL", "24h")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if config.DBHost != "localhost" || config.DBName != "test-db" {
		t.Fatalf("database configuration was not loaded: %#v", config)
	}
	if config.JWTIssuer != "ielts-vocab-test" || config.JWTTTL != 24*time.Hour {
		t.Fatalf("JWT configuration was not loaded: %#v", config)
	}
}
