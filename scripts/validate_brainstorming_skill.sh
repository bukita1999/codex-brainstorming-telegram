#!/usr/bin/env bash
set -euo pipefail

TARGET="${1:-/home/cao/.codex/skills/brainstorming/SKILL.md}"
test -f "$TARGET" || { echo "missing: $TARGET"; exit 2; }

need() {
  local p="$1"
  rg -q --fixed-strings "$p" "$TARGET" || { echo "missing rule: $p"; exit 1; }
}

need "State Machine"
need "Six-Dimension Analysis Ledger"
need "Question Protocol"
need "Soft-Gate Recovery"
need "Completion Contract"
need "One question at a time"
need "Only the writing-plans skill"

echo "OK: brainstorming skill rules present"
