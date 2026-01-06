-- 002_extend_sign_events.sql

BEGIN;

-- Thêm các cột mới cho bảng sign_events
-- Lưu ý: Một số cột có thể đã tồn tại từ 001_init.sql, 
-- nhưng để đảm bảo tính nhất quán với yêu cầu, ta dùng các câu lệnh an toàn hoặc điều chỉnh.
-- PostgreSQL không hỗ trợ ADD COLUMN IF NOT EXISTS cho đến bản 9.6+, 
-- nhưng thông thường ta nên kiểm tra trước hoặc chỉ add những cái thực sự thiếu.

ALTER TABLE sign_events
  ADD COLUMN IF NOT EXISTS signer_id     TEXT,
  ADD COLUMN IF NOT EXISTS location      JSONB,
  ADD COLUMN IF NOT EXISTS device_id     TEXT,
  ADD COLUMN IF NOT EXISTS provider      TEXT,
  ADD COLUMN IF NOT EXISTS extra         JSONB;

-- Index phục vụ query audit/report
CREATE INDEX IF NOT EXISTS idx_sign_events_document_signed_at
  ON sign_events (document_id, signed_at);

CREATE INDEX IF NOT EXISTS idx_sign_events_signer_email_signed_at
  ON sign_events (signer_email, signed_at);

COMMIT;
