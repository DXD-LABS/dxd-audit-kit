# dxd-audit-kit

`dxd-audit-kit` is a Go + Postgres toolkit for **audit logs**, **verification**, and **anomaly detection** around digital signatures and e-sign workflows.

- 🔐 Verify signed documents and signatures
- 🧾 Generate structured, queryable audit trails
- 🕵️ Detect suspicious signing activity with rules and AI

## Features

- Document & signature verification (hash, certificates, timestamps)
- Normalized audit log schema on Postgres
- CLI and Go library for easy integration into existing e-sign platforms
- Pluggable rules engine for anomaly detection (IP, geo, time, device, etc.)

## Tech stack

- Language: Go
- Database: PostgreSQL
- Interfaces: CLI (`dxd-audit-cli`) and Go package (`github.com/dxdlabs/dxd-audit-kit/pkg/dxdaudit`)

## Production Readiness

Dự án được thiết kế để có thể triển khai trên môi trường Production với các tiêu chuẩn:
- **Logging:** Structured JSON logging qua `stdout`.
- **Config:** Quản lý hoàn toàn qua biến môi trường.
- **Database:** Hỗ trợ cơ chế migration và backup.
- **Deployment:** Hỗ trợ Docker và Kubernetes.

Xem chi tiết tại [PRODUCTION_READINESS.md](./PRODUCTION_READINESS.md).

## Getting started

### 1. Khởi động Database
Dự án sử dụng Postgres. Bạn có thể khởi động nhanh bằng Docker Compose:
```bash
docker-compose up -d
```

### 2. Cấu hình Database
Mặc định CLI sẽ kết nối tới `localhost:5432`. Nếu chạy trong môi trường Docker hoặc cần kết nối tới host khác, hãy đặt biến môi trường `DATABASE_URL`:
```bash
export DATABASE_URL="postgres://dxd_audit:dxd_audit_password@postgres:5432/dxd_audit?sslmode=disable"
```

### 3. Chạy CLI
**Verify tài liệu:**
```bash
go run ./cmd/dxd-audit-cli verify --file path/to/document.pdf
```

**Ghi log sự kiện ký:**
```bash
go run ./cmd/dxd-audit-cli log-event --file path/to/document.pdf --signer-email user@example.com --ip 1.2.3.4
```

### 4. Reporting
Hỗ trợ xuất báo cáo dưới định dạng JSON (mặc định) hoặc CSV.

**Báo cáo theo tài liệu (Document Report):**
```bash
# Xuất JSON
go run ./cmd/dxd-audit-cli report document --document-id <UUID>

# Xuất CSV
go run ./cmd/dxd-audit-cli report document --document-id <UUID> --format csv
```

*Ví dụ JSON Output:*
```json
{
  "document": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "hash_algo": "sha256",
    "size": 1024,
    "created_at": "2023-10-27T10:00:00Z"
  },
  "events": [],
  "sign_count": 2,
  "first_signed_at": "2023-10-27T10:05:00Z",
  "last_signed_at": "2023-10-27T11:00:00Z",
  "unique_ips": ["1.2.3.4", "5.6.7.8"]
}
```

**Báo cáo theo người ký (Signer Report):**
```bash
# Xuất JSON với khoảng thời gian
go run ./cmd/dxd-audit-cli report signer --email user@example.com --from 2023-01-01 --to 2023-12-31

# Xuất CSV
go run ./cmd/dxd-audit-cli report signer --email user@example.com --format csv
```

*Ví dụ CSV Output:*
```csv
document_hash,signer_email,signed_at,ip_address,provider
e3b0c442...,user@example.com,2023-10-27T10:05:00Z,1.2.3.4,adobe
e3b0c442...,user@example.com,2023-10-27T11:00:00Z,5.6.7.8,docusign
```

### 5. Anomaly Detection
Hệ thống hỗ trợ phân tích các sự kiện ký để phát hiện các dấu hiệu bất thường dựa trên các quy tắc (rules):
- **Night-time signing:** Ký vào khung giờ lạ (22:00 - 06:00).
- **Multi-IP signing:** Một người ký sử dụng nhiều địa chỉ IP khác nhau trong một khoảng thời gian ngắn (1 giờ).

**Phân tích tài liệu:**
```bash
go run ./cmd/dxd-audit-cli analyze document --document-id <UUID>
```

**Xem kết quả trong Report:**
Khi chạy lệnh `report`, thông tin về anomaly sẽ được nhúng trực tiếp vào output JSON.
```bash
go run ./cmd/dxd-audit-cli report document --document-id <UUID> --format json
```

