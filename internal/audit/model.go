package audit

import "time"

type AuditLog struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      string    `json:"data"`
	Hash      string    `json:"hash"`
	Signature string    `json:"signature"`
}
