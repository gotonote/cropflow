# CorpFlow

**Multi-Agent Collaboration Platform**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Flutter-blue.svg)](https://flutter.dev)

> **中文**: [README_zh.md](./README_zh.md)

---

## Overview

CorpFlow is a **multi-agent collaboration platform** that enables you to:
- Create and manage AI agents
- Build visual workflows with drag-and-drop
- Deploy across multiple channels (Feishu, WeChat, Telegram, Discord)
- Use multiple AI models with intelligent voting

---

## Features

| Feature | Description |
|---------|-------------|
| 🤖 **AI Agents** | Create custom AI agents with different models |
| 🔀 **Flow Builder** | Visual workflow automation |
| 💬 **Multi-Channel** | Feishu, WeChat, Telegram, Discord |
| 🗳️ **Multi-Model Voting** | Multiple AI models discuss and vote |
| 📱 **Mobile App** | iOS, Android, macOS, Windows, iPadOS |

---

## Supported AI Models

| Model | Provider | Env Variable |
|-------|----------|--------------|
| GPT-4 | OpenAI | `OPENAI_API_KEY` |
| Claude 3 | Anthropic | `ANTHROPIC_API_KEY` |
| GLM-4 | Zhipu | `ZHIPU_API_KEY` |
| Kimi | Moonshot | `KIMI_API_KEY` |
| Qwen | Alibaba | `DASHSCOPE_API_KEY` |
| DeepSeek | DeepSeek | `DEEPSEEK_API_KEY` |
| MiniMax | MiniMax | `MINIMAX_API_KEY` |

---

## Quick Start

### Backend (Go + Docker)

```bash
# Clone the repo
git clone https://github.com/gotonote/corpflow.git
cd corpflow

# Copy configuration
cp .env.example .env

# Edit .env with your API keys
vim .env

# Start with Docker
docker-compose up -d
```

### Mobile App (Flutter)

```bash
cd mobile

# Install dependencies
flutter pub get

# Run in development
flutter run

# Build for Android
flutter build apk --release

# Build for iOS (macOS only)
flutter build ios --release
```

---

## How to Use

### 💬 Chat

1. Tap **"New Chat"** button
2. Type your message in the input field
3. AI responds instantly
4. Conversation is saved automatically

### 🔀 Flow

1. Go to **Flow** tab
2. Tap **"+"** to create new flow
3. Add nodes (Trigger/Agent/Tool/Condition)
4. Connect nodes by dragging
5. Save your flow
6. Execute by tapping play button

### 🤖 Agents

1. Go to **Agents** tab
2. Tap **"+"** to create new agent
3. Configure: name, model, system prompt, tools
4. Save and use in flows or chat

### 🗳️ Multi-Model Voting

Enable in **Settings** → Multi-Model Voting

**Voting Methods:**
- **Comprehensive**: Scores by Accuracy + Completeness + Clarity + Creativity
- **Cross-evaluation**: Models evaluate each other
- **Length**: By response length

**Scoring weights:**
- Accuracy - 30%
- Completeness - 30%
- Clarity - 20%
- Creativity - 20%

---

## Environment Variables

```bash
# AI Models
export OPENAI_API_KEY=sk-xxx
export ANTHROPIC_API_KEY=sk-ant-xxx
export ZHIPU_API_KEY=xxx
export KIMI_API_KEY=xxx
export DASHSCOPE_API_KEY=xxx
export DEEPSEEK_API_KEY=xxx
export MINIMAX_API_KEY=xxx

# Channels
export FEISHU_APP_ID=xxx
export FEISHU_APP_SECRET=xxx
export WECHAT_APP_ID=xxx
export TELEGRAM_BOT_TOKEN=xxx
```

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/agents` | Create agent |
| GET | `/api/agents` | List agents |
| POST | `/api/flows` | Create flow |
| POST | `/api/flows/:id/execute` | Execute flow |
| POST | `/webhook/feishu` | Feishu webhook |

---

## License

MIT License
