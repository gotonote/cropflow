# CorpFlow

**Multi-Agent Collaboration Platform**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Flutter-blue.svg)](https://flutter.dev)
[![AI Models](https://img.shields.io/badge/AI-Models-GPT--4%20%7C%20Claude%20%7C%20GLM%20%7C%20Kimi-green.svg)](https://github.com/gotonote/corpflow)

---

## Overview | 概述

CorpFlow is a **multi-agent collaboration platform** that enables you to:
- Create and manage AI agents
- Build visual workflows with drag-and-drop
- Deploy across multiple channels (Feishu, WeChat, Telegram, Discord)
- Use multiple AI models with intelligent voting

CorpFlow 是一个**多智能体协作平台**，支持：
- 创建和管理 AI 智能体
- 可视化流程编排（拖拽操作）
- 多渠道部署（飞书、微信、Telegram、Discord）
- 多模型投票决策

---

## Features | 功能

| Feature | Description | 功能 | 说明 |
|---------|-------------|------|------|
| 🤖 **AI Agents** | Create custom AI agents with different models | 智能体 | 创建自定义AI智能体 |
| 🔀 **Flow Builder** | Visual workflow automation | 流程编排 | 可视化工作流自动化 |
| 💬 **Multi-Channel** | Feishu, WeChat, Telegram, Discord | 多渠道 | 飞书、微信、Telegram、Discord |
| 🗳️ **Multi-Model Voting** | Multiple AI models discuss and vote | 多模型投票 | 多AI讨论并投票 |
| 📱 **Mobile App** | iOS, Android, macOS, Windows, iPadOS | 移动应用 | 全平台支持 |

---

## Supported AI Models | 支持的AI模型

| Model | Provider | 中文名 | Env Variable |
|-------|----------|--------|--------------|
| GPT-4 | OpenAI | - | `OPENAI_API_KEY` |
| Claude 3 | Anthropic | - | `ANTHROPIC_API_KEY` |
| GLM-4 | Zhipu | 智谱GLM | `ZHIPU_API_KEY` |
| Kimi | Moonshot | 月之暗面 | `KIMI_API_KEY` |
| Qwen | Alibaba | 通义千问 | `DASHSCOPE_API_KEY` |
| DeepSeek | DeepSeek | 深度求索 | `DEEPSEEK_API_KEY` |
| MiniMax | MiniMax | - | `MINIMAX_API_KEY` |

---

## Quick Start | 快速开始

### Backend (Go + Docker)

```bash
# Clone the repo
git clone https://github.com/gotonote/corpflow.git
cd corpflow

# Copy configuration
cp .env.example .env

# Edit .env with your API keys
# 添加你的 API Key 到 .env 文件

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

## How to Use | 使用指南

### 💬 Chat | 对话

**English:**
1. Tap **"New Chat"** button
2. Type your message in the input field
3. AI responds instantly
4. Conversation is saved automatically

**中文:**
1. 点击 **"新建对话"** 按钮
2. 在输入框输入消息
3. AI 即时回复
4. 对话自动保存

---

### 🔀 Flow | 流程编排

**English:**
1. Go to **Flow** tab
2. Tap **"+"** to create new flow
3. **Add nodes**:
   - **Trigger**: Message trigger, schedule, webhook
   - **Agent**: AI agent node
   - **Tool**: Browser, search, calculator
   - **Condition**: Branch logic
4. **Connect nodes** by dragging from output to input
5. **Save** your flow
6. **Execute** by tapping play button

**中文:**
1. 进入 **流程** 标签
2. 点击 **"+"** 创建新流程
3. **添加节点**：
   - **触发器**：消息触发、定时任务、Webhook
   - **智能体**：AI 节点
   - **工具**：浏览器、搜索、计算器
   - **条件**：分支逻辑
4. **连接节点**：从输出拖拽到输入
5. **保存**流程
6. 点击播放按钮**执行**

---

### 🤖 Agents | 智能体

**English:**
1. Go to **Agents** tab
2. Tap **"+"** to create new agent
3. Configure:
   - Name your agent
   - Select AI model
   - Set system prompt
   - Enable tools
4. Save and use in flows or chat

**中文:**
1. 进入 **智能体** 标签
2. 点击 **"+"** 创建新智能体
3. 配置：名称、模型、系统提示词、工具
4. 保存后在流程或对话中使用

---

### 🗳️ Multi-Model Voting | 多模型投票

**English:**
Enable in **Settings** → Multi-Model Voting

1. Enable voting toggle
2. Select voting method:
   - **Comprehensive**: Scores by Accuracy + Completeness + Clarity + Creativity
   - **Cross-evaluation**: Models evaluate each other
   - **Length**: Simple by response length
3. Multiple AI models will respond
4. System automatically selects the best response

**Scoring weights:**
- Accuracy - 30%
- Completeness - 30%
- Clarity - 20%
- Creativity - 20%

**中文:**
在 **设置** → 多模型投票 中启用

1. 开启投票开关
2. 选择投票方式：综合评分/交叉评估/按长度
3. 多个 AI 模型同时响应
4. 系统自动选择最佳答案

**评分权重：**
- 准确性 - 30%
- 完整性 - 30%
- 清晰度 - 20%
- 创造性 - 20%

---

## Environment Variables | 环境变量

```bash
# AI Models
export OPENAI_API_KEY=sk-xxx        # OpenAI
export ANTHROPIC_API_KEY=sk-ant-xxx # Anthropic
export ZHIPU_API_KEY=xxx            # 智谱GLM
export KIMI_API_KEY=xxx             # Kimi
export DASHSCOPE_API_KEY=xxx         # 阿里通义千问
export DEEPSEEK_API_KEY=xxx         # DeepSeek
export MINIMAX_API_KEY=xxx           # MiniMax

# Channels | 渠道
export FEISHU_APP_ID=xxx            # 飞书
export FEISHU_APP_SECRET=xxx
export WECHAT_APP_ID=xxx            # 微信
export WECHAT_APP_SECRET=xxx
export TELEGRAM_BOT_TOKEN=xxx       # Telegram
```

---

## API Endpoints | API接口

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
