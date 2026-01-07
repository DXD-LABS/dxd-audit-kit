-- 003_anomaly_scores.sql
BEGIN;

CREATE TABLE IF NOT EXISTS anomaly_scores (
    id UUID PRIMARY KEY,
    sign_event_id UUID NOT NULL REFERENCES sign_events(id),
    score REAL NOT NULL,
    labels JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_anomaly_scores_event_id ON anomaly_scores(sign_event_id);
CREATE INDEX IF NOT EXISTS idx_anomaly_scores_created_at ON anomaly_scores(created_at);

COMMIT;
