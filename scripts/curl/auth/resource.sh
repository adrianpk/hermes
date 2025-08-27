#!/bin/bash

# ==============================================================================
# API Test Script for: Resource CRUD
# ==============================================================================
#
# Description:
#   This script tests the full CRUD (Create, Read, Update, Delete)
#   functionality for the /resources endpoint of the auth API.
#
# Usage:
#   Ensure the Hermes server is running, then execute this script from the
#   project root:
#   ./scripts/curl/auth/resource.sh
#
# Requirements:
#   - curl
#   - jq
#
# ==============================================================================

# --- Configuration ---
BASE_URL="http://localhost:8081/api/v1/auth"
HEADERS="-H \"Content-Type: application/json\""
RES_NAME="test-resource-$$" # Use process ID for uniqueness
RES_DESC="A resource created via API test script."
RES_KIND="entity"
UPDATED_DESC="This resource description has been updated."

# --- Helper Functions ---
function print_header() {
    echo ""
    echo "--- $1 ---"
}

# --- Test Execution ---

# 1. Get all resources (initial state)
print_header "1. GET /resources (Initial State)"
curl -s -X GET "$BASE_URL/resources" | jq .

# 2. Create a new resource
print_header "2. POST /resources (Create New Resource)"
CREATE_RESPONSE=$(curl -s -X POST $HEADERS -d "{\"name\": \"$RES_NAME\", \"description\": \"$RES_DESC\", \"kind\": \"$RES_KIND\"}" "$BASE_URL/resources")
echo "$CREATE_RESPONSE" | jq .

# 3. Capture the new resource's ID
RES_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.resource.id')

if [ -z "$RES_ID" ] || [ "$RES_ID" == "null" ]; then
    echo "Error: Failed to create resource or capture ID. Aborting."
    exit 1
fi

echo "Captured Resource ID: $RES_ID"

# 4. Get the specific resource by ID
print_header "4. GET /resources/{id} (Verify Creation)"
curl -s -X GET "$BASE_URL/resources/$RES_ID" | jq .

# 5. Update the resource
print_header "5. PUT /resources/{id} (Update Resource)"
curl -s -X PUT $HEADERS -d "{\"name\": \"$RES_NAME\", \"description\": \"$UPDATED_DESC\", \"kind\": \"$RES_KIND\"}" "$BASE_URL/resources/$RES_ID" | jq .

# 6. Delete the resource
print_header "6. DELETE /resources/{id} (Delete Resource)"
curl -s -X DELETE "$BASE_URL/resources/$RES_ID" | jq .

# 7. Get all resources (final state)
print_header "7. GET /resources (Final State - Verify Deletion)"
curl -s -X GET "$BASE_URL/resources" | jq .

echo ""
echo "Resource CRUD test script finished."
