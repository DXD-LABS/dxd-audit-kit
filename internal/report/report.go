package report

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/dxdlabs/dxd-audit-kit/internal/audit"
	"github.com/google/uuid"
)

type DocumentReport struct {
	Document       audit.Document       `json:"document"`
	Events         []audit.SignEvent    `json:"events"`
	Anomalies      []audit.AnomalyScore `json:"anomalies,omitempty"`
	AnomalySummary *AnomalySummary      `json:"anomaly_summary,omitempty"`
	SignCount      int                  `json:"sign_count"`
	FirstSignedAt  *time.Time           `json:"first_signed_at"`
	LastSignedAt   *time.Time           `json:"last_signed_at"`
	UniqueIPs      []string             `json:"unique_ips"`
}

type AnomalySummary struct {
	AnomalyCount int                `json:"anomaly_count"`
	MaxScore     float32            `json:"max_score"`
	AvgScore     float32            `json:"avg_score"`
	CommonLabels map[string]float64 `json:"common_labels"`
}

type SignerReport struct {
	SignerEmail    string               `json:"signer_email"`
	Documents      []audit.Document     `json:"documents"`
	Events         []audit.SignEvent    `json:"events"`
	Anomalies      []audit.AnomalyScore `json:"anomalies,omitempty"`
	AnomalySummary *AnomalySummary      `json:"anomaly_summary,omitempty"`
	From           *time.Time           `json:"from,omitempty"`
	To             *time.Time           `json:"to,omitempty"`
}

type Reporter struct {
	repo audit.Repository
}

func NewReporter(repo audit.Repository) *Reporter {
	return &Reporter{repo: repo}
}

func (r *Reporter) BuildDocumentReport(ctx context.Context, docID string) (*DocumentReport, error) {
	id, err := uuid.Parse(docID)
	if err != nil {
		return nil, fmt.Errorf("invalid document id: %w", err)
	}

	doc, err := r.repo.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	events, err := r.repo.ListEventsByDocument(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	report := &DocumentReport{
		Document:  doc,
		Events:    events,
		SignCount: len(events),
	}

	ipMap := make(map[string]bool)
	for _, ev := range events {
		if ev.IPAddress != "" {
			ipMap[ev.IPAddress] = true
		}

		if report.FirstSignedAt == nil || ev.SignedAt.Before(*report.FirstSignedAt) {
			t := ev.SignedAt
			report.FirstSignedAt = &t
		}
		if report.LastSignedAt == nil || ev.SignedAt.After(*report.LastSignedAt) {
			t := ev.SignedAt
			report.LastSignedAt = &t
		}
	}

	for ip := range ipMap {
		report.UniqueIPs = append(report.UniqueIPs, ip)
	}

	// Fetch Anomalies
	anomalies, err := r.repo.ListAnomaliesByDocument(ctx, id)
	if err == nil && len(anomalies) > 0 {
		report.Anomalies = anomalies
		report.AnomalySummary = calculateAnomalySummary(anomalies)
	}

	return report, nil
}

func calculateAnomalySummary(anomalies []audit.AnomalyScore) *AnomalySummary {
	if len(anomalies) == 0 {
		return nil
	}

	summary := &AnomalySummary{
		AnomalyCount: len(anomalies),
		CommonLabels: make(map[string]float64),
	}

	var totalScore float32
	for _, a := range anomalies {
		if a.Score > summary.MaxScore {
			summary.MaxScore = a.Score
		}
		totalScore += a.Score

		var labels map[string]interface{}
		if err := json.Unmarshal(a.Labels, &labels); err == nil {
			for k := range labels {
				summary.CommonLabels[k]++
			}
		}
	}
	summary.AvgScore = totalScore / float32(len(anomalies))

	return summary
}

func (r *Reporter) BuildSignerReport(ctx context.Context, signerEmail string, from, to *time.Time) (*SignerReport, error) {
	events, err := r.repo.ListEventsBySigner(ctx, signerEmail, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	report := &SignerReport{
		SignerEmail: signerEmail,
		Events:      events,
		From:        from,
		To:          to,
	}

	docMap := make(map[uuid.UUID]bool)
	for _, ev := range events {
		if !docMap[ev.DocumentID] {
			doc, err := r.repo.GetDocumentByID(ctx, ev.DocumentID)
			if err == nil {
				report.Documents = append(report.Documents, doc)
				docMap[ev.DocumentID] = true
			}
		}
	}

	// Fetch Anomalies
	anomalies, err := r.repo.ListAnomaliesBySigner(ctx, signerEmail)
	if err == nil && len(anomalies) > 0 {
		// Lọc anomalies theo thời gian nếu cần
		var filtered []audit.AnomalyScore
		for _, a := range anomalies {
			if (from == nil || a.CreatedAt.After(*from) || a.CreatedAt.Equal(*from)) &&
				(to == nil || a.CreatedAt.Before(*to) || a.CreatedAt.Equal(*to)) {
				filtered = append(filtered, a)
			}
		}
		if len(filtered) > 0 {
			report.Anomalies = filtered
			report.AnomalySummary = calculateAnomalySummary(filtered)
		}
	}

	return report, nil
}

func (r *Reporter) ExportJSON(w io.Writer, report interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func (r *Reporter) ExportCSV(ctx context.Context, w io.Writer, events []audit.SignEvent) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Header
	header := []string{"document_hash", "signer_email", "signed_at", "ip_address", "provider"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write csv header: %w", err)
	}

	docCache := make(map[uuid.UUID]string)

	for _, ev := range events {
		hash, ok := docCache[ev.DocumentID]
		if !ok {
			doc, err := r.repo.GetDocumentByID(ctx, ev.DocumentID)
			if err != nil {
				hash = "unknown"
			} else {
				hash = doc.Hash
			}
			docCache[ev.DocumentID] = hash
		}

		provider := ""
		if ev.Provider != nil {
			provider = *ev.Provider
		}

		row := []string{
			hash,
			ev.SignerEmail,
			ev.SignedAt.Format(time.RFC3339),
			ev.IPAddress,
			provider,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write csv row: %w", err)
		}
	}

	return nil
}
