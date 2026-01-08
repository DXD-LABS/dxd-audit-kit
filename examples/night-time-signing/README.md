# Scenario 1: Night-time Signing (Ký khuya bất thường)

Mục đích: Minh họa rule `night_time` phát hiện các sự kiện ký ngoài giờ hành chính (22:00 - 06:00).

## Các bước thực hiện

1. **Verify tài liệu**: Đăng ký tài liệu vào hệ thống.
   ```bash
   dxd-audit-cli verify --file ../../sample.pdf
   ```
   *Lưu ý: Lấy Document Hash từ output.*

2. **Log sự kiện ký lúc nửa đêm**: Mô phỏng việc ký tài liệu lúc 23:30.
   ```bash
   dxd-audit-cli log-event --file ../../sample.pdf --signer-email "night-owl@example.com" --ip "1.2.3.4" --signed-at "2026-01-08T23:30:00Z"
   ```

3. **Phân tích tài liệu**: Chạy bộ máy phân tích để tìm anomalies.
   ```bash
   dxd-audit-cli analyze document --document-id <DOCUMENT_ID>
   ```

4. **Kiểm tra báo cáo**: Xem kết quả phân tích.
   ```bash
   dxd-audit-cli report document --document-id <DOCUMENT_ID> --format json
   ```

## Kỳ vọng kết quả
Trong phần `anomalies` của báo cáo, bạn sẽ thấy:
- Một entry với nhãn `"night_time": true`.
- Điểm số `score` lớn hơn 0.
- `anomaly_summary.anomaly_count` > 0.
