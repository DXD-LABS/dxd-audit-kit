package ingest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dxdlabs/dxd-audit-kit/internal/audit"
	"github.com/google/uuid"
)

type DefaultIngestService struct {
	repo audit.Repository
}

func NewIngestService(repo audit.Repository) *DefaultIngestService {
	return &DefaultIngestService{repo: repo}
}

func (s *DefaultIngestService) HandleSigningEvent(ctx context.Context, p SigningEventPayload) (Result, error) {
	// 1. Idempotency check
	existing, err := s.repo.GetIngestEvent(ctx, p.Source, p.EventID)
	if err == nil {
		// Tìm được record cũ -> deduplicated
		return Result{
			SignEventID:  existing.SignEventID,
			Deduplicated: true,
		}, nil
	}

	// 2. Map Document
	var docID uuid.UUID
	doc, err := s.findDocument(ctx, p.Target)
	if err != nil {
		// Create new document
		newDoc := audit.Document{
			Hash:     p.Target.Hash,
			HashAlgo: p.Target.HashAlgo,
			Size:     0, // Unknown from payload
		}
		if p.Target.ExternalID != "" {
			newDoc.ExternalID = &p.Target.ExternalID
		}
		if p.Target.Title != "" {
			newDoc.Title = &p.Target.Title
		}

		createdDoc, err := s.repo.CreateDocument(ctx, newDoc)
		if err != nil {
			return Result{}, fmt.Errorf("failed to create document: %w", err)
		}
		docID = createdDoc.ID
	} else {
		docID = doc.ID
	}

	// 3. Create SignEvent
	locationJSON, _ := json.Marshal(p.Context.Location)
	extraJSON, _ := json.Marshal(map[string]any{
		"event_name":      p.EventName,
		"channel":         p.Context.Channel,
		"auth_method":     p.Context.AuthMethod,
		"onchain_tx_hash": p.Context.OnchainTx,
		"trace_id":        p.Context.TraceID,
		"request":         p.Context.Request,
	})

	signEv := audit.SignEvent{
		DocumentID:  docID,
		SignerID:    &p.Actor.ID,
		SignerEmail: p.Actor.Email,
		IPAddress:   p.Context.IPAddress,
		UserAgent:   p.Context.UserAgent,
		Location:    locationJSON,
		DeviceID:    &p.Context.DeviceID,
		Provider:    &p.Source,
		Extra:       extraJSON,
		SignedAt:    p.EventTime,
	}

	loggedEv, err := s.repo.LogSignEvent(ctx, signEv)
	if err != nil {
		return Result{}, fmt.Errorf("failed to log sign event: %w", err)
	}

	// 4. Save IngestEvent for idempotency
	err = s.repo.CreateIngestEvent(ctx, audit.IngestEvent{
		Source:        p.Source,
		SourceEventID: p.EventID,
		SignEventID:   loggedEv.ID,
	})
	if err != nil {
		// Log error but don't fail the whole request?
		// Actually, if we fail here, the next retry might succeed in creating another sign_event if not careful.
		// For robustness, this should be in a transaction.
	}

	return Result{
		DocumentID:   docID,
		SignEventID:  loggedEv.ID,
		Deduplicated: false,
	}, nil
}

func (s *DefaultIngestService) findDocument(ctx context.Context, t Target) (audit.Document, error) {
	// Try external_id first
	if t.ExternalID != "" {
		doc, err := s.repo.GetDocumentByExternalID(ctx, t.ExternalID)
		if err == nil {
			return doc, nil
		}
	}
	// Then hash
	if t.Hash != "" {
		doc, err := s.repo.GetDocumentByHash(ctx, t.Hash)
		if err == nil {
			return doc, nil
		}
	}
	return audit.Document{}, fmt.Errorf("document not found")
}
