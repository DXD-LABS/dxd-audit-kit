# DXD Audit Kit

Bộ công cụ Audit Log chuyên dụng cho DXD Labs.

## Cấu trúc dự án

- `cmd/`: Entry point cho CLI và Server.
- `internal/`: Logic cốt lõi (verify, audit, ingest, analyze, config, logger).
- `pkg/dxdaudit/`: Thư viện API công khai.
- `migrations/`: Cấu trúc cơ sở dữ liệu Postgres.
- `api/`: OpenAPI Specification.
- `docker/`: Cấu hình Docker & Docker Compose.

## Khởi chạy nhanh

```bash
cd docker
docker-compose up -d
```
