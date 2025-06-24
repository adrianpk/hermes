#!/bin/bash

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <new-module-path>"
  exit 1
fi

NEW_MODULE="$1"
TEMPLATE_REPO="https://github.com/aquamarinepk/hermes.git"
ORIGINAL_DIR="$(pwd)"

WORKDIR=$(mktemp -d)
git clone --depth 1 "$TEMPLATE_REPO" "$WORKDIR"
cd "$WORKDIR"

ORIGINAL_MODULE=$(grep "^module " go.mod | awk '{print $2}')
go mod edit -module "$NEW_MODULE"
find . -type f -name "*.go" -exec sed -i "s|$ORIGINAL_MODULE|$NEW_MODULE|g" {} +

PROJECT_NAME_RAW=$(basename "$NEW_MODULE")
PROJECT_NAME_UPPER=$(echo "$PROJECT_NAME_RAW" | tr '[:lower:]' '[:upper:]')

find . -type f ! -path "./.git/*" -exec sed -i \
  -e "s/hermes/$PROJECT_NAME_RAW/g" \
  -e "s/HERMES/$PROJECT_NAME_UPPER/g" \
  {} +

find . -maxdepth 1 \( -name "Makefile" -o -name "makefile" \) -exec sed -i "s|APP_NAME[[:space:]]*=[[:space:]]*.*|APP_NAME = $PROJECT_NAME_RAW|g" {} +

rm -rf .git

TARGET_DIR="$ORIGINAL_DIR/$PROJECT_NAME_RAW"
cp -R . "$TARGET_DIR"

cd "$ORIGINAL_DIR"
rm -rf "$WORKDIR"

