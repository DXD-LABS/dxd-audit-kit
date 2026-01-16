package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dxdlabs/dxd-audit-kit/internal/config"
	"github.com/dxdlabs/dxd-audit-kit/internal/ingest"
)

func main() {
	fmt.Println("DXD Audit Server starting...")

	cfg := config.Load()

	ingestSvc := ingest.NewIngestService()
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

	fmt.Printf("Listening on %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
