---
name: telegram-brainstorming
description: Use when the user explicitly requests brainstorming through Telegram chat and the workflow must run via packaged Linux binaries instead of terminal-only interaction.
---

# Telegram Brainstorming

## Purpose

Run the brainstorming workflow through Telegram while preserving core discipline from `brainstorming`: one question per turn, option-first prompts, and progressive requirement clarification.

This production version is binary-driven: the skill calls packaged Linux binaries from its own `bin/` folder.

## Trigger

Use only when the user explicitly asks for `telegram-brainstorming`.

Do not use for unrelated coding tasks.

## Platform Support

- Supported runtimes: `linux/amd64`, `linux/arm64`
- Not supported: Windows, macOS
- Before any execution, check current OS/arch.
- If unsupported, stop and print exactly:
  - `当前仅支持 Linux amd64/arm64，当前还不支持 Windows 和 macOS。`

## Packaging Contract

Skill folder structure must include:

```text
telegram-brainstorming/
  SKILL.md
  SKILL.zh-CN.md
  bin/
    telegram-brainstorming-linux-amd64
    telegram-brainstorming-linux-arm64
    .env
    .env.example
    .env.examples
```

Binary selection rule:
- `linux/amd64` -> `bin/telegram-brainstorming-linux-amd64`
- `linux/arm64` -> `bin/telegram-brainstorming-linux-arm64`

## Runtime Rules

1. Detect OS/arch first.
2. Resolve binary path from `bin/` by architecture.
3. Validate `.env` existence in `bin/`.
4. Validate `.env.examples` existence in `bin/` (if only `.env.example` exists, mirror it as `.env.examples`).
5. If `.env` is missing, print a clear instruction to create it from `.env.example` and stop.
6. Execute binary in text-only Telegram mode.
7. All brainstorming questions must be delivered in Telegram chat only.
8. Terminal output must contain status only (for example: running/waiting/completed), and must not print question bodies.

## Brainstorming Interaction Rules

- Ask exactly one question per message.
- Prefer 2-3 numbered options (`1/2/3`) before open-ended prompts.
- Accept either option number or short text.
- Continue until purpose, constraints, and success criteria are explicit.
- Never render brainstorming question text in terminal.
- Use concise prompt ending:
  - `请回复 1/2/3，或直接输入你的说明。`

## Network and Proxy

- Load Telegram and proxy settings from `.env`.
- Support `HTTPS_PROXY`, `HTTP_PROXY`, `ALL_PROXY`, `NO_PROXY`.
- Allow explicit proxy URL override when provided.
- Use bounded timeout/retry and return actionable errors.

## Reliability Baseline

- De-duplicate updates by `update_id`.
- Keep session transitions idempotent.
- Allow restart recovery from persisted state.
- Return explicit failure reason on network/proxy issues.

## Completion Criteria

This skill is production-ready only when:
- The correct Linux binary is selected and runs from `bin/`.
- Unsupported platforms are blocked with the required message.
- Missing `.env` produces an actionable setup error.
- A full Telegram text conversation can complete one brainstorming cycle.
