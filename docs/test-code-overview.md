# 测试代码说明（逐文件）

本文档说明当前仓库中“测试相关代码”分别在验证什么，方便你快速判断改动影响范围。

## 1. 自动化测试（`go test ./...`）

### `cmd/virtual-codex/main_test.go`
- 验证 `virtual-codex` CLI 主流程的行为。
- 主要覆盖：
  - 正常输入时，`run()` 会输出预期的虚拟 Codex 文本。
  - 未提供 `--input` 时，返回用法错误（退出码 `2`），并提示 `input is required`。
  - 设置极短 `--timeout` 时，能正确触发超时错误（`deadline exceeded`，退出码 `1`）。

### `cmd/telegram-echo-test/main_test.go`
- 验证 `telegram-echo-test` CLI 在 `.env` 缺失时的错误提示是否足够明确。
- 主要覆盖：
  - `--env` 指向不存在文件时，退出码为 `2`。
  - 错误信息包含 `.env.example` 以及“创建对应的 `.env` 文件”的引导文案。

### `cmd/telegram-brainstorming/main_test.go`
- 验证 `telegram-brainstorming` CLI 的运行模式是否符合“单轮 prompt->reply”要求。
- 主要覆盖：
  - `.env` 缺失时返回可操作错误（包含 `.env.example` 提示）。
  - 未传入 prompt 时返回参数错误（退出码 `2`）。
  - 正常运行时：状态输出不包含 prompt 正文，`stdout` 仅返回 Telegram 回复文本。

### `internal/config/dotenv_test.go`
- 验证配置加载逻辑 `LoadTelegramConfig()`。
- 主要覆盖：
  - 能从 `.env` 正确读取 `TELEGRAM_BOT_TOKEN`、`TELEGRAM_CHAT_ID`、`TELEGRAM_PROXY_URL`、`TELEGRAM_REPLY_TIMEOUT`。
  - `.env` 文件不存在时，返回可操作的错误信息（提示参考 `.env.example`）。
  - 缺少必填项（token/chat id）时会报错。
  - 未配置超时时间时使用默认值 `5m`。

### `internal/virtualcodex/engine_test.go`
- 验证虚拟 Codex 引擎 `Engine.Respond()` 的核心规则。
- 主要覆盖：
  - 命中已知 prompt（如 `brainstorm login flow`）时，返回固定引导回复。
  - 未命中规则时，走回退逻辑（echo 输入）。
  - 空输入时返回 `ErrEmptyInput`。
  - 在处理延迟大于上下文超时时，返回 `context.DeadlineExceeded`。

### `internal/telegramapi/client_test.go`
- 验证 Telegram API 客户端封装是否正确组装请求并解析响应。
- 通过自定义 `RoundTripper` 模拟 HTTP，不依赖真实网络。
- 主要覆盖：
  - `SendMessage()`：请求路径、`Content-Type`、表单参数（`chat_id`/`text`）和响应 `message_id` 解析。
  - `GetUpdates()`：查询参数（`offset`/`timeout`）以及返回 update 列表解析。

### `internal/telegramtest/challenge_test.go`
- 验证挑战码与文本匹配相关的纯逻辑函数。
- 主要覆盖：
  - `GenerateCode()` 生成 6 位数字字符串。
  - `BuildChallengeMessage()` 生成固定格式消息：`这是一个测试，请回复 "[123456]"`。
  - `IsMatchingReply()` 能识别带引号、带方括号或带空白的等价回复，并拒绝错误验证码。

### `internal/telegramtest/runner_test.go`
- 验证挑战流程编排函数 `RunChallenge()` 的端到端逻辑（使用 fake API）。
- 主要覆盖：
  - 成功路径：先发送挑战消息，再轮询 updates，收到匹配回复后成功结束。
  - 超时路径：在指定时限内未收到匹配回复时返回 `ErrChallengeTimeout`。

### `internal/telegrambrainstorm/runner_test.go`
- 验证 Telegram 单轮问答编排逻辑（使用 fake API）。
- 主要覆盖：
  - 成功路径：发送一条传入 prompt，收到第一条有效回复后立即返回。
  - 返回值包含 `RawReply` 与 `NormalizedReply`。
  - 超时路径：在时限内未收到有效回复时返回 `ErrSessionTimeout`。

## 2. 手工联调脚本

### `scripts/run_telegram_echo_test.sh`
- 这是“手工测试入口脚本”，不是单元测试。
- 作用：
  - 切换到仓库根目录。
  - 以 `go run ./cmd/telegram-echo-test --env .env` 方式启动 Telegram echo 测试程序。
  - 允许把额外参数透传给主程序，便于你本地联调。

### `scripts/run_telegram_brainstorming.sh`
- 这是“Telegram 单轮问答手工联调入口脚本”，不是单元测试。
- 作用：
  - 切换到仓库根目录。
  - 以 `go run ./cmd/telegram-brainstorming --env .env "$@"` 方式启动程序。
  - 每次调用传入一条 prompt 到 Telegram，等待一条回复后返回。

## 3. 当前测试覆盖的重点与边界

- 已覆盖重点：
  - 配置读取与默认值。
  - 虚拟 Codex 关键分支（正常/空输入/超时）。
  - Telegram API 请求拼装与响应解析。
  - 挑战码流程（成功与超时）。
- 尚未覆盖（可后续补充）：
  - 真实 Telegram 网络环境下的自动化集成测试（目前由你手工测试）。
  - 代理异常、网络抖动、Telegram API 非 200 响应的更细粒度场景。
