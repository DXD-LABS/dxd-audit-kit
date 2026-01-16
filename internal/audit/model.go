package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Document đại diện cho 1 tài liệu được verify
type Document struct {
	ID         uuid.UUID `json:"id"`
	Hash       string    `json:"hash"`
	HashAlgo   string    `json:"hash_algo"`
	ExternalID *string   `json:"external_id,omitempty"`
	Title      *string   `json:"title,omitempty"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"created_at"`
}

type IngestEvent struct {
	Source        string    `json:"source"`
	SourceEventID string    `json:"source_event_id"`
	SignEventID   uuid.UUID `json:"sign_event_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// SignEvent đại diện cho một sự kiện ký
type SignEvent struct {
	ID          uuid.UUID       `json:"id"`
	DocumentID  uuid.UUID       `json:"document_id"`
	SignerID    *string         `json:"signer_id,omitempty"`
	SignerEmail string          `json:"signer_email"`
	IPAddress   string          `json:"ip_address"`
	UserAgent   string          `json:"user_agent"`
	Location    json.RawMessage `json:"location,omitempty"`
	DeviceID    *string         `json:"device_id,omitempty"`
	Provider    *string         `json:"provider,omitempty"`
	Extra       json.RawMessage `json:"extra,omitempty"`
	SignedAt    time.Time       `json:"signed_at"`
}

// AnomalyScore đại diện cho điểm số bất thường của một sự kiện ký
type AnomalyScore struct {
	ID          uuid.UUID       `json:"id"`
	SignEventID uuid.UUID       `json:"sign_event_id"`
	Score       float32         `json:"score"`
	Labels      json.RawMessage `json:"labels"`
	CreatedAt   time.Time       `json:"created_at"`
}
