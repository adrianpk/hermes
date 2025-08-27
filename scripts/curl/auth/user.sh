#!/bin/bash

# ==============================================================================
# API Test Script for: User CRUD
# ==============================================================================
#
# Description:
#   This script tests the full CRUD (Create, Read, Update, Delete)
#   functionality for the /users endpoint of the auth API.
#
#   It performs the following steps:
#   1. Fetches all existing users (initial state).
#   2. Creates a new test user.
#   3. Captures the ID of the newly created user.
#   4. Fetches the specific user by their new ID to verify creation.
#   5. Updates the user's name.
#   6. Deletes the user.
#   7. Fetches all users again to ensure it was deleted.
#
# Usage:
#   Ensure the Hermes server is running, then execute this script from the
#   project root:
#   ./scripts/curl/auth/user.sh
#
# Requirements:
#   - curl
#   - jq
#
# ==============================================================================

# --- Configuration ---
BASE_URL="http://localhost:8081/api/v1/auth"
HEADERS="-H \"Content-Type: application/json\""
USERNAME="api-test-user-$$" # Use process ID for uniqueness
EMAIL="testuser-$$@example.com"
NAME="API Test User"
PASSWORD="a-secure-password"
UPDATED_NAME="API Test User (Updated)"

# --- Helper Functions ---
function print_header() {
    echo ""
    echo "--- $1 ---"
}

# --- Test Execution ---

# 1. Get all users (initial state)
print_header "1. GET /users (Initial State)"
curl -s -X GET "$BASE_URL/users" | jq .

# 2. Create a new user
print_header "2. POST /users (Create New User)"
CREATE_PAYLOAD="{\"username\": \"$USERNAME\", \"email\": \"$EMAIL\", \"name\": \"$NAME\", \"password\": \"$PASSWORD\"}"
CREATE_RESPONSE=$(curl -s -X POST $HEADERS -d "$CREATE_PAYLOAD" "$BASE_URL/users")
echo "$CREATE_RESPONSE" | jq .

# 3. Capture the new user's ID
USER_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.user.id')

if [ -z "$USER_ID" ] || [ "$USER_ID" == "null" ]; then
    echo "Error: Failed to create user or capture ID. Aborting."
    exit 1
fi

echo "Captured User ID: $USER_ID"

# 4. Get the specific user by ID
print_header "4. GET /users/{id} (Verify Creation)"
curl -s -X GET "$BASE_URL/users/$USER_ID" | jq .

# 5. Update the user
print_header "5. PUT /users/{id} (Update User)"
UPDATE_PAYLOAD="{\"name\": \"$UPDATED_NAME\"}"
curl -s -X PUT $HEADERS -d "$UPDATE_PAYLOAD" "$BASE_URL/users/$USER_ID" | jq .

# 6. Delete the user
print_header "6. DELETE /users/{id} (Delete User)"
curl -s -X DELETE "$BASE_URL/users/$USER_ID" | jq .

# 7. Get all users (final state)
print_header "7. GET /users (Final State - Verify Deletion)"
curl -s -X GET "$BASE_URL/users" | jq .

echo ""
echo "User CRUD test script finished."
