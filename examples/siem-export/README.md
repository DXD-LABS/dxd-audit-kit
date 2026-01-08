# Scenario 5: SIEM-friendly Export (CSV & NDJSON)

Mục đích: Minh họa cách trích xuất dữ liệu từ Toolkit để đưa vào các hệ thống giám sát an ninh (SIEM) hoặc nền tảng quản lý log (Logstash, Splunk, v.v.)

## Các bước thực hiện

1. **Tạo dữ liệu mẫu**: Chạy một vài lệnh `verify` và `log-event` (có thể dùng lại dữ liệu từ các scenario khác).

2. **Xuất báo cáo định dạng CSV**:
   ```bash
   dxd-audit-cli report document --document-id <DOCUMENT_ID> --format csv > audit_logs.csv
   ```

3. **Xuất báo cáo định dạng NDJSON (SIEM-friendly)**:
   ```bash
   dxd-audit-cli report document --document-id <DOCUMENT_ID> --format ndjson > audit_logs.ndjson
   ```

## Tích hợp SIEM
Hầu hết các SIEM hiện đại đều hỗ trợ định dạng NDJSON (Newline Delimited JSON). Bạn có thể gửi dữ liệu này qua HTTP API.

Ví dụ gửi log đến một giả lập endpoint:
```bash
curl -X POST -H "Content-Type: application/x-ndjson" --data-binary @audit_logs.ndjson http://siem-collector.local/ingest
```

## Kỳ vọng kết quả
- File `audit_logs.csv`: Phù hợp để mở bằng Excel hoặc các công cụ phân tích dữ liệu bảng.
- File `audit_logs.ndjson`: Mỗi dòng là một JSON object hoàn chỉnh của một event, cực kỳ tối ưu cho việc ingest log tự động.
