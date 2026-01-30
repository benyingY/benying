#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKFLOW="$ROOT_DIR/.github/workflows/build-images.yml"

if [ ! -f "$WORKFLOW" ]; then
  echo "Missing workflow: $WORKFLOW" >&2
  exit 1
fi

if command -v rg >/dev/null 2>&1; then
  SEARCH=(rg -F -q)
else
  SEARCH=(grep -F -q)
fi

assert_contains() {
  local pattern="$1"
  if ! "${SEARCH[@]}" "$pattern" "$WORKFLOW"; then
    echo "Expected pattern not found: $pattern" >&2
    exit 1
  fi
}

assert_contains "push:"
assert_contains "docker/metadata-action"
assert_contains "type=sha"
assert_contains "docker/build-push-action"
assert_contains "push: true"
assert_contains "services/\${{ matrix.service }}"

echo "Step 2 checks passed."
