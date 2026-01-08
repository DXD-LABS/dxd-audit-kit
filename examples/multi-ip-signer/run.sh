#!/bin/bash

# Build CLI if not exists
go build -o ../../dxd-audit-cli ../../cmd/dxd-audit-cli/main.go

CLI="../../dxd-audit-cli"
SAMPLE_PDF="../../sample.pdf"

echo "1. Verifying document..."
$CLI verify --file $SAMPLE_PDF

echo "2. Logging 3 events with different IPs in 1 hour..."
$CLI log-event --file $SAMPLE_PDF --signer-email "traveler@example.com" --ip "1.1.1.1" --signed-at "2026-01-08T10:00:00Z"
$CLI log-event --file $SAMPLE_PDF --signer-email "traveler@example.com" --ip "8.8.8.8" --signed-at "2026-01-08T10:30:00Z"
$CLI log-event --file $SAMPLE_PDF --signer-email "traveler@example.com" --ip "203.0.113.1" --signed-at "2026-01-08T10:50:00Z"

echo "3. Analyzing document..."
DOC_ID=$($CLI report signer --email "traveler@example.com" | grep "document_id" | head -n 1 | awk -F'"' '{print $4}')
echo "Found Document ID: $DOC_ID"
$CLI analyze document --document-id $DOC_ID

echo "4. Generating report..."
$CLI report document --document-id $DOC_ID
