# Scenario 3: Normal Office Signing (Baseline - Không có bất thường)

Mục đích: Cung cấp một trường hợp "sạch" để người dùng so sánh với các trường hợp bất thường.

## Các bước thực hiện

1. **Verify tài liệu**:
   ```bash
   dxd-audit-cli verify --file ../../sample.pdf
   ```

2. **Log các sự kiện ký trong giờ hành chính**:
   ```bash
   # Ký lúc 09:00 sáng
   dxd-audit-cli log-event --file ../../sample.pdf --signer-email "staff@company.com" --ip "192.168.1.10" --signed-at "2026-01-08T09:00:00Z"
   
   # Ký lúc 14:00 chiều từ cùng dải IP nội bộ
   dxd-audit-cli log-event --file ../../sample.pdf --signer-email "manager@company.com" --ip "192.168.1.15" --signed-at "2026-01-08T14:00:00Z"
   ```

3. **Phân tích tài liệu**:
   ```bash
   dxd-audit-cli analyze document --document-id <DOCUMENT_ID>
   ```

4. **Kiểm tra báo cáo**:
   ```bash
   dxd-audit-cli report document --document-id <DOCUMENT_ID>
   ```

## Kỳ vọng kết quả
Báo cáo sẽ không có bất thường:
- `anomalies`: `null` hoặc danh sách rỗng.
- `anomaly_summary.anomaly_count`: 0.
- `score` của các sự kiện đều bằng 0.
