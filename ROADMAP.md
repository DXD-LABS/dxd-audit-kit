# Roadmap – dxd-audit-kit

dxd-audit-kit là Go + Postgres toolkit cho **audit logs**, **verification** và **anomaly detection** quanh chữ ký số và e-sign workflows.

Mục tiêu tổng: trở thành **hạ tầng chứng cứ (evidence backend)** cho dsign.foundation và các nền tảng e-sign/fintech khác.

---

## Phase 1 – Core MVP (ĐÃ HOÀN THÀNH)

**Trạng thái:** Done ✅

**Mục tiêu:** Có pipeline end-to-end: verify → lưu audit → report.

### Deliverables

- Core modules:
    - `verify`: tính hash tài liệu, tạo `VerifyResult`.
    - `audit`: schema `documents`, `sign_events` trên Postgres, repository CRUD.
    - `report`: sinh `DocumentReport` / `SignerReport` ra JSON/CSV.

- CLI:
    - `verify --file`: tính hash, tạo/lấy document.
    - `log-event`: tạo sign_event cơ bản cho document.
    - `report document` / `report signer`: xuất report JSON/CSV.

- Infra:
    - `docker-compose.yml` cho Postgres.
    - Migration `001` / `002` cho schema cơ bản.
    - GitHub Actions: `go test ./...`.

---

## Phase 2 – Rich audit trail & anomaly (HIỆN TẠI)

**Trạng thái:** In progress 🚧

**Mục tiêu:** Biến audit trail thành **công cụ compliance & phân tích rủi ro**, không chỉ là log.

### 2.1. Mở rộng schema & reporting

**Done:**

- Thêm field cho `sign_events`:
    - `signer_id`, `signer_email`, `ip_address`, `user_agent`,
    - `location` (JSONB), `device_id`, `provider`, `extra` (JSONB).
- Report:
    - `DocumentReport` & `SignerReport` xuất JSON/CSV.
    - CLI:
        - `report document --document-id <id> --format json|csv`
        - `report signer --email <mail> --from ... --to ... --format json|csv`

**Tiếp tục tinh chỉnh:**

- Bổ sung các field cần cho compliance thực tế (vd: reason revoke/decline, channel, auth method).

### 2.2. Rule-based anomaly detection

**Done:**

- Migration `003_anomaly_scores` với bảng:
    - `anomaly_scores(id, sign_event_id, score, labels JSONB, created_at)`.
- Package `analyze`:
    - Rule `night_time`: ký ngoài khung giờ làm việc.
    - Rule `multi_ip`: một signer dùng nhiều IP trong khoảng thời gian ngắn.
    - Hàm `AnalyzeDocument(documentID)` lưu kết quả vào `anomaly_scores`.
- CLI:
    - `analyze document --document-id <id>`.

- Report:
    - Nhúng `anomalies` + `anomaly_summary` vào JSON output.

### 2.3. Sample scenarios

**Done:**

- 5 scenario trong `examples/`:
    - Night-time signing.
    - Multi-IP signer.
    - Normal office signing (no anomaly).
    - Multi-document signer timeline.
    - SIEM-friendly export demo.

---

## Phase 3 – Ingest layer cho dsign.foundation

**Trạng thái:** Planned 🗺

**Mục tiêu:** dxd-audit-kit trở thành **audit backend mặc định** của dsign.foundation.

### 3.1. HTTP ingest API

- Tạo service `dxd-audit-server`:
    - `POST /v1/events`
        - Auth: Bearer token hoặc HMAC.
        - Body: `SigningEvent` chuẩn.
        - Response:
          ```json
          {
            "status": "ok",
            "document_id": "UUID",
            "sign_event_id": "UUID",
            "deduplicated": true
          }
          ```

