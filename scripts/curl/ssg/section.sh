#!/bin/bash

# Description: Tests the Section API endpoints (CRUD cycle).
# Usage: ./section.sh

# --- Configuration ---
source "$(dirname "$0")/_config.sh"
RESOURCE="sections"
RANDOM_SUFFIX=$((RANDOM % 10000))
SECTION_NAME="api-test-section-$RANDOM_SUFFIX"

# --- Helper Function to manage dependencies ---
setup_dependency() {
    echo "--- Setting up dependency (Layout) ---"
    LAYOUT_PAYLOAD="{\"name\": \"dep-layout-$RANDOM_SUFFIX\", \"description\": \"Dependency for section test\", \"code\": \"<p>Test</p>\"}"
    response=$(curl -s -X POST "$BASE_URL/layouts" -H "Content-Type: application/json" -d "$LAYOUT_PAYLOAD")
    LAYOUT_ID=$(echo "$response" | jq -r '.data.layout.id')
    if [ -z "$LAYOUT_ID" ] || [ "$LAYOUT_ID" == "null" ]; then
        echo "Failed to create dependency layout. Aborting."
        exit 1
    fi
    echo "Created dependency Layout ID: $LAYOUT_ID"
    # Export for cleanup
    export LAYOUT_ID_CLEANUP=$LAYOUT_ID
}

cleanup_dependency() {
    if [ -n "$LAYOUT_ID_CLEANUP" ]; then
        echo "--- Cleaning up dependency (Layout ID: $LAYOUT_ID_CLEANUP) ---"
        curl -s -X DELETE "$BASE_URL/layouts/$LAYOUT_ID_CLEANUP" > /dev/null
    fi
}

# --- Main Functions ---
create_section() {
    local layout_id=$1
    echo "--- 1. POST /$RESOURCE (Create New Section) ---" >&2 # Write to stderr
    PAYLOAD=$(cat <<EOF
{
    "name": "$SECTION_NAME",
    "description": "A test section created via API.",
    "path": "/$SECTION_NAME",
    "layout_id": "$layout_id"
}
EOF
)
    response=$(curl -s -X POST "$BASE_URL/$RESOURCE" -H "Content-Type: application/json" -d "$PAYLOAD")
    echo "$response"
}

get_section() {
    local id=$1
    echo "--- GET /$RESOURCE/{id} (Verify Creation) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE/$id" | jq .
}

list_sections() {
    echo "--- GET /$RESOURCE (List Sections) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE" | jq .
}

update_section() {
    local id=$1
    local layout_id=$2
    echo "--- PUT /$RESOURCE/{id} (Update Section) ---"
    PAYLOAD=$(cat <<EOF
{
    "name": "$SECTION_NAME (Updated)",
    "description": "An updated test section.",
    "path": "/$SECTION_NAME-updated",
    "layout_id": "$layout_id"
}
EOF
)
    curl -s -X PUT "$BASE_URL/$RESOURCE/$id" -H "Content-Type: application/json" -d "$PAYLOAD" | jq .
}

delete_section() {
    local id=$1
    echo "--- DELETE /$RESOURCE/{id} (Delete Section) ---"
    curl -s -X DELETE "$BASE_URL/$RESOURCE/$id" | jq .
}

# --- Main Execution ---
trap cleanup_dependency EXIT

echo "--- Running Full Section CRUD Test Cycle ---"
setup_dependency
LAYOUT_ID=$LAYOUT_ID_CLEANUP

response_data=$(create_section "$LAYOUT_ID")
SECTION_ID=$(echo "$response_data" | jq -r '.data.section.id')

if [ -z "$SECTION_ID" ] || [ "$SECTION_ID" == "null" ]; then
    echo "Failed to create section or capture ID. Aborting."
    exit 1
fi
echo "Captured Section ID: $SECTION_ID"
sleep 1

get_section "$SECTION_ID"
sleep 1

list_sections
sleep 1

update_section "$SECTION_ID" "$LAYOUT_ID"
sleep 1

delete_section "$SECTION_ID"
echo "--- Section CRUD Test Cycle Finished ---"
