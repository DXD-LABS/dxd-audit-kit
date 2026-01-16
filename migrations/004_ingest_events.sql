-- 004_ingest_events.sql

BEGIN;

-- Thêm external_id và title vào bảng documents
ALTER TABLE documents 
  ADD COLUMN IF NOT EXISTS external_id TEXT,
  ADD COLUMN IF NOT EXISTS title       TEXT;

-- Tạo index cho external_id để tìm kiếm nhanh
CREATE INDEX IF NOT EXISTS idx_documents_external_id ON documents(external_id);

-- Tạo bảng ingest_events phục vụ Idempotency
CREATE TABLE IF NOT EXISTS ingest_events (
    source           TEXT NOT NULL,
    source_event_id  TEXT NOT NULL,
    sign_event_id    UUID NOT NULL,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (source, source_event_id)
);

COMMIT;
