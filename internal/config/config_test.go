package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("Default value", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")
		cfg := Load()
		expected := "postgres://dxd_audit:dxd_audit_password@localhost:5432/dxd_audit?sslmode=disable"
		if cfg.DatabaseURL != expected {
			t.Errorf("expected %s, got %s", expected, cfg.DatabaseURL)
		}
	})

	t.Run("From environment variable", func(t *testing.T) {
		customURL := "postgres://user:pass@remote:5432/db"
		os.Setenv("DATABASE_URL", customURL)
		defer os.Unsetenv("DATABASE_URL")

		cfg := Load()
		if cfg.DatabaseURL != customURL {
			t.Errorf("expected %s, got %s", customURL, cfg.DatabaseURL)
		}
	})
}