- JSON `SigningEvent` (draft):
    - `event_id`: ID duy nhất từ dsign (dùng idempotency).
    - `event_type`: `document.created | document.uploaded | document.sent | document.viewed | document.signed | document.revoked | signer.auth.* | signer.declined ...`
    - `occurred_at`: ISO8601.
    - `document`: `{ external_id, hash, hash_algo, title }`
    - `signer`: `{ id, email }`
    - `context`: `{ ip_address, user_agent, location{country,city}, device_id, provider, onchain_tx_hash, trace_id, extra{...} }`

### 3.2. Idempotency & mapping

- Bảng `ingest_events`:
    - `source` (`"dsign"`), `source_event_id`, `sign_event_id`, `created_at`, `UNIQUE (source, source_event_id)`.
- Flow:
    - Check `(source, source_event_id)` → nếu tồn tại, trả `deduplicated: true`.
    - Upsert document (theo `hash + external_id`).
    - Insert `sign_events` với đầy đủ context.

### 3.3. Staging integration với dsign.foundation

- Viết adapter client bên dsign:
    - Gửi event thực (create/upload/sign/decline/revoke).
- E2E:
    - Event từ dsign → dxd-audit-kit → `report document` + `analyze document` hoạt động.
- Định nghĩa SLO cho ingest:
    - Ví dụ: 99% event được ghi vào DB trong ≤ 10 giây.

---

## Phase 4 – Production readiness & SIEM integration

**Trạng thái:** Planned

**Mục tiêu:** Đủ tin cậy để chạy trong môi trường production của dsign và khách hàng doanh nghiệp.

### 4.1. Production readiness

- Healthcheck:
    - `GET /healthz` (DB, migrations, queue).
- Metrics (Prometheus):
    - Request count/latency, ingest success/fail, analyze duration.
- Logging:
    - JSON log (level, trace_id, document_id, event_type).
- Config:
    - Env-based config (`DATABASE_URL`, `INGEST_API_TOKEN`, log level, etc.).

- Tài liệu:
    - `PRODUCTION_READINESS.md`:
        - Deploy guideline (Docker/K8s).
        - Backup/restore Postgres.
        - Migration strategy.

### 4.2. SIEM / log platform integration

- Export:
    - NDJSON line-based output cho SIEM.
    - Option: push trực tiếp tới Kafka/queue để hệ thống log khác tiêu thụ.
- Docs:
    - Ví dụ ingest vào ELK/Splunk/Datadog.

---

## Phase 5 – Multi-tenant & external adopters

**Trạng thái:** Future

**Mục tiêu:** Cho phép nhiều ứng dụng (không chỉ dsign) dùng chung dxd-audit-kit.

- Multi-tenant schema:
    - Thêm `tenant_id`/`source` rõ ràng cho mọi bảng (`documents`, `sign_events`, `anomaly_scores`, `ingest_events`).
- Tenant-level API keys:
    - Mỗi client/app có token riêng, phân quyền và rate-limit.
- Template integration:
    - Example client cho: e-sign platform khác, fintech app, web3 dApp.

---

## Phase 6 – AI & advanced fraud analytics

**Trạng thái:** Future

**Mục tiêu:** Nâng cấp từ rule-based lên **AI-assisted risk engine**.

- Data preparation:
    - Tạo feature store từ audit log: tần suất ký, số IP/country, lịch sử tranh chấp, giá trị hợp đồng.
- Model:
    - Unsupervised anomaly (clustering, outlier detection).
    - Risk scoring per document/signer/transaction.
- LLM layer:
    - Tóm tắt risk story cho 1 document/signer (“tại sao bị chấm điểm cao?”).
- UX:
    - API trả `risk_score` + giải thích ngắn.

---

## Phase song song – Open source & cộng đồng

- Chính thức hóa:
    - `LICENSE`, `CODE_OF_CONDUCT.md`, `CONTRIBUTING.md`, `SECURITY.md`.
- Issues & Projects:
    - `good first issue`, Project board cho từng phase.
- Content:
    - Blog/talk: “Building an e-sign audit trail engine with Go + Postgres”, “Rule-based anomaly detection for digital signatures”.

---
