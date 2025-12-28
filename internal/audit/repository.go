package audit

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	CreateDocument(ctx context.Context, doc Document) (Document, error)
	GetDocumentByHash(ctx context.Context, hash string) (Document, error)
	LogSignEvent(ctx context.Context, ev SignEvent) (SignEvent, error)
	ListEventsByDocument(ctx context.Context, docID uuid.UUID) ([]SignEvent, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateDocument(ctx context.Context, doc Document) (Document, error) {
	if doc.ID == uuid.Nil {
		doc.ID = uuid.New()
	}
	query := `INSERT INTO documents (id, hash, hash_algo, size, created_at) 
              VALUES ($1, $2, $3, $4, CASE WHEN $5 = '0001-01-01 00:00:00+00'::timestamptz THEN NOW() ELSE $5 END) 
              RETURNING created_at`

	err := r.db.QueryRowContext(ctx, query, doc.ID, doc.Hash, doc.HashAlgo, doc.Size, doc.CreatedAt).Scan(&doc.CreatedAt)
	if err != nil {
		return Document{}, fmt.Errorf("failed to create document: %w", err)
	}
	return doc, nil
}

func (r *postgresRepo) GetDocumentByHash(ctx context.Context, hash string) (Document, error) {
	var doc Document
	query := `SELECT id, hash, hash_algo, size, created_at FROM documents WHERE hash = $1`
	err := r.db.QueryRowContext(ctx, query, hash).Scan(&doc.ID, &doc.Hash, &doc.HashAlgo, &doc.Size, &doc.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Document{}, fmt.Errorf("document not found: %w", err)
		}
		return Document{}, fmt.Errorf("failed to get document: %w", err)
	}
	return doc, nil
}

func (r *postgresRepo) LogSignEvent(ctx context.Context, ev SignEvent) (SignEvent, error) {
	if ev.ID == uuid.Nil {
		ev.ID = uuid.New()
	}
	query := `INSERT INTO sign_events (id, document_id, signer_email, ip_address, user_agent, signed_at) 
              VALUES ($1, $2, $3, $4, $5, CASE WHEN $6 = '0001-01-01 00:00:00+00'::timestamptz THEN NOW() ELSE $6 END) 
              RETURNING signed_at`

	err := r.db.QueryRowContext(ctx, query, ev.ID, ev.DocumentID, ev.SignerEmail, ev.IPAddress, ev.UserAgent, ev.SignedAt).Scan(&ev.SignedAt)
	if err != nil {
		return SignEvent{}, fmt.Errorf("failed to log sign event: %w", err)
	}
	return ev, nil
}

func (r *postgresRepo) ListEventsByDocument(ctx context.Context, docID uuid.UUID) ([]SignEvent, error) {
	query := `SELECT id, document_id, signer_email, ip_address, user_agent, signed_at FROM sign_events WHERE document_id = $1`
	rows, err := r.db.QueryContext(ctx, query, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	var events []SignEvent
	for rows.Next() {
		var ev SignEvent
		if err := rows.Scan(&ev.ID, &ev.DocumentID, &ev.SignerEmail, &ev.IPAddress, &ev.UserAgent, &ev.SignedAt); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, ev)
	}
	return events, nil
}
