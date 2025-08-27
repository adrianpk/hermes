#!/bin/bash

# ==============================================================================
# API Test Script for: Role CRUD
# ==============================================================================
#
# Description:
#   This script tests the full CRUD (Create, Read, Update, Delete)
#   functionality for the /roles endpoint of the auth API.
#
#   It performs the following steps:
#   1. Fetches all existing roles (initial state).
#   2. Creates a new test role.
#   3. Captures the ID of the newly created role.
#   4. Fetches the specific role by its new ID to verify creation.
#   5. Updates the role's description.
#   6. Deletes the role.
#   7. Fetches all roles again to ensure it was deleted.
#
# Usage:
#   Ensure the Hermes server is running, then execute this script from the
#   project root:
#   ./scripts/curl/auth/role.sh
#
# Requirements:
#   - curl
#   - jq (for parsing JSON responses)
#
# ==============================================================================

# --- Configuration ---
BASE_URL="http://localhost:8081/api/v1/auth"
HEADERS="-H \"Content-Type: application/json\""
ROLE_NAME="api-test-role-$$" # Use process ID for uniqueness
ROLE_DESC="A role created via API test script."
UPDATED_DESC="This description has been updated."

# --- Helper Functions ---
function print_header() {
    echo ""
    echo "--- $1 ---"
}

# --- Test Execution ---

# 1. Get all roles (initial state)
print_header "1. GET /roles (Initial State)"
curl -s -X GET "$BASE_URL/roles" | jq .

# 2. Create a new role
print_header "2. POST /roles (Create New Role)"
CREATE_RESPONSE=$(curl -s -X POST $HEADERS -d "{\"name\": \"$ROLE_NAME\", \"description\": \"$ROLE_DESC\"}" "$BASE_URL/roles")
echo "$CREATE_RESPONSE" | jq .

# 3. Capture the new role's ID
ROLE_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.role.id')

if [ -z "$ROLE_ID" ] || [ "$ROLE_ID" == "null" ]; then
    echo "Error: Failed to create role or capture ID. Aborting."
    exit 1
fi

echo "Captured Role ID: $ROLE_ID"

# 4. Get the specific role by ID
print_header "4. GET /roles/{id} (Verify Creation)"
curl -s -X GET "$BASE_URL/roles/$ROLE_ID" | jq .

# 5. Update the role
print_header "5. PUT /roles/{id} (Update Role)"
curl -s -X PUT $HEADERS -d "{\"name\": \"$ROLE_NAME\", \"description\": \"$UPDATED_DESC\"}" "$BASE_URL/roles/$ROLE_ID" | jq .

# 6. Delete the role
print_header "6. DELETE /roles/{id} (Delete Role)"
curl -s -X DELETE "$BASE_URL/roles/$ROLE_ID" | jq .

# 7. Get all roles (final state)
print_header "7. GET /roles (Final State - Verify Deletion)"
curl -s -X GET "$BASE_URL/roles" | jq .

echo ""
echo "Role CRUD test script finished."
