# codex-brainstorming-telegram

English version: [README.md](README.md)

这是一个基于 Go 的 Telegram 文本验证与 brainstorming 技能打包项目。

### 快速开始

1. 安装 Go（建议 `1.25+`）并确认可用：

```bash
go version
```

2. 基于 `.env.example` 创建 `.env`，并填写 Telegram 参数：

```bash
cp .env.example .env
```

必填项：
- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_CHAT_ID`

可选项：
- `TELEGRAM_PROXY_URL`：仅在需要代理时填写（例如中国大陆网络环境）；若不需要代理可留空或删除该行。
- `TELEGRAM_REPLY_TIMEOUT`：默认 `5m`。

如果 `.env` 不存在，程序会提示你根据 `.env.example` 创建。

3. 先跑一次连通性测试（手动在 Telegram 回复相同六位数字）：

```bash
scripts/run_telegram_echo_test.sh
```

4. 用你的 AI Agent 安装 skill（如 `Codex` / `Claude Code` / `Opencode`）：

- 让 AI 读取 `instruction_for_AI.md`
- 示例提示词：`请你参照 instruction_for_AI.md 这个文档来尝试安装`

### 生成与许可声明

- 本项目整体代码由 `Codex: GPT-5.3-Codex` 生成。
- 本项目采用 `MIT` 许可证。

### 主要内容

- `cmd/telegram-echo-test`：向 Telegram 发送挑战消息并校验回包是否一致。
- `cmd/telegram-brainstorming`：单轮 prompt->reply 桥接器（传入一段 A/B/C 文本到 Telegram，等待一条回复并输出结果）。
- `skills/telegram-brainstorming/`：生产版 skill 文档（英文）与中文对照。
- `instruction_for_AI.md`：指导 AI 构建、打包、安装完整 skill。

### 平台支持

- 仅支持：Linux `amd64` / `arm64`
- 暂不支持：Windows / macOS

### 常用命令（开发/调试）

```bash
# 运行 Telegram 回环测试
scripts/run_telegram_echo_test.sh

# 运行单轮 Telegram brainstorming（--prompt 方式）
GOCACHE=/tmp/go-build go run ./cmd/telegram-brainstorming --env .env --prompt "请选择方案：A) 稳健 B) 平衡 C) 激进。请回复 A/B/C 或补充说明。"

# 运行单轮 Telegram brainstorming（位置参数方式）
GOCACHE=/tmp/go-build go run ./cmd/telegram-brainstorming --env .env "请选择方案：A) 稳健 B) 平衡 C) 激进。请回复 A/B/C 或补充说明。"

# 运行单轮 Telegram brainstorming（\n 会被解释为真实换行）
GOCACHE=/tmp/go-build go run ./cmd/telegram-brainstorming --env .env --prompt "请选择方案：\nA) 稳健\nB) 平衡\nC) 激进\n请回复 A/B/C。"

# 全量测试
GOCACHE=/tmp/go-build go test ./...

# 构建 Linux 二进制
mkdir -p build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/telegram-brainstorming-linux-amd64 ./cmd/telegram-brainstorming
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/telegram-brainstorming-linux-arm64 ./cmd/telegram-brainstorming
```

### 致谢

感谢大家使用本项目！🎉🙏  
如果遇到任何问题，欢迎提交 issue，并尽量附上复现步骤与日志信息。🐛🧪🛠️🚀