*Ví dụ JSON Output với Anomaly:*
```json
{
  "document": {
    "id": "...",
    "hash": "..."
  },
  "events": [
    { "id": "..." }
  ],
  "anomalies": [
    {
      "sign_event_id": "...",
      "score": 0.3,
      "labels": { "night_time": true },
      "created_at": "..."
    }
  ],
  "anomaly_summary": {
    "anomaly_count": 1,
    "max_score": 0.3,
    "avg_score": 0.3,
    "common_labels": {
      "night_time": 1
    }
  }
}
```

## Security Considerations

Vấn đề bảo mật là ưu tiên hàng đầu trong các hệ thống Audit và Chữ ký số. Dưới đây là các nguyên tắc bảo mật được áp dụng và khuyến nghị:

### 1. Data Logging (Dữ liệu nên và không nên log)
- **Nên log:** Metadata của sự kiện (Timestamp, Document Hash, Signer ID/Email, IP Address, User Agent, Provider).
- **Không nên log:** 
    - Nội dung nhạy cảm bên trong tài liệu (trừ khi có yêu cầu nghiệp vụ đặc thù và đã được mã hóa).
    - Thông tin định danh cá nhân (PII) không cần thiết.
    - Token hoặc Secret Key của các hệ thống tích hợp.
- **Log Level:** Sử dụng structured logging (JSON) để dễ dàng tích hợp với các hệ thống SIEM/SOC nhằm phát hiện sớm các hành vi tấn công.

### 2. Database Protection
- **Kết nối:** Luôn sử dụng kết nối bảo mật (TLS/SSL) giữa ứng dụng và Postgres.
- **Phân vùng mạng:** Database nên được đặt trong mạng nội bộ (private subnet), không mở public port (5432) ra ngoài internet.
- **Mã hóa:** Khuyến nghị bật cơ chế mã hóa dữ liệu khi lưu trữ (Encryption at rest) ở tầng storage.

### 3. Application RBAC (Phân quyền Role)
Hệ thống khuyến nghị phân chia 3 nhóm quyền chính:
- **Signer (Người ký):** Chỉ có quyền gửi yêu cầu ký và ghi nhận sự kiện ký thông qua hệ thống tích hợp. Không có quyền truy cập trực tiếp vào Audit Log.
- **Auditor (Người kiểm toán):** Có quyền đọc (Read-only) các báo cáo Audit Trail, xem kết quả Anomaly Detection để kiểm tra tính toàn vẹn của giao dịch.
- **Admin (Quản trị viên):** Có quyền cấu hình hệ thống, quản lý quy tắc (Rules) phân tích bất thường và quản lý các tích hợp đầu vào.

---

## API / CLI Reference

### CLI Reference

Tóm tắt các lệnh chính của `dxd-audit-cli`:

| Lệnh | Input chính | Output | Mô tả |
| :--- | :--- | :--- | :--- |
| `verify` | `--file` | Document ID, Hash | Xác thực file và đăng ký tài liệu vào DB. |
| `log-event` | `--document-id` hoặc `--file`, `--signer-email` | Event ID, Signed At | Ghi lại một sự kiện ký tài liệu. |
| `report document` | `--document-id`, `--format` | JSON/CSV/NDJSON | Xuất báo cáo lịch sử của một tài liệu. |
| `report signer` | `--email`, `--from`, `--to`, `--format` | JSON/CSV/NDJSON | Xuất báo cáo các hoạt động của một người ký. |
| `analyze document` | `--document-id` | Danh sách Anomaly | Phân tích các dấu hiệu bất thường cho tài liệu. |

Sử dụng `--help` sau mỗi lệnh để xem chi tiết tất cả các flag.

### API Reference (OpenAPI)

Dự án cung cấp một HTTP server (`dxd-audit-server`) để tích hợp qua API. Chi tiết đặc tả API (OpenAPI 3.0) có thể tìm thấy tại:

- 📄 [api/openapi.yaml](./api/openapi.yaml)

HTTP Server mặc định lắng nghe tại cổng `8080`.

## Community & Contributing

Chúng tôi hoan nghênh mọi đóng góp từ cộng đồng!

- 🤝 **[CONTRIBUTING.md](./CONTRIBUTING.md):** Hướng dẫn đóng góp và các ý tưởng cho người mới.
- 📜 **[CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md):** Quy tắc ứng xử trong cộng đồng.
- 🛡️ **[SECURITY.md](./SECURITY.md):** Chính sách bảo mật và báo cáo lỗ hổng.
- 🎫 **[Open an Issue](https://github.com/dxdlabs/dxd-audit-kit/issues/new/choose):** Báo lỗi hoặc yêu cầu tính năng mới.

