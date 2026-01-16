package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dxdlabs/dxd-audit-kit/internal/audit"
	"github.com/dxdlabs/dxd-audit-kit/internal/config"
	"github.com/dxdlabs/dxd-audit-kit/internal/db"
	"github.com/dxdlabs/dxd-audit-kit/internal/ingest"
	"github.com/dxdlabs/dxd-audit-kit/internal/logger"
	"github.com/dxdlabs/dxd-audit-kit/migrations"
)

func main() {
	fmt.Println("DXD Audit Server starting...")

	cfg := config.Load()

	// Khởi tạo Database
	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Tự động chạy migrations
	if err := migrations.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Khởi tạo Repository & Service
	auditRepo := audit.NewRepository(database)
	ingestSvc := ingest.NewIngestService(auditRepo)
	ingestHandler := ingest.NewHTTPHandler(cfg, ingestSvc)

	mux := http.NewServeMux()
	ingestHandler.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	logger.Info("Server listening", "port", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
