#!/bin/bash

# Build CLI if not exists
go build -o ../../dxd-audit-cli ../../cmd/dxd-audit-cli/main.go

CLI="../../dxd-audit-cli"
SAMPLE_PDF="../../sample.pdf"

echo "1. Verifying documents..."
$CLI verify --file $SAMPLE_PDF

echo "2. Logging events for a single signer across multiple documents and days..."
# Giả sử ta dùng lại cùng 1 file cho đơn giản, hoặc người dùng có thể dùng file khác
$CLI log-event --file $SAMPLE_PDF --signer-email "active-user@example.com" --ip "1.1.1.1" --signed-at "2026-01-07T10:00:00Z"
$CLI log-event --file $SAMPLE_PDF --signer-email "active-user@example.com" --ip "1.1.1.1" --signed-at "2026-01-08T09:00:00Z"
$CLI log-event --file $SAMPLE_PDF --signer-email "active-user@example.com" --ip "1.1.1.1" --signed-at "2026-01-08T15:00:00Z"

echo "3. Generating signer report..."
$CLI report signer --email "active-user@example.com" --from "2026-01-07" --to "2026-01-09" --format json
