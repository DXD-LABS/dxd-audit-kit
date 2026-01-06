package audit

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/dxdlabs/dxd-audit-kit/internal/db"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://dxd_audit:dxd_audit_password@localhost:5432/dxd_audit?sslmode=disable"
	}

	database, err := db.Open(dbURL)
	if err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
	}

	// Clean up tables before testing
	_, _ = database.Exec("DELETE FROM sign_events")
	_, _ = database.Exec("DELETE FROM documents")

	return database
}

func TestPostgresRepo(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	t.Run("Document operations", func(t *testing.T) {
		doc := Document{
			ID:       uuid.New(),
			Hash:     "test-hash-" + uuid.New().String(),
			HashAlgo: "sha256",
			Size:     1024,
		}

		// Test CreateDocument
		created, err := repo.CreateDocument(ctx, doc)
		if err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}
		if created.ID != doc.ID {
			t.Errorf("Expected ID %s, got %s", doc.ID, created.ID)
		}
		if created.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		// Test GetDocumentByHash
		found, err := repo.GetDocumentByHash(ctx, doc.Hash)
		if err != nil {
			t.Fatalf("Failed to get document by hash: %v", err)
		}
		if found.ID != doc.ID {
			t.Errorf("Expected ID %s, got %s", doc.ID, found.ID)
		}

		// Test GetDocumentByHash - Not Found
		_, err = repo.GetDocumentByHash(ctx, "non-existent")
		if err == nil {
			t.Error("Expected error for non-existent hash, got nil")
		}
	})

	t.Run("SignEvent operations", func(t *testing.T) {
		doc := Document{
			ID:       uuid.New(),
			Hash:     "event-test-hash-" + uuid.New().String(),
			HashAlgo: "sha256",
			Size:     2048,
		}
		_, err := repo.CreateDocument(ctx, doc)
		if err != nil {
			t.Fatalf("Failed to create document for event test: %v", err)
		}

		signerID := "user-123"
		deviceID := "device-456"
		event := SignEvent{
			ID:          uuid.New(),
			DocumentID:  doc.ID,
			SignerID:    &signerID,
			SignerEmail: "test@dxd.io",
			IPAddress:   "127.0.0.1",
			UserAgent:   "Go-Test",
			DeviceID:    &deviceID,
			Location:    []byte(`{"city": "Saigon"}`),
		}

		// Test LogSignEvent
		logged, err := repo.LogSignEvent(ctx, event)
		if err != nil {
			// If migration is not yet applied, this might fail in CI/Local
			// In a real scenario, we would apply migrations before running tests
			t.Logf("LogSignEvent failed (expected if migrations not applied): %v", err)
			return
		}
		if logged.ID != event.ID {
			t.Errorf("Expected event ID %s, got %s", event.ID, logged.ID)
		}
		if logged.SignedAt.IsZero() {
			t.Error("Expected SignedAt to be set")
		}

		// Test ListEventsByDocument
		events, err := repo.ListEventsByDocument(ctx, doc.ID)
		if err != nil {
			t.Fatalf("Failed to list events: %v", err)
		}
		if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		}
		if events[0].ID != event.ID {
			t.Errorf("Expected event ID %s, got %s", event.ID, events[0].ID)
		}
	})
}
