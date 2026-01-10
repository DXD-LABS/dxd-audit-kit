# Production Readiness Guide

Tài liệu này hướng dẫn cách vận hành `dxd-audit-kit` trong môi trường Production.

## 1. Logging (JSON)

Hệ thống sử dụng Structured Logging với định dạng JSON (thông qua package `log/slog` của Go). Điều này giúp dễ dàng tích hợp với các hệ thống quản lý log như ELK Stack, Datadog, hoặc CloudWatch.

- **Định dạng:** JSON
- **Output:** `stdout` (theo nguyên tắc Twelve-Factor App)
- **Level:** INFO, WARN, ERROR

Ví dụ log:
```json
{"time":"2023-10-27T10:00:00Z","level":"INFO","msg":"Starting server","port":8080}
```

## 2. Cấu hình (Environment Variables)

Ứng dụng được cấu hình hoàn toàn thông qua biến môi trường.

| Biến | Mô tả | Mặc định |
| --- | --- | --- |
| `DATABASE_URL` | Chuỗi kết nối PostgreSQL | `postgres://dxd_audit:dxd_audit_password@localhost:5432/dxd_audit?sslmode=disable` |
| `LOG_LEVEL` | Mức độ log (debug, info, warn, error) | `info` |

## 3. Chiến lược Migration (Database)

Chúng tôi sử dụng thư mục `migrations/` để quản lý thay đổi schema.

- **Môi trường Dev:** Chạy tự động khi khởi động hoặc qua Docker Compose.
- **Môi trường Prod:** 
    - Khuyến khích sử dụng các công cụ như `golang-migrate` hoặc chạy SQL thủ công trong CI/CD pipeline trước khi deploy ứng dụng mới.
    - Luôn backup database trước khi chạy migration.
    - Các file migration được đánh số thứ tự (ví dụ: `001_init.sql`, `002_extend_sign_events.sql`).

## 4. Backup & Recovery

- **Dữ liệu:** PostgreSQL là nguồn dữ liệu duy nhất.
- **Chiến lược:**
    - Sử dụng `pg_dump` để backup định kỳ.
    - Lưu trữ bản backup ở nơi an toàn (S3, GCS) có thiết lập vòng đời (retention policy).
    - Kiểm tra khả năng khôi phục (restore test) ít nhất một lần mỗi quý.

Ví dụ lệnh backup:
```bash
docker exec dxd-audit-postgres pg_dump -U dxd_audit dxd_audit > backup_$(date +%F).sql
```

## 5. Container & Kubernetes

### Docker
Dự án cung cấp `Dockerfile` để đóng gói ứng dụng.

```bash
docker build -t dxd-audit-kit:latest -f docker/Dockerfile .
```

### Kubernetes (K8s)
Để chạy trong K8s, bạn cần chuẩn bị:
1. **Deployment:** Chạy ứng dụng API.
2. **Service:** Expose API ra ngoài hoặc nội bộ.
3. **Secret/ConfigMap:** Lưu trữ `DATABASE_URL`.
4. **Health Checks:** 
    - Liveness Probe: Kiểm tra process còn sống.
    - Readiness Probe: Kiểm tra kết nối DB thành công.

Gợi ý cấu hình Resources:
- CPU: `100m` (request) / `500m` (limit)
- Memory: `128Mi` (request) / `512Mi` (limit)

## 6. Giám sát (Monitoring)

- Kiểm tra các metric cơ bản của container (CPU, RAM).
- Theo dõi tỷ lệ lỗi (Error Rate) từ log JSON.
- Thiết lập cảnh báo (Alerting) khi Database connection bị lỗi quá lâu.
