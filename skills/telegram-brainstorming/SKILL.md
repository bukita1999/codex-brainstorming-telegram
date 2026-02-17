---
name: telegram-brainstorming
description: Use when the user explicitly requests brainstorming through Telegram chat and the workflow must run via packaged Linux binaries instead of terminal-only interaction.
---

# Telegram Brainstorming

## Purpose

Run the full pre-execution collaboration workflow through Telegram: one question per turn, option-first prompts, progressive requirement clarification, plan proposal, and execution confirmation.

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
6. Execute binary in text-only Telegram mode with one prompt argument per invocation.
7. All user-facing interaction content must be delivered in Telegram chat only.
8. This includes option prompts, clarification questions, full plan descriptions, and pre-execution confirmation requests.
9. Terminal output must contain status only (for example: running/waiting/completed), and must not print question or plan bodies.

## Binary I/O Contract

- Input: pass one full prompt string to the binary (may include long A/B/C options).
- Prompt content must be sent to Telegram only, never echoed in terminal.
- Output: binary returns the received Telegram reply as `stdout` text.
- Status logs must go to terminal status stream only, without prompt body.
- Use one binary invocation per question round.

## Embedded Collaboration Rules

- Skill side builds each question prompt (single question per round).
- Prefer options-first prompts (A/B/C or `1/2/3`) plus short free-text fallback.
- Binary waits for one reply and returns it to the caller.
- Caller decides next question based on returned reply.
- Continue rounds until purpose, constraints, and success criteria are explicit.
- Before any implementation command, send a full execution plan to Telegram and ask whether to proceed.
- Only execute when explicit approval is received in Telegram.
- If approval is missing or unclear, continue Telegram clarification and do not start execution.

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
- Prompt text appears only in Telegram, not terminal.
- One binary invocation can complete one ask/reply round and return reply text.
- Plan-to-execution confirmation is completed in Telegram, and terminal never becomes the decision channel.
