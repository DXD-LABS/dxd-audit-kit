package ingest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dxdlabs/dxd-audit-kit/internal/config"
)

type HTTPHandler struct {
	cfg           config.Config
	ingestService IngestService
}

func NewHTTPHandler(cfg config.Config, svc IngestService) *HTTPHandler {
	return &HTTPHandler{
		cfg:           cfg,
		ingestService: svc,
	}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/events", h.handlePostEvent)
}

func (h *HTTPHandler) handlePostEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Auth: Bearer <token>
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" || token != h.cfg.IngestAPIToken {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var payload SigningEventPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if payload.EventID == "" || payload.EventName == "" || payload.EventTime.IsZero() ||
		payload.Target.Hash == "" || payload.Actor.Email == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	res, err := h.ingestService.HandleSigningEvent(ctx, payload)
	if err != nil {
		// In a real app, we'd use a logger and include trace_id
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"status":        "ok",
		"document_id":   res.DocumentID.String(),
		"sign_event_id": res.SignEventID.String(),
		"deduplicated":  res.Deduplicated,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
