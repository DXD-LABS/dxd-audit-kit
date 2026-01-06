package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Document đại diện cho 1 tài liệu được verify
type Document struct {
	ID        uuid.UUID `json:"id"`
	Hash      string    `json:"hash"`
	HashAlgo  string    `json:"hash_algo"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
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
