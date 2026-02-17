---
name: telegram-brainstorming
description: Use when the user explicitly asks to run brainstorming through Telegram chat instead of terminal interaction.
---

# Telegram Brainstorming

## Scope

MVP only. Use text messages only (no inline keyboard buttons).

## Trigger

Use this skill only when the user explicitly says `telegram-brainstorming`.

## Workflow

1. Load bot and proxy settings from `.env`.
2. Ask one question per message.
3. Prefer 2-3 numbered options in plain text.
4. Accept reply as either option number or short text.
5. Continue until purpose, constraints, and success criteria are clear.

## Message Pattern

- Prompt format: concise question + numbered options.
- End with: `请回复 1/2/3，或直接输入你的说明。`

## Reliability

- Use timeout for each wait step.
- Keep session state recoverable after process restart.
- On proxy/network failure, return actionable error text.
