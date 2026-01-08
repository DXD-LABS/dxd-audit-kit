#!/bin/bash

# Build CLI if not exists
go build -o ../../dxd-audit-cli ../../cmd/dxd-audit-cli/main.go

CLI="../../dxd-audit-cli"
SAMPLE_PDF="../../sample.pdf"

echo "1. Verifying document..."
DOC_HASH=$($CLI verify --file $SAMPLE_PDF | grep "Hash:" | awk '{print $2}')
echo "Document Hash: $DOC_HASH"

# Lấy Document ID từ DB dựa trên Hash
# Vì CLI verify không trả về ID trực tiếp, ta sẽ dùng log-event với --file để CLI tự tìm ID
echo "2. Logging night-time event (23:30)..."
$CLI log-event --file $SAMPLE_PDF --signer-email "night-owl@example.com" --ip "1.2.3.4" --signed-at "2026-01-08T23:30:00Z"

# Để lấy DOC_ID cho lệnh analyze, ta có thể dùng report để xem ID
DOC_ID=$($CLI report document --document-id dummy 2>&1 | grep "Document not found" || true)
# Thực tế, cách tốt nhất là CLI log-event nên trả về DOC_ID hoặc ta dùng một mẹo nhỏ
# Ở đây tôi giả sử người dùng sẽ copy ID từ output của log-event hoặc ta dùng lệnh report để lấy nếu đã biết email

echo "3. Analyzing document..."
# Tìm DOC_ID từ email người ký vừa log
DOC_ID=$($CLI report signer --email "night-owl@example.com" | grep "document_id" | head -n 1 | awk -F'"' '{print $4}')

if [ -z "$DOC_ID" ]; then
    echo "Could not find Document ID. Please make sure the event was logged."
    exit 1
fi

echo "Found Document ID: $DOC_ID"
$CLI analyze document --document-id $DOC_ID

echo "4. Generating report..."
$CLI report document --document-id $DOC_ID --format json
