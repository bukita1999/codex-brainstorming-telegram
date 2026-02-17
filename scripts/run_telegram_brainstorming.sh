#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if [ "$#" -eq 0 ]; then
  echo "Usage: $0 \"prompt text\""
  exit 1
fi

GOCACHE=/tmp/go-build go run ./cmd/telegram-brainstorming --env .env "$@"
