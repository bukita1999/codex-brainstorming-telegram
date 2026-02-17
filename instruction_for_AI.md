# Instruction For AI: Build and Install a Production Skill Package

This instruction defines how to create a new production skill package from this Go project.

## Objective

Create a fully packaged skill that:
- Uses Linux binaries built from this repository
- Includes runtime environment files in `bin/`
- Can be installed or reinstalled into the current coding agent skill directory
- Fully overwrites any existing skill with the same name (including `bin/` and env files)

## Required Output Structure

Create a new skill folder with this exact layout:

```text
<skill-name>/
  SKILL.md
  SKILL.zh-CN.md
  bin/
    telegram-brainstorming-linux-amd64
    telegram-brainstorming-linux-arm64
    .env
    .env.example
    .env.examples
```

## Steps

1. Build Linux binaries from this Go project.
2. Create the target skill folder and `bin/` subfolder.
3. Copy skill docs (`SKILL.md` and `SKILL.zh-CN.md`) into the skill root.
4. Copy binaries into `bin/` and keep executable permissions.
5. Copy `.env`, `.env.example`, and `.env.examples` into `bin/` (ensure all required variants exist).
6. If the repository uses `.env.example` naming, copy it and also provide `.env.examples` as an alias.
7. Pack the skill as an archive (for delivery and installation).
8. If an old skill already exists at install target, remove it first, then unpack the new one.
9. Verify installed files exist and executable permissions are correct.

## Build Commands (Reference)

```bash
# from repository root
mkdir -p build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/telegram-brainstorming-linux-amd64 ./cmd/telegram-brainstorming
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/telegram-brainstorming-linux-arm64 ./cmd/telegram-brainstorming
```

## Packaging Commands (Reference)

```bash
SKILL_NAME="telegram-brainstorming"
PKG_DIR="dist/${SKILL_NAME}"
mkdir -p "${PKG_DIR}/bin"

cp skills/telegram-brainstorming/SKILL.md "${PKG_DIR}/SKILL.md"
cp skills/telegram-brainstorming/SKILL.zh-CN.md "${PKG_DIR}/SKILL.zh-CN.md"
cp build/telegram-brainstorming-linux-amd64 "${PKG_DIR}/bin/"
cp build/telegram-brainstorming-linux-arm64 "${PKG_DIR}/bin/"
cp .env "${PKG_DIR}/bin/.env"

if [ -f .env.example ]; then
  cp .env.example "${PKG_DIR}/bin/.env.example"
else
  echo ".env.example is required"
  exit 1
fi

if [ -f .env.examples ]; then
  cp .env.examples "${PKG_DIR}/bin/.env.examples"
else
  cp .env.example "${PKG_DIR}/bin/.env.examples"
fi

chmod +x "${PKG_DIR}/bin/telegram-brainstorming-linux-amd64"
chmod +x "${PKG_DIR}/bin/telegram-brainstorming-linux-arm64"

tar -C dist -czf "dist/${SKILL_NAME}.tar.gz" "${SKILL_NAME}"
```

## Install Into Current Coding Agent

Install or update by extracting to the active skills directory used by your coding agent.

Example target (adapt to environment):
- `~/.codex/skills/<skill-name>/`

Reference:

```bash
mkdir -p ~/.codex/skills
rm -rf ~/.codex/skills/telegram-brainstorming
tar -C ~/.codex/skills -xzf dist/telegram-brainstorming.tar.gz
```

The remove-and-extract sequence is mandatory for update/reinstall. It guarantees complete overwrite of:
- `SKILL.md` and `SKILL.zh-CN.md`
- all binaries under `bin/`
- `bin/.env`, `bin/.env.example`

## Validation Checklist

- `SKILL.md` exists in installed skill root
- `SKILL.zh-CN.md` exists in installed skill root
- both Linux binaries exist in `bin/`
- `bin/.env`, `bin/.env.example` exist
- binaries are executable
- unsupported platform message is documented in `SKILL.md`

## Binary Runtime Contract

Installed `telegram-brainstorming` binaries must follow this contract:
- Input: one prompt string per invocation (`--prompt "..."` or positional text).
- Behavior: send prompt to Telegram chat and wait for one reply.
- Output: print only the received reply text to `stdout`.
- Logging: terminal status logs must not include prompt body.

## Important Runtime Note

If `.env` is missing at runtime, the program must print an actionable message telling the user to create `.env` from `.env.example`.
