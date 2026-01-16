package ingest

import (
	"context"

	"github.com/google/uuid"
)

type DefaultIngestService struct {
}

func NewIngestService() *DefaultIngestService {
	return &DefaultIngestService{}
}

func (s *DefaultIngestService) HandleSigningEvent(ctx context.Context, p SigningEventPayload) (Result, error) {
	// For now, return a dummy result as per the skeleton requirement
	return Result{
		DocumentID:   uuid.New(),
		SignEventID:  uuid.New(),
		Deduplicated: false,
	}, nil
}
