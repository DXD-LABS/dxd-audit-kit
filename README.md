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

