# Telegram Brainstorming Reference

This document keeps detailed technical information that is intentionally removed from `README.md` to keep the README clean.

## What Is Included

- `cmd/telegram-echo-test`: CLI entry for the challenge/echo integrity test.
- `cmd/telegram-brainstorming`: CLI entry for one prompt->one reply Telegram interaction.
- `cmd/virtual-codex`: local virtual Codex binary for non-network testing.
- `internal/config`: `.env` parser and runtime config validation.
- `internal/telegramapi`: Telegram Bot API client (`sendMessage`, `getUpdates`).
- `internal/telegramtest`: challenge code generation and echo test orchestration.
- `internal/telegrambrainstorm`: brainstorming prompt/reply orchestration.
- `skills/telegram-brainstorming/`: production skill docs (English + Chinese translation).
- `instruction_for_AI.md`: build/package/install/update instructions for AI agents.
- `scripts/run_telegram_echo_test.sh`: manual entry script for challenge test.
- `scripts/run_telegram_brainstorming.sh`: manual entry script for one-round Telegram brainstorming.

## Runtime Logic

### 1) Configuration and startup

Both CLIs load `.env` through `internal/config.LoadTelegramConfig`.

- Required:
  - `TELEGRAM_BOT_TOKEN`
  - `TELEGRAM_CHAT_ID`
- Optional:
  - `TELEGRAM_PROXY_URL`
  - `TELEGRAM_REPLY_TIMEOUT` (default `5m`)

If `.env` is missing, the program returns an actionable error telling the user to create it from `.env.example`.

### 2) Proxy handling

`buildHTTPClient` clones `http.DefaultTransport`.

- Default proxy behavior is `http.ProxyFromEnvironment`.
- If `TELEGRAM_PROXY_URL` is set, it overrides the proxy with that explicit URL.
- HTTP timeout is set to `35s` for Telegram API requests.

This supports direct internet environments and proxied environments with the same binary.

### 3) Telegram-only interaction channel (brainstorming)

The brainstorming CLI prints status lines to terminal, but the actual question content is sent only to Telegram.

- Terminal: running status only.
- Telegram: prompt/options/plan/confirmation content.

This ensures the decision channel is Telegram, not terminal.

### 4) Prompt input contract

`cmd/telegram-brainstorming` accepts prompt text in one of two ways:

- `--prompt "..."`
- positional text argument

Rules:
- exactly one input mode must be used (not both).
- prompt must be non-empty after trim.

Escape handling:
- `\n` -> newline
- `\r` -> carriage return
- `\t` -> tab
- `\\` -> backslash

This allows structured multiline prompts to be passed in one command argument.

### 5) Update offset and polling model

Before sending a new message, the runner first reads the latest Telegram update offset with `getUpdates(offset=0, timeout=0)`.

Then it:
- sends the new prompt/challenge message,
- polls `getUpdates` with rolling `offset`,
- ignores old updates,
- ignores messages from other chats,
- returns on first valid reply.

Polling timeout is dynamic:
- minimum `1s`
- maximum `20s`
- based on remaining session time

This avoids replaying stale responses and reduces unnecessary polling load.

### 6) Brainstorming session lifecycle

`telegrambrainstorm.RunPrompt` flow:

1. Validate `chatID`, `prompt`, and timeout.
2. Snapshot latest offset.
3. Send prompt to Telegram.
4. Poll updates until timeout.
5. Return:
   - `RawReply` (trimmed original text)
   - `NormalizedReply` (currently same trim behavior)

Timeout returns `ErrSessionTimeout`.

CLI behavior:
- `stderr`: session status and error messages
- `stdout`: final normalized reply text only

### 7) Echo integrity test lifecycle

`telegram-echo-test` flow:

1. Generate random 6-digit challenge code.
2. Build fixed challenge sentence.
3. Send challenge to Telegram chat.
4. Poll updates until a matching reply is found or timeout occurs.
5. Report success/failure.

Reply matching accepts equivalent forms with wrappers (quotes/brackets/whitespace) but still requires the same underlying code.

### 8) Exit code contract

- `0`: success
- `1`: runtime failure or timeout
- `2`: usage/config/argument error

This keeps scripting integration predictable.

## Common Commands (Dev/Debug)

```bash
# Run Telegram echo test
scripts/run_telegram_echo_test.sh

# Run single-round Telegram brainstorming (--prompt)
GOCACHE=/tmp/go-build go run ./cmd/telegram-brainstorming --env .env --prompt "Choose one: A) Conservative B) Balanced C) Aggressive. Reply with A/B/C or short notes."

# Run single-round Telegram brainstorming (positional prompt)
GOCACHE=/tmp/go-build go run ./cmd/telegram-brainstorming --env .env "Choose one: A) Conservative B) Balanced C) Aggressive. Reply with A/B/C or short notes."

# Run single-round Telegram brainstorming (\n is converted to real line breaks)
GOCACHE=/tmp/go-build go run ./cmd/telegram-brainstorming --env .env --prompt "Choose one:\nA) Conservative\nB) Balanced\nC) Aggressive\nReply with A/B/C."

# Run all tests
GOCACHE=/tmp/go-build go test ./...

# Build Linux binaries
mkdir -p build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/telegram-brainstorming-linux-amd64 ./cmd/telegram-brainstorming
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/telegram-brainstorming-linux-arm64 ./cmd/telegram-brainstorming
```
