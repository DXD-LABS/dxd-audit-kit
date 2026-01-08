#!/bin/bash

# Build CLI if not exists
go build -o ../../dxd-audit-cli ../../cmd/dxd-audit-cli/main.go

CLI="../../dxd-audit-cli"
SAMPLE_PDF="../../sample.pdf"

echo "1. Creating sample data..."
$CLI verify --file $SAMPLE_PDF
$CLI log-event --file $SAMPLE_PDF --signer-email "siem-test@example.com" --ip "1.2.3.4" --signed-at "2026-01-08T10:00:00Z"

DOC_ID=$($CLI report signer --email "siem-test@example.com" | grep "document_id" | head -n 1 | awk -F'"' '{print $4}')

echo "2. Exporting CSV..."
$CLI report document --document-id $DOC_ID --format csv > audit_logs.csv
echo "CSV exported to audit_logs.csv"

echo "3. Exporting NDJSON..."
$CLI report document --document-id $DOC_ID --format ndjson > audit_logs.ndjson
echo "NDJSON exported to audit_logs.ndjson"

echo "Sample NDJSON content:"
head -n 5 audit_logs.ndjson
