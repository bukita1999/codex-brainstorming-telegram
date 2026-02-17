# codex-brainstorming-telegram

English version: [README.md](README.md)

这个项目主要是为了解决一个实际问题：`brainstorming` skill 很好，但我不想一直坐在电脑前逐个回答问题。既然很多场景只需要在几个选项中选择一个，就可以直接在手机上完成。因此这里采用一种更简单的方法，让 coding agent 把交互内容发送到 Telegram，并等待用户在 Telegram 中回复，避免使用需要额外配置和部署时间的重型方案。

### 快速开始

1. 克隆仓库：

```bash
git clone https://github.com/bukita1999/codex-brainstorming-telegram.git
cd codex-brainstorming-telegram
```

2. 安装 Go（建议 `1.25+`）并确认可用：

```bash
go version
```

3. 基于 `.env.example` 创建 `.env`，并填写 Telegram 参数：

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

4. 先跑一次连通性测试（手动在 Telegram 回复相同六位数字）：

```bash
scripts/run_telegram_echo_test.sh
```

5. 用你的 AI Agent 安装 skill（如 `Codex` / `Claude Code` / `Opencode`）：

- 让 AI 读取 `instruction_for_AI.md`
- 示例提示词：`请你参照 instruction_for_AI.md 这个文档来尝试安装`

### 详细文档

- `docs/telegram-brainstorming-reference.md`：程序架构、运行逻辑、代码组成与开发调试命令（英文）。

### 生成与许可声明

- 本项目整体代码由 `Codex: GPT-5.3-Codex` 生成。
- 本项目采用 `MIT` 许可证。

### 平台支持

- 仅支持：Linux `amd64` / `arm64`
- 暂不支持：Windows / macOS

### 致谢

感谢大家使用本项目！🎉🙏  
如果遇到任何问题，欢迎提交 issue，并尽量附上复现步骤与日志信息。🐛🧪🛠️🚀
