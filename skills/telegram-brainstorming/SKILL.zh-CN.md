---
name: telegram-brainstorming
description: 当用户明确要求通过 Telegram 进行 brainstorming，并且需要通过打包后的 Linux 二进制运行时使用。
---

# Telegram 头脑风暴（中文对照）

## 目标

通过 Telegram 完成执行前的完整协作流程：单轮单问题、优先选项、渐进式需求澄清、方案提议与执行确认。

该生产版以二进制为入口：skill 直接调用自身 `bin/` 下的 Linux 可执行文件。

## 触发条件

仅当用户明确要求 `telegram-brainstorming` 时使用。

与该流程无关的编码任务不使用此 skill。

## 平台支持

- 支持：`linux/amd64`、`linux/arm64`
- 不支持：Windows、macOS
- 在执行前必须先检测当前 OS/架构。
- 若不支持，必须停止并输出：
  - `当前仅支持 Linux amd64/arm64，当前还不支持 Windows 和 macOS。`

## 打包约定

skill 目录必须包含：

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

二进制选择规则：
- `linux/amd64` -> `bin/telegram-brainstorming-linux-amd64`
- `linux/arm64` -> `bin/telegram-brainstorming-linux-arm64`

## 运行规则

1. 先检测 OS/架构。
2. 根据架构从 `bin/` 选择对应二进制。
3. 检查 `bin/.env` 是否存在。
4. 检查 `bin/.env.examples` 是否存在（若仅有 `.env.example`，需同步生成 `.env.examples` 别名）。
5. 若 `.env` 缺失，明确提示根据 `.env.example` 创建 `.env`，并停止。
6. 以 Telegram 文本模式执行二进制，每次调用只处理一条 prompt。
7. 所有面向用户的交互内容都必须只在 Telegram 对话中发送。
8. 包括选项提问、澄清问题、完整执行方案描述、执行前确认问题。
9. 终端只能显示运行状态（运行中/等待中/已完成），不得打印提问或方案正文。

## 二进制输入输出契约

- 输入：调用二进制时传入一整段 prompt 文本（可包含较长 A/B/C 方案）。
- prompt 正文只能发送到 Telegram，不能在终端回显。
- 输出：二进制把 Telegram 收到的回复文本写到 `stdout`。
- 状态信息只写到终端状态流，且不得包含 prompt 正文。
- 每轮一个问题对应一次二进制调用。

## 内嵌协作规则

- 由 skill 侧组织每一轮问题（每轮只问一个问题）。
- 优先选项式提问（A/B/C 或 `1/2/3`），并允许简短自由输入。
- 二进制等待一条回复并返回给调用方。
- 调用方依据回复决定下一轮问题。
- 直到目标、约束、成功标准全部明确才收敛。
- 在执行任何实现命令前，必须先把完整执行方案发到 Telegram 并询问是否继续。
- 仅当 Telegram 中收到明确同意后才能进入执行。
- 若同意不明确或未给出同意，继续在 Telegram 澄清，不能开始执行。

## 网络与代理

- 从 `.env` 加载 Telegram 与代理配置。
- 支持 `HTTPS_PROXY`、`HTTP_PROXY`、`ALL_PROXY`、`NO_PROXY`。
- 若有显式代理 URL，允许覆盖默认代理行为。
- 请求必须有超时/有限重试，失败要返回可执行错误信息。

## 可靠性基线

- 按 `update_id` 去重。
- 状态迁移幂等。
- 进程重启后可恢复会话。
- 网络/代理失败时要有明确原因。

## 生产完成标准

满足以下才算生产可用：
- 可从 `bin/` 正确选择并执行 Linux 二进制。
- 不支持平台被阻断并输出固定提示。
- `.env` 缺失时有明确初始化指引。
- 提问正文只出现在 Telegram，不出现在终端。
- 一次二进制调用可完成一轮问答并返回回复文本。
- 从方案到执行前确认的决策链全部在 Telegram 完成，终端不承担决策交互。
