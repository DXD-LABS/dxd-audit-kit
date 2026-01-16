package ingest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dxdlabs/dxd-audit-kit/internal/config"
)

type mockIngestService struct{}

func (m *mockIngestService) HandleSigningEvent(ctx any, p SigningEventPayload) (Result, error) {
	return Result{}, nil
}

func TestHandlePostEvent(t *testing.T) {
	cfg := config.Config{IngestAPIToken: "test-token"}
	svc := NewIngestService()
	handler := NewHTTPHandler(cfg, svc)

	payload := SigningEventPayload{
		EventID:   "dsign-evt-12345",
		EventName: "document.signed",
		EventTime: time.Now(),
		Actor: Actor{
			Email: "user@example.com",
		},
		Target: Target{
			Hash: "abc123...",
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/v1/events", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer test-token")

	rr := httptest.NewRecorder()
	handler.handlePostEvent(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]any
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp["status"] != "ok" {
		t.Errorf("handler returned unexpected status: got %v want %v", resp["status"], "ok")
	}
}

func TestHandlePostEvent_Unauthorized(t *testing.T) {
	cfg := config.Config{IngestAPIToken: "test-token"}
	svc := NewIngestService()
	handler := NewHTTPHandler(cfg, svc)

	req := httptest.NewRequest("POST", "/v1/events", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")

	rr := httptest.NewRecorder()
	handler.handlePostEvent(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestHandlePostEvent_MissingFields(t *testing.T) {
	cfg := config.Config{IngestAPIToken: "test-token"}
	svc := NewIngestService()
	handler := NewHTTPHandler(cfg, svc)

	payload := SigningEventPayload{
		EventID: "dsign-evt-12345",
		// Missing other fields
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/v1/events", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer test-token")

	rr := httptest.NewRecorder()
	handler.handlePostEvent(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
