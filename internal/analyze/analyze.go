package analyze

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dxdlabs/dxd-audit-kit/internal/audit"
	"github.com/google/uuid"
)

// AnalyzeDocument lấy tất cả sign_events của document, chạy rule, lưu vào anomaly_scores.
func AnalyzeDocument(ctx context.Context, repo audit.Repository, docID uuid.UUID) ([]audit.AnomalyScore, error) {
	events, err := repo.ListEventsByDocument(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	var results []audit.AnomalyScore
	for _, ev := range events {
		score, labels := analyzeEvent(ev, events)
		if score > 0 {
			anomaly := audit.AnomalyScore{
				ID:          uuid.New(),
				SignEventID: ev.ID,
				Score:       score,
				Labels:      labels,
				CreatedAt:   time.Now(),
			}
			err := repo.SaveAnomalyScore(ctx, anomaly)
			if err != nil {
				return nil, fmt.Errorf("failed to save anomaly score: %w", err)
			}
			results = append(results, anomaly)
		}
	}

	return results, nil
}

func analyzeEvent(ev audit.SignEvent, allEvents []audit.SignEvent) (float32, json.RawMessage) {
	var totalScore float32
	labels := make(map[string]interface{})

	// Rule 1: Ký ngoài giờ làm việc (22:00–06:00)
	hour := ev.SignedAt.Hour()
	if hour >= 22 || hour < 6 {
		totalScore += 0.3
		labels["night_time"] = true
	}

	// Rule 2: Cùng signer sử dụng nhiều ip_address trong khoảng thời gian ngắn (ví dụ 1 giờ)
	// Tìm các event của cùng signer trong vòng 1 giờ trước/sau event hiện tại
	differentIPs := make(map[string]bool)
	for _, other := range allEvents {
		if other.SignerEmail == ev.SignerEmail && other.ID != ev.ID {
			diff := other.SignedAt.Sub(ev.SignedAt)
			if diff < 0 {
				diff = -diff
			}
			if diff <= time.Hour {
				if other.IPAddress != ev.IPAddress {
					differentIPs[other.IPAddress] = true
				}
			}
		}
	}
	if len(differentIPs) > 0 {
		totalScore += 0.5
		labels["multi_ip"] = true
		labels["ip_count"] = len(differentIPs) + 1
	}

	// Giới hạn score tối đa là 1.0
	if totalScore > 1.0 {
		totalScore = 1.0
	}

	if totalScore > 0 {
		lbBytes, _ := json.Marshal(labels)
		return totalScore, lbBytes
	}

	return 0, nil
}
