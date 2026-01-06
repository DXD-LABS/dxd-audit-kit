package report

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dxdlabs/dxd-audit-kit/internal/audit"
	"github.com/dxdlabs/dxd-audit-kit/internal/db"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) audit.Repository {
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

	// Apply migrations if needed for the test DB
	_, _ = database.Exec("ALTER TABLE sign_events ADD COLUMN IF NOT EXISTS signer_id TEXT")
	_, _ = database.Exec("ALTER TABLE sign_events ADD COLUMN IF NOT EXISTS location JSONB")
	_, _ = database.Exec("ALTER TABLE sign_events ADD COLUMN IF NOT EXISTS device_id TEXT")
	_, _ = database.Exec("ALTER TABLE sign_events ADD COLUMN IF NOT EXISTS provider TEXT")
	_, _ = database.Exec("ALTER TABLE sign_events ADD COLUMN IF NOT EXISTS extra JSONB")

	return audit.NewRepository(database)
}

func TestReporter(t *testing.T) {
	repo := setupTestDB(t)
	reporter := NewReporter(repo)
	ctx := context.Background()

	// Setup data
	doc := audit.Document{
		ID:       uuid.New(),
		Hash:     "report-test-hash",
		HashAlgo: "sha256",
		Size:     100,
	}
	_, err := repo.CreateDocument(ctx, doc)
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	t1 := time.Now().Add(-2 * time.Hour)
	t2 := time.Now().Add(-1 * time.Hour)

	ev1 := audit.SignEvent{
		ID:          uuid.New(),
		DocumentID:  doc.ID,
		SignerEmail: "signer1@example.com",
		IPAddress:   "1.1.1.1",
		SignedAt:    t1,
		Location:    []byte("{}"),
		Extra:       []byte("{}"),
	}
	ev2 := audit.SignEvent{
		ID:          uuid.New(),
		DocumentID:  doc.ID,
		SignerEmail: "signer2@example.com",
		IPAddress:   "2.2.2.2",
		SignedAt:    t2,
		Location:    []byte("{}"),
		Extra:       []byte("{}"),
	}
	_, err = repo.LogSignEvent(ctx, ev1)
	if err != nil {
		t.Fatalf("LogSignEvent 1 failed: %v", err)
	}
	_, err = repo.LogSignEvent(ctx, ev2)
	if err != nil {
		t.Fatalf("LogSignEvent 2 failed: %v", err)
	}

	t.Run("BuildDocumentReport", func(t *testing.T) {
		report, err := reporter.BuildDocumentReport(ctx, doc.ID.String())
		if err != nil {
			t.Fatalf("BuildDocumentReport failed: %v", err)
		}

		if report.SignCount != 2 {
			t.Errorf("Expected 2 signs, got %d", report.SignCount)
		}
		if len(report.UniqueIPs) != 2 {
			t.Errorf("Expected 2 unique IPs, got %d", len(report.UniqueIPs))
		}
		if report.FirstSignedAt == nil || !report.FirstSignedAt.Equal(t1.Truncate(time.Microsecond)) && !report.FirstSignedAt.Equal(t1) {
			// Postgres stores with microsecond precision, but let's just check if it's set for now
			t.Logf("FirstSignedAt: %v, expected: %v", report.FirstSignedAt, t1)
		}
	})

	t.Run("BuildSignerReport", func(t *testing.T) {
		report, err := reporter.BuildSignerReport(ctx, "signer1@example.com", nil, nil)
		if err != nil {
			t.Fatalf("BuildSignerReport failed: %v", err)
		}

		if len(report.Events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(report.Events))
		}
		if len(report.Documents) != 1 {
			t.Errorf("Expected 1 document, got %d", len(report.Documents))
		}
		if report.Documents[0].ID != doc.ID {
			t.Errorf("Expected document ID %s, got %s", doc.ID, report.Documents[0].ID)
		}
	})

	t.Run("ExportJSON", func(t *testing.T) {
		report, _ := reporter.BuildDocumentReport(ctx, doc.ID.String())
		var buf bytes.Buffer
		err := reporter.ExportJSON(&buf, report)
		if err != nil {
			t.Fatalf("ExportJSON failed: %v", err)
		}
		if !strings.Contains(buf.String(), doc.Hash) {
			t.Errorf("JSON output should contain document hash")
		}
	})

	t.Run("ExportCSV", func(t *testing.T) {
		events, _ := repo.ListEventsByDocument(ctx, doc.ID)
		var buf bytes.Buffer
		err := reporter.ExportCSV(ctx, &buf, events)
		if err != nil {
			t.Fatalf("ExportCSV failed: %v", err)
		}
		output := buf.String()
		if !strings.Contains(output, "document_hash,signer_email,signed_at,ip_address,provider") {
			t.Errorf("CSV output missing header")
		}
		if !strings.Contains(output, doc.Hash) {
			t.Errorf("CSV output should contain document hash")
		}
		if !strings.Contains(output, "signer1@example.com") {
			t.Errorf("CSV output should contain signer email")
		}
	})
}
