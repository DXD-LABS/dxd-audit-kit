package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
}

func Load() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default cho phát triển local với Docker Compose (ánh xạ port 5432)
		dbURL = "postgres://dxd_audit:dxd_audit_password@localhost:5432/dxd_audit?sslmode=disable"
	}
	return Config{
		DatabaseURL: dbURL,
	}
}
