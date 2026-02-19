# Brainstorming Skill Process Redesign Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rewrite `brainstorming` skill text into an analysis-driven, convergence-oriented workflow without changing Telegram runtime behavior.

**Architecture:** Keep runtime/binaries untouched and change only skill documentation behavior contract. Add a lightweight validator script in this repo to enforce mandatory sections/rules in `SKILL.md`, then rewrite the skill and verify via rule checks plus scenario walkthrough.

**Tech Stack:** Markdown, shell (`bash`, `rg`), Git

---

### Task 1: Create Skill Rule Validator

**Files:**
- Create: `scripts/validate_brainstorming_skill.sh`
- Test: `scripts/validate_brainstorming_skill.sh`

**Step 1: Write the failing validator scaffold**

```bash
cat > scripts/validate_brainstorming_skill.sh <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
TARGET="${1:-/home/cao/.codex/skills/brainstorming/SKILL.md}"
test -f "$TARGET"
echo "TODO"
exit 1
EOF
chmod +x scripts/validate_brainstorming_skill.sh
```

**Step 2: Run validator to verify it fails**

Run: `scripts/validate_brainstorming_skill.sh`  
Expected: `exit code 1`

**Step 3: Implement minimal rule checks**

```bash
cat > scripts/validate_brainstorming_skill.sh <<'EOF'
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
EOF
chmod +x scripts/validate_brainstorming_skill.sh
```

**Step 4: Run validator to verify it passes after rewrite (currently may fail)**

Run: `scripts/validate_brainstorming_skill.sh`  
Expected: `OK: brainstorming skill rules present` after Task 2 is done.

**Step 5: Commit**

```bash
git add scripts/validate_brainstorming_skill.sh
git commit -m "chore: add brainstorming skill rule validator"
```

### Task 2: Rewrite Brainstorming Skill Document

**Files:**
- Modify: `/home/cao/.codex/skills/brainstorming/SKILL.md`
- Test: `scripts/validate_brainstorming_skill.sh`

**Step 1: Write a backup copy**

Run: `cp /home/cao/.codex/skills/brainstorming/SKILL.md /tmp/brainstorming-SKILL.backup.md`  
Expected: backup file exists.

**Step 2: Rewrite SKILL.md with new structure**

Required sections to include:
- `Overview`
- `State Machine` with entry/actions/exit for 6 states
- `Six-Dimension Analysis Ledger`
- `Question Protocol`
- `Option Management`
- `Soft-Gate Recovery`
- `Completion Contract`
- `After Design` with writing-plans as the only next skill

**Step 3: Run validator and verify pass**

Run: `scripts/validate_brainstorming_skill.sh /home/cao/.codex/skills/brainstorming/SKILL.md`  
Expected: `OK: brainstorming skill rules present`

**Step 4: Manual walkthrough for 3 scenarios**

Run:
- simple request scenario
- conflicting constraints scenario
- changing preference scenario

Expected:
- one question per turn
- ledger updated each turn
- candidate options promoted/demoted/eliminated
- convergence criteria reached before design finalization

**Step 5: Commit**

```bash
git add /home/cao/.codex/skills/brainstorming/SKILL.md
git commit -m "docs: redesign brainstorming workflow around analysis-driven convergence"
```

### Task 3: Validate Compatibility With Existing Telegram Workflow

**Files:**
- Modify: `docs/plans/2026-02-19-brainstorming-skill-process-redesign-design.md` (append verification notes)
- Test: `scripts/run_telegram_brainstorming.sh`

**Step 1: Run a smoke test with existing Telegram binary path**

Run: `scripts/run_telegram_brainstorming.sh --prompt "A/B/C test"`  
Expected: prompt sent via Telegram, terminal remains status-oriented.

**Step 2: Confirm no runtime code changes are needed**

Run: `git diff --name-only`  
Expected: only skill/doc/script files changed.

**Step 3: Record verification notes in design doc**

Append:
- smoke test date/time
- result summary
- compatibility confirmation

**Step 4: Commit**

```bash
git add docs/plans/2026-02-19-brainstorming-skill-process-redesign-design.md
git commit -m "docs: record compatibility verification for brainstorming skill redesign"
```
