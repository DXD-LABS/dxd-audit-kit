# Scenario 4: Multi-document Signer Timeline

Mục đích: Hiển thị báo cáo lịch sử ký của một cá nhân trên nhiều tài liệu khác nhau.

## Các bước thực hiện

1. **Verify nhiều tài liệu**:
   ```bash
   dxd-audit-cli verify --file ../../sample.pdf
   # Giả sử bạn có thêm file khác
   # dxd-audit-cli verify --file ../../another.pdf
   ```

2. **Một người ký nhiều tài liệu trong vài ngày**:
   ```bash
   # Ký tài liệu 1 ngày hôm trước
   dxd-audit-cli log-event --document-id <DOC_1_ID> --signer-email "active-user@example.com" --ip "1.1.1.1" --signed-at "2026-01-07T10:00:00Z"
   
   # Ký tài liệu 1 lần nữa (ví dụ: bổ sung chữ ký)
   dxd-audit-cli log-event --document-id <DOC_1_ID> --signer-email "active-user@example.com" --ip "1.1.1.1" --signed-at "2026-01-08T09:00:00Z"

   # Ký tài liệu 2 ngày hôm nay
   dxd-audit-cli log-event --document-id <DOC_2_ID> --signer-email "active-user@example.com" --ip "1.1.1.1" --signed-at "2026-01-08T15:00:00Z"
   ```

3. **Chạy báo cáo theo email người ký**:
   ```bash
   dxd-audit-cli report signer --email "active-user@example.com" --from "2026-01-07" --to "2026-01-09" --format json
   ```

## Kỳ vọng kết quả
Báo cáo cho người ký sẽ bao gồm:
- Danh sách tất cả các tài liệu đã ký (`documents`).
- Danh sách các sự kiện ký theo thời gian (`events`).
- Các bất thường liên quan đến người ký này (nếu có).
