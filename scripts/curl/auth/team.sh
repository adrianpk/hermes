#!/bin/bash

# Load environment variables
source .env

# Set base URL
BASE_URL="http://localhost:8080/api/v1/auth"

# Function to print a header
print_header() {
    echo ""
    echo "================================================================================"
    echo "$1"
    echo "================================================================================"
}

# Get the first Org ID to use in tests
ORG_ID=$(curl -s -X GET "$BASE_URL/orgs" | jq -r '.data[0].id')
if [ -z "$ORG_ID" ] || [ "$ORG_ID" == "null" ]; then
    echo "Error: Could not retrieve an ORG_ID. Please create an organization first."
    exit 1
fi
echo "Using ORG_ID: $ORG_ID"


# Test case 1: Get all teams for an organization
print_header "GET /orgs/{orgId}/teams - Get all teams"
curl -X GET "$BASE_URL/orgs/$ORG_ID/teams" | jq .


# Test case 2: Create a new team
print_header "POST /orgs/{orgId}/teams - Create a new team"
TEAM_PAYLOAD='{
    "name": "Team Rocket",
    "shortDescription": "A team that is always blasting off again.",
    "description": "Team Rocket is a criminal organization from the Pokémon series."
}'
CREATED_TEAM=$(curl -s -X POST -H "Content-Type: application/json" -d "$TEAM_PAYLOAD" "$BASE_URL/orgs/$ORG_ID/teams")
echo $CREATED_TEAM | jq .
TEAM_ID=$(echo $CREATED_TEAM | jq -r '.data.team.id')


# Test case 3: Get a specific team
print_header "GET /teams/{teamId} - Get a specific team"
if [ -z "$TEAM_ID" ] || [ "$TEAM_ID" == "null" ]; then
    echo "Skipping GET /teams/{teamId} because TEAM_ID is not set."
else
    curl -X GET "$BASE_URL/teams/$TEAM_ID" | jq .
fi


# Test case 4: Update a team
print_header "PUT /teams/{teamId} - Update a team"
UPDATE_PAYLOAD='{
    "name": "Team Awesome",
    "shortDescription": "A truly awesome team.",
    "description": "The best team in the world."
}'
if [ -z "$TEAM_ID" ] || [ "$TEAM_ID" == "null" ]; then
    echo "Skipping PUT /teams/{teamId} because TEAM_ID is not set."
else
    curl -X PUT -H "Content-Type: application/json" -d "$UPDATE_PAYLOAD" "$BASE_URL/teams/$TEAM_ID" | jq .
fi


# Test case 5: Delete a team
print_header "DELETE /teams/{teamId} - Delete a team"
if [ -z "$TEAM_ID" ] || [ "$TEAM_ID" == "null" ]; then
    echo "Skipping DELETE /teams/{teamId} because TEAM_ID is not set."
else
    curl -X DELETE "$BASE_URL/teams/$TEAM_ID" | jq .
fi
