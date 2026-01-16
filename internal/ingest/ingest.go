package ingest

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SigningEventPayload struct {
	EventID   string    `json:"event_id"`
	EventName string    `json:"event_name"` // vd: "document.signed"
	EventTime time.Time `json:"event_time"` // ISO8601
	Source    string    `json:"source"`     // vd: "dsign.foundation"

	Actor   Actor   `json:"actor"`
	Target  Target  `json:"target"`
	Context Context `json:"context"`
}

type Actor struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"` // signer | approver | admin...
	Type  string `json:"type"` // user | system | service_account
}

type Target struct {
	Type       string `json:"type"`        // document | transaction | account | policy...
	ID         string `json:"id"`          // internal id phía source
	ExternalID string `json:"external_id"` // nếu cần phân biệt
	Hash       string `json:"hash"`
	HashAlgo   string `json:"hash_algo"`
	Title      string `json:"title"`
	Version    int    `json:"version"`
}

type Context struct {
	IPAddress  string         `json:"ip_address"`
	UserAgent  string         `json:"user_agent"`
	Location   map[string]any `json:"location"` // country, city...
	DeviceID   string         `json:"device_id"`
	Channel    string         `json:"channel"`     // web | mobile | api...
	AuthMethod string         `json:"auth_method"` // otp_sms | ekyc | password...
	OnchainTx  string         `json:"onchain_tx_hash"`
	TraceID    string         `json:"trace_id"`
	Request    map[string]any `json:"request"` // contract_id, business_unit, ...
}

type IngestService interface {
	HandleSigningEvent(ctx context.Context, p SigningEventPayload) (Result, error)
}

type Result struct {
	DocumentID   uuid.UUID
	SignEventID  uuid.UUID
	Deduplicated bool
}
