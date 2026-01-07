package analyze

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dxdlabs/dxd-audit-kit/internal/audit"
	"github.com/google/uuid"
)

type mockRepo struct {
	audit.Repository
	events []audit.SignEvent
	scores []audit.AnomalyScore
}

func (m *mockRepo) ListEventsByDocument(ctx context.Context, docID uuid.UUID) ([]audit.SignEvent, error) {
	return m.events, nil
}

func (m *mockRepo) SaveAnomalyScore(ctx context.Context, s audit.AnomalyScore) error {
	m.scores = append(m.scores, s)
	return nil
}

func TestAnalyzeDocument(t *testing.T) {
	docID := uuid.New()

	// Event 1: Night time signing
	signedAt1, _ := time.Parse(time.RFC3339, "2023-10-27T23:00:00Z")
	ev1 := audit.SignEvent{
		ID:          uuid.New(),
		DocumentID:  docID,
		SignerEmail: "user@example.com",
		IPAddress:   "1.1.1.1",
		SignedAt:    signedAt1,
	}

	// Event 2: Normal time, same signer, DIFFERENT IP, within 1 hour
	signedAt2, _ := time.Parse(time.RFC3339, "2023-10-27T23:30:00Z")
	ev2 := audit.SignEvent{
		ID:          uuid.New(),
		DocumentID:  docID,
		SignerEmail: "user@example.com",
		IPAddress:   "2.2.2.2",
		SignedAt:    signedAt2,
	}

	repo := &mockRepo{
		events: []audit.SignEvent{ev1, ev2},
	}

	results, err := AnalyzeDocument(context.Background(), repo, docID)
	if err != nil {
		t.Fatalf("AnalyzeDocument failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 anomalies, got %d", len(results))
	}

	// Check ev1 anomaly (should have night_time and multi_ip)
	var foundEv1, foundEv2 bool
	for _, r := range results {
		if r.SignEventID == ev1.ID {
			foundEv1 = true
			var labels map[string]interface{}
			json.Unmarshal(r.Labels, &labels)
			if labels["night_time"] != true {
				t.Errorf("Event 1 should have night_time label")
			}
			if labels["multi_ip"] != true {
				t.Errorf("Event 1 should have multi_ip label")
			}
		}
		if r.SignEventID == ev2.ID {
			foundEv2 = true
			var labels map[string]interface{}
			json.Unmarshal(r.Labels, &labels)
			if labels["night_time"] != true {
				t.Errorf("Event 2 should have night_time label (23:30 is night)")
			}
			if labels["multi_ip"] != true {
				t.Errorf("Event 2 should have multi_ip label")
			}
		}
	}

	if !foundEv1 || !foundEv2 {
		t.Errorf("Missing anomaly results for events")
	}
}
