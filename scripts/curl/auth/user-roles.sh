#!/bin/bash

# ==============================================================================
# API Test Script for: User-Roles Relationship
# ==============================================================================
#
# Description:
#   This script tests the API endpoints for managing the relationship
#   between Users and Roles.
#
#   It performs the following steps:
#   1. Fetches the list of all available users.
#   2. Captures the ID of the first user from the list.
#   3. Fetches the list of all available roles.
#   4. Captures the ID of a role that is NOT 'admin' or 'user' (e.g., 'editor').
#   5. Lists the initial roles assigned to the user.
#   6. Assigns the captured role to the user.
#   7. Verifies the assignment by listing the user's roles again.
#   8. Removes the role from the user.
#   9. Verifies the removal by listing the user's roles one last time.
#
# Usage:
#   Ensure the Hermes server is running, then execute this script from the
#   project root:
#   ./scripts/curl/auth/user-roles.sh
#
# Requirements:
#   - curl
#   - jq
#
# ==============================================================================

# --- Configuration ---
BASE_URL="http://localhost:8081/api/v1/auth"
HEADERS="-H \"Content-Type: application/json\""

# --- Helper Functions ---
function print_header() {
    echo ""
    echo "--- $1 ---"
}

# --- Test Execution ---

# 1. Get all available users
print_header "1. GET /users (Fetch All Users)"
USERS_RESPONSE=$(curl -s -X GET "$BASE_URL/users")
echo "$USERS_RESPONSE" | jq .

# 2. Capture the first user's ID
USER_ID=$(echo "$USERS_RESPONSE" | jq -r '.data[0].id')
if [ -z "$USER_ID" ] || [ "$USER_ID" == "null" ]; then
    echo "Error: Failed to fetch users or capture a user ID. Aborting."
    exit 1
fi
echo "Captured User ID for testing: $USER_ID"

# 3. Get all available roles
print_header "3. GET /roles (Fetch All Roles)"
ROLES_RESPONSE=$(curl -s -X GET "$BASE_URL/roles")
echo "$ROLES_RESPONSE" | jq .

# 4. Capture the 'editor' role's ID
ROLE_ID=$(echo "$ROLES_RESPONSE" | jq -r '.data[] | select(.name == "editor").id')
if [ -z "$ROLE_ID" ] || [ "$ROLE_ID" == "null" ]; then
    echo "Error: Failed to find the 'editor' role. Aborting."
    exit 1
fi
echo "Captured Role ID for testing ('editor'): $ROLE_ID"

# 5. Get initial roles for the user
print_header "5. GET /users/{userId}/roles (Initial State)"
curl -s -X GET "$BASE_URL/users/$USER_ID/roles" | jq .

# 6. Assign the role to the user
print_header "6. POST /users/{userId}/roles (Assign Role)"
curl -s -X POST $HEADERS -d "{\"role_id\": \"$ROLE_ID\"}" "$BASE_URL/users/$USER_ID/roles" | jq .

# 7. Verify the role was assigned
print_header "7. GET /users/{userId}/roles (Verify Assignment)"
curl -s -X GET "$BASE_URL/users/$USER_ID/roles" | jq .

# 8. Remove the role from the user
print_header "8. DELETE /users/{userId}/roles/{roleId} (Remove Role)"
curl -s -X DELETE "$BASE_URL/users/$USER_ID/roles/$ROLE_ID" | jq .

# 9. Verify the role was removed
print_header "9. GET /users/{userId}/roles (Verify Removal)"
curl -s -X GET "$BASE_URL/users/$USER_ID/roles" | jq .

echo ""
echo "User-Roles relationship test script finished."
