#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "Usage: $0 <python|go> <service-name>" >&2
}

if [ $# -ne 2 ]; then
  usage
  exit 1
fi

LANGUAGE="$1"
NAME="$2"
TEMPLATE_DIR="services/${LANGUAGE}-service-template"
TARGET_DIR="services/${NAME}"

if [ ! -d "$TEMPLATE_DIR" ]; then
  echo "Template not found: $TEMPLATE_DIR" >&2
  exit 1
fi

if [ -e "$TARGET_DIR" ]; then
  echo "Target already exists: $TARGET_DIR" >&2
  exit 1
fi

cp -R "$TEMPLATE_DIR" "$TARGET_DIR"

if [ "$LANGUAGE" = "go" ]; then
  MODULE="example.com/${NAME}"
  perl -pi -e "s/MODULE_PLACEHOLDER/${MODULE}/g" \
    "$TARGET_DIR/go.mod" \
    "$TARGET_DIR/cmd/server/main.go"
fi

if [ "$LANGUAGE" = "python" ]; then
  perl -pi -e "s/service-template/${NAME}/g" \
    "$TARGET_DIR/app/main.py"
fi

echo "Created service at $TARGET_DIR"
