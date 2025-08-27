#!/bin/bash

# Description: Tests the Layout API endpoints (CRUD cycle).
# Usage: ./layout.sh [create|get|list|update|delete] [ID]

# --- Configuration ---
source "$(dirname "$0")/_config.sh"
RESOURCE="layouts"
RANDOM_SUFFIX=$((RANDOM % 10000))
LAYOUT_NAME="api-test-layout-$RANDOM_SUFFIX"

# --- Functions ---
create_layout() {
    echo "--- 1. POST /$RESOURCE (Create New Layout) ---" >&2 # Write to stderr
    PAYLOAD=$(cat <<EOF
{
    "name": "$LAYOUT_NAME",
    "description": "A test layout created via API.",
    "code": "<html><body>{{ .Content }}</body></html>"
}
EOF
)
    response=$(curl -s -X POST "$BASE_URL/$RESOURCE" -H "Content-Type: application/json" -d "$PAYLOAD")
    echo "$response"
}

get_layout() {
    local id=$1
    echo "--- GET /$RESOURCE/{id} (Verify Creation) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE/$id" | jq .
}

list_layouts() {
    echo "--- GET /$RESOURCE (List Layouts) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE" | jq .
}

update_layout() {
    local id=$1
    echo "--- PUT /$RESOURCE/{id} (Update Layout) ---"
    PAYLOAD=$(cat <<EOF
{
    "name": "$LAYOUT_NAME (Updated)",
    "description": "An updated test layout.",
    "code": "<html><head><title>Updated</title></head><body>{{ .Content }}</body></html>"
}
EOF
)
    curl -s -X PUT "$BASE_URL/$RESOURCE/$id" -H "Content-Type: application/json" -d "$PAYLOAD" | jq .
}

delete_layout() {
    local id=$1
    echo "--- DELETE /$RESOURCE/{id} (Delete Layout) ---"
    curl -s -X DELETE "$BASE_URL/$RESOURCE/$id" | jq .
}

# --- Main Execution ---
COMMAND=$1
ID=$2

case "$COMMAND" in
    create)
        create_layout
        ;;
    get)
        [ -z "$ID" ] && { echo "Usage: $0 get <ID>"; exit 1; }
        get_layout "$ID"
        ;;
    list)
        list_layouts
        ;;
    update)
        [ -z "$ID" ] && { echo "Usage: $0 update <ID>"; exit 1; }
        update_layout "$ID"
        ;;
    delete)
        [ -z "$ID" ] && { echo "Usage: $0 delete <ID>"; exit 1; }
        delete_layout "$ID"
        ;;
    *)
        echo "--- Running Full Layout CRUD Test Cycle ---"
        response_data=$(create_layout)
        LAYOUT_ID=$(echo "$response_data" | jq -r '.data.layout.id')

        if [ -z "$LAYOUT_ID" ] || [ "$LAYOUT_ID" == "null" ]; then
            echo "Failed to create layout or capture ID. Aborting."
            exit 1
        fi
        echo "Captured Layout ID: $LAYOUT_ID"
        sleep 1

        get_layout "$LAYOUT_ID"
        sleep 1

        list_layouts
        sleep 1

        update_layout "$LAYOUT_ID"
        sleep 1

        delete_layout "$LAYOUT_ID"
        echo "--- Layout CRUD Test Cycle Finished ---"
        ;;
esac
