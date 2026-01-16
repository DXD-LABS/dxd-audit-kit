package config

import (
	"os"
)

type Config struct {
	DatabaseURL    string
	IngestAPIToken string
}

func Load() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default cho phát triển local với Docker Compose (ánh xạ port 5432)
		dbURL = "postgres://dxd_audit:dxd_audit_password@localhost:5432/dxd_audit?sslmode=disable"
	}

	apiToken := os.Getenv("INGEST_API_TOKEN")
	if apiToken == "" {
		apiToken = "default-token"
	}

	return Config{
		DatabaseURL:    dbURL,
		IngestAPIToken: apiToken,
	}
}
