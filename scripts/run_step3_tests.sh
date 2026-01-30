#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Python tests (template)
PY_VENV="${PY_VENV:-/tmp/benying-venv}"
if [ ! -x "$PY_VENV/bin/python" ]; then
  python3 -m venv "$PY_VENV"
  "$PY_VENV/bin/python" -m pip install -r "$ROOT_DIR/services/python-service-template/requirements.txt"
fi
"$PY_VENV/bin/python" -m pytest "$ROOT_DIR/services/python-service-template/tests"

# Go tests (template)
export GOCACHE="${GOCACHE:-/tmp/go-build-cache}"
( cd "$ROOT_DIR/services/go-service-template" && go test ./... )
