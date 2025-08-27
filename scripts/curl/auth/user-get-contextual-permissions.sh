#!/bin/bash

# Description: Get a user with their global and all team-specific permissions.
# Usage: ./user-get-contextual-permissions.sh [USER_ID] [TOKEN]
# Example: ./user-get-contextual-permissions.sh 12345-abcde-67890 my_token

# --- Configuration ---
source "$(dirname "$0")/_config.sh"

# --- Command-line arguments ---
USER_ID=${1:?"Usage: $0 USER_ID [TOKEN]"}
TOKEN=${2:-"your_auth_token_here"} # Replace with a valid auth token or pass as second argument

# --- API Call ---
echo "Getting user $USER_ID with all global and team permissions..."
curl -X GET \
  "$BASE_URL/auth/users/$USER_ID?embed=team_permissions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq .

echo -e "\n--- Done ---"

