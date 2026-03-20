package ingest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dxdlabs/dxd-audit-kit/internal/audit"
	"github.com/google/uuid"
)

type mockState struct {
	documentsByHash  map[string]audit.Document
	documentsByExtID map[string]audit.Document
	ingestByKey      map[string]audit.IngestEvent
	signEvents       []audit.SignEvent
}

type txMockRepo struct {
	audit.Repository
	state            *mockState
	failCreateIngest bool
}

func newTxMockRepo() *txMockRepo {
	return &txMockRepo{
		state: &mockState{
			documentsByHash:  make(map[string]audit.Document),
			documentsByExtID: make(map[string]audit.Document),
			ingestByKey:      make(map[string]audit.IngestEvent),
			signEvents:       nil,
		},
	}
}

func (m *txMockRepo) cloneState() *mockState {
	next := &mockState{
		documentsByHash:  make(map[string]audit.Document, len(m.state.documentsByHash)),
		documentsByExtID: make(map[string]audit.Document, len(m.state.documentsByExtID)),
		ingestByKey:      make(map[string]audit.IngestEvent, len(m.state.ingestByKey)),
		signEvents:       append([]audit.SignEvent(nil), m.state.signEvents...),
	}
	for k, v := range m.state.documentsByHash {
		next.documentsByHash[k] = v
	}
	for k, v := range m.state.documentsByExtID {
		next.documentsByExtID[k] = v
	}
	for k, v := range m.state.ingestByKey {
		next.ingestByKey[k] = v
	}
	return next
}

func (m *txMockRepo) WithinTransaction(ctx context.Context, fn func(audit.Repository) error) error {
	_ = ctx
	txRepo := &txMockRepo{
		state:            m.cloneState(),
		failCreateIngest: m.failCreateIngest,
	}
	if err := fn(txRepo); err != nil {
		return err
	}
	m.state = txRepo.state
	return nil
}

func ingestKey(source, sourceEventID string) string {
	return source + "::" + sourceEventID
}

func (m *txMockRepo) GetIngestEvent(ctx context.Context, source, sourceEventID string) (audit.IngestEvent, error) {
	_ = ctx
	ev, ok := m.state.ingestByKey[ingestKey(source, sourceEventID)]
	if !ok {
		return audit.IngestEvent{}, fmt.Errorf("ingest event not found")
	}
	return ev, nil
}

func (m *txMockRepo) CreateIngestEvent(ctx context.Context, ev audit.IngestEvent) error {
	_ = ctx
	if m.failCreateIngest {
		return fmt.Errorf("forced ingest event failure")
	}
	key := ingestKey(ev.Source, ev.SourceEventID)
	if _, exists := m.state.ingestByKey[key]; exists {
		return fmt.Errorf("duplicate ingest event")
	}
	m.state.ingestByKey[key] = ev
	return nil
}

func (m *txMockRepo) GetDocumentByHash(ctx context.Context, hash string) (audit.Document, error) {
	_ = ctx
	doc, ok := m.state.documentsByHash[hash]
	if !ok {
		return audit.Document{}, fmt.Errorf("document not found")
	}
	return doc, nil
}

func (m *txMockRepo) GetDocumentByExternalID(ctx context.Context, externalID string) (audit.Document, error) {
	_ = ctx
	doc, ok := m.state.documentsByExtID[externalID]
	if !ok {
		return audit.Document{}, fmt.Errorf("document not found")
	}
	return doc, nil
}

func (m *txMockRepo) CreateDocument(ctx context.Context, doc audit.Document) (audit.Document, error) {
	_ = ctx
	if existing, ok := m.state.documentsByHash[doc.Hash]; ok {
		return existing, nil
	}
	if doc.ID == uuid.Nil {
		doc.ID = uuid.New()
	}
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
	}
	m.state.documentsByHash[doc.Hash] = doc
	if doc.ExternalID != nil && *doc.ExternalID != "" {
		m.state.documentsByExtID[*doc.ExternalID] = doc
	}
	return doc, nil
}

func (m *txMockRepo) LogSignEvent(ctx context.Context, ev audit.SignEvent) (audit.SignEvent, error) {
	_ = ctx
	if ev.ID == uuid.Nil {
		ev.ID = uuid.New()
	}
	if ev.SignedAt.IsZero() {
		ev.SignedAt = time.Now()
	}
	m.state.signEvents = append(m.state.signEvents, ev)
	return ev, nil
}

func TestHandleSigningEvent_RollbackWhenIngestEventFails(t *testing.T) {
	repo := newTxMockRepo()
	repo.failCreateIngest = true
	svc := NewIngestService(repo)

	payload := SigningEventPayload{
		EventID:   "evt-1",
		EventName: "document.signed",
		EventTime: time.Now(),
		Source:    "dsign.foundation",
		Actor: Actor{
			ID:    "u1",
			Email: "user@example.com",
		},
		Target: Target{
			Hash:     "hash-1",
			HashAlgo: "sha256",
		},
		Context: Context{
			IPAddress: "127.0.0.1",
		},
	}

	_, err := svc.HandleSigningEvent(context.Background(), payload)
	if err == nil {
		t.Fatalf("expected error when ingest event creation fails")
	}

	if len(repo.state.signEvents) != 0 {
		t.Fatalf("expected sign event to be rolled back, got %d", len(repo.state.signEvents))
	}
	if len(repo.state.ingestByKey) != 0 {
		t.Fatalf("expected ingest markers to stay empty after rollback")
	}
}

func TestHandleSigningEvent_DeduplicatesAfterFirstSuccess(t *testing.T) {
	repo := newTxMockRepo()
	svc := NewIngestService(repo)

	payload := SigningEventPayload{
		EventID:   "evt-2",
		EventName: "document.signed",
		EventTime: time.Now(),
		Source:    "dsign.foundation",
		Actor: Actor{
			ID:    "u2",
			Email: "user@example.com",
		},
		Target: Target{
			Hash:     "hash-2",
			HashAlgo: "sha256",
		},
		Context: Context{
			IPAddress: "127.0.0.2",
		},
	}

	first, err := svc.HandleSigningEvent(context.Background(), payload)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if first.Deduplicated {
		t.Fatalf("first call should not be deduplicated")
	}

	second, err := svc.HandleSigningEvent(context.Background(), payload)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	if !second.Deduplicated {
		t.Fatalf("second call should be deduplicated")
	}
	if first.SignEventID != second.SignEventID {
		t.Fatalf("expected same sign event id for duplicate call")
	}
	if len(repo.state.signEvents) != 1 {
		t.Fatalf("expected exactly one sign event persisted, got %d", len(repo.state.signEvents))
	}
}

