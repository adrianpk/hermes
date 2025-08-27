#!/bin/bash

# Description: Tests the Content API endpoints (CRUD cycle).
# Usage: ./content.sh

# --- Configuration ---
source "$(dirname "$0")/_config.sh"
RESOURCE="contents"
RANDOM_SUFFIX=$((RANDOM % 10000))
CONTENT_HEADING="API Test Content $RANDOM_SUFFIX"

# --- Helper Function to manage dependencies ---
setup_dependencies() {
    echo "--- Setting up dependencies (Layout & Section) ---"
    # Create Layout
    LAYOUT_PAYLOAD="{\"name\": \"dep-layout-$RANDOM_SUFFIX\", \"description\": \"Dependency for content test\", \"code\": \"<p>{{ .Body }}</p>\"}"
    layout_response=$(curl -s -X POST "$BASE_URL/layouts" -H "Content-Type: application/json" -d "$LAYOUT_PAYLOAD")
    LAYOUT_ID=$(echo "$layout_response" | jq -r '.data.layout.id')
    if [ -z "$LAYOUT_ID" ] || [ "$LAYOUT_ID" == "null" ]; then
        echo "Failed to create dependency layout. Aborting."
        exit 1
    fi
    echo "Created dependency Layout ID: $LAYOUT_ID"
    export LAYOUT_ID_CLEANUP=$LAYOUT_ID

    # Create Section
    SECTION_PAYLOAD="{\"name\": \"dep-section-$RANDOM_SUFFIX\", \"description\": \"Dependency for content test\", \"path\": \"/dep-section\", \"layout_id\": \"$LAYOUT_ID\"}"
    section_response=$(curl -s -X POST "$BASE_URL/sections" -H "Content-Type: application/json" -d "$SECTION_PAYLOAD")
    SECTION_ID=$(echo "$section_response" | jq -r '.data.section.id')
    if [ -z "$SECTION_ID" ] || [ "$SECTION_ID" == "null" ]; then
        echo "Failed to create dependency section. Aborting."
        exit 1
    fi
    echo "Created dependency Section ID: $SECTION_ID"
    export SECTION_ID_CLEANUP=$SECTION_ID
}

cleanup_dependencies() {
    if [ -n "$SECTION_ID_CLEANUP" ]; then
        echo "--- Cleaning up dependency (Section ID: $SECTION_ID_CLEANUP) ---"
        curl -s -X DELETE "$BASE_URL/sections/$SECTION_ID_CLEANUP" > /dev/null
    fi
    if [ -n "$LAYOUT_ID_CLEANUP" ]; then
        echo "--- Cleaning up dependency (Layout ID: $LAYOUT_ID_CLEANUP) ---"
        curl -s -X DELETE "$BASE_URL/layouts/$LAYOUT_ID_CLEANUP" > /dev/null
    fi
}

# --- Main Functions ---
create_content() {
    local section_id=$1
    echo "--- 1. POST /$RESOURCE (Create New Content) ---" >&2 # Write to stderr
    PAYLOAD=$(cat <<EOF
{
    "heading": "$CONTENT_HEADING",
    "body": "This is test content created via API.",
    "section_id": "$section_id",
    "status": "published"
}
EOF
)
    response=$(curl -s -X POST "$BASE_URL/$RESOURCE" -H "Content-Type: application/json" -d "$PAYLOAD")
    echo "$response"
}

get_content() {
    local id=$1
    echo "--- GET /$RESOURCE/{id} (Verify Creation) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE/$id" | jq .
}

list_content() {
    echo "--- GET /$RESOURCE (List Content) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE" | jq .
}

update_content() {
    local id=$1
    local section_id=$2
    echo "--- PUT /$RESOURCE/{id} (Update Content) ---"
    PAYLOAD=$(cat <<EOF
{
    "heading": "$CONTENT_HEADING (Updated)",
    "body": "This is updated test content.",
    "section_id": "$section_id",
    "status": "draft"
}
EOF
)
    curl -s -X PUT "$BASE_URL/$RESOURCE/$id" -H "Content-Type: application/json" -d "$PAYLOAD" | jq .
}

delete_content() {
    local id=$1
    echo "--- DELETE /$RESOURCE/{id} (Delete Content) ---"
    curl -s -X DELETE "$BASE_URL/$RESOURCE/$id" | jq .
}

# --- Main Execution ---
trap cleanup_dependencies EXIT

echo "--- Running Full Content CRUD Test Cycle ---"
setup_dependencies
SECTION_ID=$SECTION_ID_CLEANUP

response_data=$(create_content "$SECTION_ID")
CONTENT_ID=$(echo "$response_data" | jq -r '.data.content.id')

if [ -z "$CONTENT_ID" ] || [ "$CONTENT_ID" == "null" ]; then
    echo "Failed to create content or capture ID. Aborting."
    exit 1
fi
echo "Captured Content ID: $CONTENT_ID"
sleep 1

get_content "$CONTENT_ID"
sleep 1

list_content
sleep 1

update_content "$CONTENT_ID" "$SECTION_ID"
sleep 1

delete_content "$CONTENT_ID"
echo "--- Content CRUD Test Cycle Finished ---"
