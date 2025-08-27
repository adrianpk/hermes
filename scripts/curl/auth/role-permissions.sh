#!/bin/bash

# ==============================================================================
# API Test Script for: Role-Permissions Relationship
# ==============================================================================
#
# Description:
#   This script tests the API endpoints for managing the relationship
#   between Roles and Permissions.
#
#   It performs the following steps:
#   1. Creates a new temporary role to use for testing.
#   2. Fetches the list of all available permissions.
#   3. Captures the ID of the first permission from the list.
#   4. Lists the initial permissions assigned to the new role (should be none).
#   5. Assigns the captured permission to the role.
#   6. Verifies the assignment by listing the role's permissions again.
#   7. Removes the permission from the role.
#   8. Verifies the removal by listing the role's permissions one last time.
#   9. Cleans up by deleting the temporary role.
#
# Usage:
#   Ensure the Hermes server is running, then execute this script from the
#   project root:
#   ./scripts/curl/auth/role-permissions.sh
#
# Requirements:
#   - curl
#   - jq
#
# ==============================================================================

# --- Configuration ---
BASE_URL="http://localhost:8081/api/v1/auth"
HEADERS="-H \"Content-Type: application/json\""
ROLE_NAME="api-test-role-perm-$$"
ROLE_DESC="A temporary role for permission relationship testing."

# --- Helper Functions ---
function print_header() {
    echo ""
    echo "--- $1 ---"
}

# --- Test Execution ---

# 1. Create a new temporary role
print_header "1. POST /roles (Create Temporary Role)"
CREATE_RESPONSE=$(curl -s -X POST $HEADERS -d "{\"name\": \"$ROLE_NAME\", \"description\": \"$ROLE_DESC\"}" "$BASE_URL/roles")
ROLE_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.id')

if [ -z "$ROLE_ID" ] || [ "$ROLE_ID" == "null" ]; then
    echo "Error: Failed to create temporary role. Aborting."
    exit 1
fi
echo "Created temporary role with ID: $ROLE_ID"

# 2. Get all available permissions
print_header "2. GET /permissions (Fetch All Permissions)"
PERMISSIONS_RESPONSE=$(curl -s -X GET "$BASE_URL/permissions")
echo "$PERMISSIONS_RESPONSE" | jq .

# 3. Capture the first permission's ID
PERMISSION_ID=$(echo "$PERMISSIONS_RESPONSE" | jq -r '.data[0].id')
if [ -z "$PERMISSION_ID" ] || [ "$PERMISSION_ID" == "null" ]; then
    echo "Error: Failed to fetch permissions or capture a permission ID. Aborting."
    # Clean up the created role before exiting
    curl -s -X DELETE "$BASE_URL/roles/$ROLE_ID" > /dev/null
    exit 1
fi
echo "Captured Permission ID for testing: $PERMISSION_ID"

# 4. Get initial permissions for the new role
print_header "4. GET /roles/{roleId}/permissions (Initial State)"
curl -s -X GET "$BASE_URL/roles/$ROLE_ID/permissions" | jq .

# 5. Assign the permission to the role
print_header "5. POST /roles/{roleId}/permissions (Assign Permission)"
curl -s -X POST $HEADERS -d "{\"permission_id\": \"$PERMISSION_ID\"}" "$BASE_URL/roles/$ROLE_ID/permissions" | jq .

# 6. Verify the permission was assigned
print_header "6. GET /roles/{roleId}/permissions (Verify Assignment)"
curl -s -X GET "$BASE_URL/roles/$ROLE_ID/permissions" | jq .

# 7. Remove the permission from the role
print_header "7. DELETE /roles/{roleId}/permissions/{permissionId} (Remove Permission)"
curl -s -X DELETE "$BASE_URL/roles/$ROLE_ID/permissions/$PERMISSION_ID" | jq .

# 8. Verify the permission was removed
print_header "8. GET /roles/{roleId}/permissions (Verify Removal)"
curl -s -X GET "$BASE_URL/roles/$ROLE_ID/permissions" | jq .

# 9. Clean up: Delete the temporary role
print_header "9. DELETE /roles/{id} (Cleanup)"
curl -s -X DELETE "$BASE_URL/roles/$ROLE_ID" | jq .

echo ""
echo "Role-Permissions relationship test script finished."
