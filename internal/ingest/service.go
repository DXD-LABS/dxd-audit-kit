package ingest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dxdlabs/dxd-audit-kit/internal/audit"
	"github.com/google/uuid"
)

// DefaultIngestService is the default implementation of IngestService.
type DefaultIngestService struct {
	repo audit.Repository
}

// NewIngestService creates a new DefaultIngestService.
func NewIngestService(repo audit.Repository) *DefaultIngestService {
	return &DefaultIngestService{repo: repo}
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func (s *DefaultIngestService) HandleSigningEvent(ctx context.Context, p SigningEventPayload) (Result, error) {
	// Prepare immutable payload-derived fields once.
	locationJSON, _ := json.Marshal(p.Context.Location)
	extraJSON, _ := json.Marshal(map[string]any{
		"event_name":      p.EventName,
		"channel":         p.Context.Channel,
		"auth_method":     p.Context.AuthMethod,
		"onchain_tx_hash": p.Context.OnchainTx,
		"trace_id":        p.Context.TraceID,
		"request":         p.Context.Request,
	})
	actorID := toString(p.Actor.ID)

	var result Result
	err := s.repo.WithinTransaction(ctx, func(txRepo audit.Repository) error {
		// 1. Idempotency check
		existing, err := txRepo.GetIngestEvent(ctx, p.Source, p.EventID)
		if err == nil {
			result = Result{
				SignEventID:  existing.SignEventID,
				Deduplicated: true,
			}
			return nil
		}

		// 2. Find/Create document
		var docID uuid.UUID
		doc, err := s.findDocumentWithRepo(ctx, txRepo, p.Target)
		if err != nil {
			newDoc := audit.Document{
				Hash:     p.Target.Hash,
				HashAlgo: p.Target.HashAlgo,
				Size:     0,
			}
			if p.Target.ExternalID != nil && toString(p.Target.ExternalID) != "" {
				extID := toString(p.Target.ExternalID)
				newDoc.ExternalID = &extID
			}
			if p.Target.Title != "" {
				newDoc.Title = &p.Target.Title
			}

			createdDoc, err := txRepo.CreateDocument(ctx, newDoc)
			if err != nil {
				return fmt.Errorf("failed to create document: %w", err)
			}
			docID = createdDoc.ID
		} else {
			docID = doc.ID
		}

		// 3. Create SignEvent
		signEv := audit.SignEvent{
			DocumentID:  docID,
			SignerID:    &actorID,
			SignerEmail: p.Actor.Email,
			IPAddress:   p.Context.IPAddress,
			UserAgent:   p.Context.UserAgent,
			Location:    locationJSON,
			DeviceID:    &p.Context.DeviceID,
			Provider:    &p.Source,
			Extra:       extraJSON,
			SignedAt:    p.EventTime,
		}

		loggedEv, err := txRepo.LogSignEvent(ctx, signEv)
		if err != nil {
			return fmt.Errorf("failed to log sign event: %w", err)
		}

		// 4. Save IngestEvent idempotency marker (must succeed, otherwise rollback).
		err = txRepo.CreateIngestEvent(ctx, audit.IngestEvent{
			Source:        p.Source,
			SourceEventID: p.EventID,
			SignEventID:   loggedEv.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to create ingest event: %w", err)
		}

		result = Result{
			DocumentID:   docID,
			SignEventID:  loggedEv.ID,
			Deduplicated: false,
		}
		return nil
	})
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

func (s *DefaultIngestService) findDocument(ctx context.Context, t Target) (audit.Document, error) {
	return s.findDocumentWithRepo(ctx, s.repo, t)
}

func (s *DefaultIngestService) findDocumentWithRepo(ctx context.Context, repo audit.Repository, t Target) (audit.Document, error) {
	// Try external_id first
	extID := toString(t.ExternalID)
	if extID != "" {
		doc, err := repo.GetDocumentByExternalID(ctx, extID)
		if err == nil {
			return doc, nil
		}
	}
	// Then hash
	if t.Hash != "" {
		doc, err := repo.GetDocumentByHash(ctx, t.Hash)
		if err == nil {
			return doc, nil
		}
	}
	return audit.Document{}, fmt.Errorf("document not found")
}
