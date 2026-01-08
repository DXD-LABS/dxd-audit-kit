#!/bin/bash

# Build CLI if not exists
go build -o ../../dxd-audit-cli ../../cmd/dxd-audit-cli/main.go

CLI="../../dxd-audit-cli"
SAMPLE_PDF="../../sample.pdf"

echo "1. Verifying document..."
$CLI verify --file $SAMPLE_PDF

echo "2. Logging normal office events..."
$CLI log-event --file $SAMPLE_PDF --signer-email "staff@company.com" --ip "192.168.1.10" --signed-at "2026-01-08T09:00:00Z"
$CLI log-event --file $SAMPLE_PDF --signer-email "manager@company.com" --ip "192.168.1.15" --signed-at "2026-01-08T14:00:00Z"

echo "3. Analyzing document..."
DOC_ID=$($CLI report signer --email "staff@company.com" | grep "document_id" | head -n 1 | awk -F'"' '{print $4}')
echo "Found Document ID: $DOC_ID"
$CLI analyze document --document-id $DOC_ID

echo "4. Generating report (expecting 0 anomalies)..."
$CLI report document --document-id $DOC_ID
