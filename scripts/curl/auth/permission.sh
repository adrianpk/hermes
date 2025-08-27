#!/bin/bash

# ==============================================================================
# API Test Script for: Permission CRUD
# ==============================================================================
#
# Description:
#   This script tests the full CRUD (Create, Read, Update, Delete)
#   functionality for the /permissions endpoint of the auth API.
#
# Usage:
#   Ensure the Hermes server is running, then execute this script from the
#   project root:
#   ./scripts/curl/auth/permission.sh
#
# Requirements:
#   - curl
#   - jq
#
# ==============================================================================

# --- Configuration ---
BASE_URL="http://localhost:8081/api/v1/auth"
HEADERS="-H \"Content-Type: application/json\""
PERM_NAME="test:permission-$$" # Use process ID for uniqueness
PERM_DESC="A permission created via API test script."
UPDATED_DESC="This permission description has been updated."

# --- Helper Functions ---
function print_header() {
    echo ""
    echo "--- $1 ---"
}

# --- Test Execution ---

# 1. Get all permissions (initial state)
print_header "1. GET /permissions (Initial State)"
curl -s -X GET "$BASE_URL/permissions" | jq .

# 2. Create a new permission
print_header "2. POST /permissions (Create New Permission)"
CREATE_RESPONSE=$(curl -s -X POST $HEADERS -d "{\"name\": \"$PERM_NAME\", \"description\": \"$PERM_DESC\"}" "$BASE_URL/permissions")
echo "$CREATE_RESPONSE" | jq .

# 3. Capture the new permission's ID
PERM_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.permission.id')

if [ -z "$PERM_ID" ] || [ "$PERM_ID" == "null" ]; then
    echo "Error: Failed to create permission or capture ID. Aborting."
    exit 1
fi

echo "Captured Permission ID: $PERM_ID"

# 4. Get the specific permission by ID
print_header "4. GET /permissions/{id} (Verify Creation)"
curl -s -X GET "$BASE_URL/permissions/$PERM_ID" | jq .

# 5. Update the permission
print_header "5. PUT /permissions/{id} (Update Permission)"
curl -s -X PUT $HEADERS -d "{\"name\": \"$PERM_NAME\", \"description\": \"$UPDATED_DESC\"}" "$BASE_URL/permissions/$PERM_ID" | jq .

# 6. Delete the permission
print_header "6. DELETE /permissions/{id} (Delete Permission)"
curl -s -X DELETE "$BASE_URL/permissions/$PERM_ID" | jq .

# 7. Get all permissions (final state)
print_header "7. GET /permissions (Final State - Verify Deletion)"
curl -s -X GET "$BASE_URL/permissions" | jq .

echo ""
echo "Permission CRUD test script finished."
