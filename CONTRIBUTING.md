# Hướng dẫn đóng góp cho dxd-audit-kit

Chào mừng bạn! Chúng tôi rất vui khi bạn quan tâm đến việc đóng góp cho `dxd-audit-kit`.

## Cách bắt đầu

1. **Tìm Issue:** Xem danh sách [Issues](https://github.com/dxdlabs/dxd-audit-kit/issues). Những issue được gắn nhãn `good first issue` là nơi tốt nhất để bắt đầu.
2. **Thảo luận:** Nếu bạn muốn thêm tính năng mới, hãy tạo một issue để thảo luận trước khi bắt tay vào làm.
3. **Fork & Branch:** Fork repo và tạo một branch mới cho tính năng/sửa lỗi của bạn.
4. **Viết Code:** Đảm bảo code tuân thủ phong cách hiện có và có kèm theo test.
5. **Pull Request:** Gửi PR và mô tả rõ những gì bạn đã thay đổi.

## Các ý tưởng "Good First Issue"

Dưới đây là một số mục bạn có thể đóng góp ngay:

- **Thêm Rule mới cho Anomaly Detection:**
  - Ví dụ: Phát hiện việc ký từ các quốc gia (Geo-fencing) khác nhau trong thời gian ngắn.
  - Ví dụ: Phát hiện việc sử dụng VPN/Proxy (dựa trên IP list).
- **Thêm Exporter mới:**
  - Hiện tại chỉ có JSON, CSV, NDJSON. Bạn có thể thêm exporter cho **Slack**, **Discord** hoặc **Elasticsearch**.
- **Cải thiện Tài liệu:**
  - Sửa lỗi chính tả, làm rõ các bước hướng dẫn hoặc thêm ví dụ code.

## Quy trình phát triển

```bash
# Chạy test
go test ./...

# Chạy lint (nếu có)
# golangci-lint run
```

Cảm ơn bạn đã góp sức làm cho dự án tốt hơn!
