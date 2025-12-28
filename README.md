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

