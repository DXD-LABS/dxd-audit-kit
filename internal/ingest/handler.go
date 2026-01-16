package ingest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dxdlabs/dxd-audit-kit/internal/config"
	"github.com/dxdlabs/dxd-audit-kit/internal/logger"
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
	mux.HandleFunc("GET /healthz", h.handleHealthCheck)
	mux.HandleFunc("GET /swagger.yaml", h.handleSwaggerYAML)
	mux.HandleFunc("GET /swagger", h.handleSwaggerUI)
}

func (h *HTTPHandler) handleSwaggerYAML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "api/openapi.yaml")
}

func (h *HTTPHandler) handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" >
    <style>
      html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
      *, *:before, *:after { box-sizing: inherit; }
      body { margin:0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"> </script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: "/swagger.yaml",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      window.ui = ui
    }
    </script>
</body>
</html>
`
	_, _ = w.Write([]byte(html))
}

func (h *HTTPHandler) handlePostEvent(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()

	// Auth: Bearer <token>
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" || token == "" || token != h.cfg.IngestAPIToken {
		logger.Warn("unauthorized access attempt", "ip", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	var payload SigningEventPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}

	// Validate required fields
	if err := h.validatePayload(payload); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	res, err := h.ingestService.HandleSigningEvent(ctx, payload)
	latency := time.Since(start)

	if err != nil {
		logger.Error("failed to handle signing event", err,
			"event_id", payload.EventID,
			"trace_id", payload.Context.TraceID,
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	// Logging requirement
	logger.Info("event ingested",
		"event_id", payload.EventID,
		"event_name", payload.EventName,
		"source", payload.Source,
		"document_id", res.DocumentID.String(),
		"sign_event_id", res.SignEventID.String(),
		"deduplicated", res.Deduplicated,
		"status_code", http.StatusOK,
		"latency_ms", latency.Milliseconds(),
		"trace_id", payload.Context.TraceID,
	)

	resp := map[string]any{
		"status":        "ok",
		"document_id":   res.DocumentID.String(),
		"sign_event_id": res.SignEventID.String(),
		"deduplicated":  res.Deduplicated,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

type validationError string

func (e validationError) Error() string { return string(e) }

func (h *HTTPHandler) validatePayload(p SigningEventPayload) error {
	if p.EventID == "" {
		return validationError("missing event_id")
	}
	if p.EventName == "" {
		return validationError("missing event_name")
	}
	if p.EventTime.IsZero() {
		return validationError("missing event_time")
	}
	if p.Source == "" {
		return validationError("missing source")
	}

	// Document related validation
	isDocEvent := strings.HasPrefix(p.EventName, "document.") || strings.HasPrefix(p.EventName, "signer.")
	if isDocEvent {
		if p.Target.Hash == "" {
			return validationError("missing target.hash for document event")
		}
		if p.Actor.Email == "" {
			return validationError("missing actor.email for document event")
		}
	}

	return nil
}

func (h *HTTPHandler) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
