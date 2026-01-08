# Scenario 2: Multi-IP Signer (Nhiều IP trong thời gian ngắn)

Mục đích: Minh họa rule `multi_ip` phát hiện một người dùng ký từ nhiều địa chỉ IP khác nhau trong vòng 1 giờ.

## Các bước thực hiện

1. **Verify tài liệu**:
   ```bash
   dxd-audit-cli verify --file ../../sample.pdf
   ```

2. **Log 3 sự kiện ký với IP khác nhau trong 1 giờ**:
   ```bash
   # IP 1
   dxd-audit-cli log-event --file ../../sample.pdf --signer-email "traveler@example.com" --ip "1.1.1.1" --signed-at "2026-01-08T10:00:00Z"
   
   # IP 2 (30 phút sau)
   dxd-audit-cli log-event --file ../../sample.pdf --signer-email "traveler@example.com" --ip "8.8.8.8" --signed-at "2026-01-08T10:30:00Z"
   
   # IP 3 (50 phút sau)
   dxd-audit-cli log-event --file ../../sample.pdf --signer-email "traveler@example.com" --ip "203.0.113.1" --signed-at "2026-01-08T10:50:00Z"
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
Báo cáo sẽ hiển thị anomaly với nhãn:
- `"multi_ip": true`.
- Tóm tắt bất thường sẽ thể hiện số lượng IP khác nhau đã được sử dụng.
