# CorpFlow

<p align="center">
  <img src="docs/logo.jpg" width="200" alt="CorpFlow Logo">
</p>

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


## Architecture

![Architecture Diagram](docs/architecture.svg)

---

## Demo: Build Your AI Team

![Demo: Build AI Team](docs/demo.svg)

### Why CorpFlow

1. **🔀 Visual Flow Editor** - Unlike CLI tools, CorpFlow provides a visual drag-and-drop workflow builder that's easy to use

2. **🗳️ Multi-Model Voting** - Unique feature! Let multiple AI models discuss and vote on the best answer

3. **📱 Mobile-First** - Full mobile app support for iOS, Android, Windows, macOS

4. **💬 Multi-Channel** - Deploy to Feishu, WeChat, Telegram, Discord simultaneously

5. **🧠 Hierarchical Memory** - Supervisors can view subordinate work history and generate reports

6. **🔧 Built-in Tools** - Shell, Git, Code Review, Test Generation, Web Search, File operations

7. **⚡ Ready Templates** - 8+ pre-built workflows: Chat, Voting, Research, Customer Service, Code Review, Content Creator, Data Analyzer, News Summarizer

---

## Features

| Feature | Description |
|---------|-------------|
| 🤖 **AI Agents** | Create custom AI agents with different models |
| 🔀 **Flow Builder** | Visual workflow automation with drag-and-drop |
| 💬 **Multi-Channel** | Feishu, WeChat, Telegram, Discord |
| 🗳️ **Multi-Model Voting** | Multiple AI models discuss and vote |
| 📱 **Mobile App** | iOS, Android, macOS, Windows, iPadOS |
| 🔧 **Tool Marketplace** | Shell, Git, Code Review, Test Gen, Calculator |
| 📋 **Execution Logs** | Step-by-step execution tracking |
| 🧠 **Memory System** | Hierarchical agent relationships |

---

## Supported AI Models

| Model | Provider | Env Variable |
|-------|----------|--------------|
| GPT-4 / GPT-4 Turbo | OpenAI | `OPENAI_API_KEY` |
| Claude 3 Opus / Sonnet | Anthropic | `ANTHROPIC_API_KEY` |
| GLM-4 / GLM-4 Flash | Zhipu | `ZHIPU_API_KEY` |
| Kimi | Moonshot | `KIMI_API_KEY` |
| Qwen Turbo / Plus | Alibaba | `DASHSCOPE_API_KEY` |
| DeepSeek Chat / Coder | DeepSeek | `DEEPSEEK_API_KEY` |
| MiniMax | MiniMax | `MINIMAX_API_KEY` |

---

## Quick Start

### Requirements

| Component | Version | Description |
|-----------|---------|-------------|
| **Go** | 1.21+ | Backend runtime |
| **Node.js** | 18+ | Frontend build |
| **PostgreSQL** | 14+ | Database |
| **Redis** | 7+ | Cache & sessions |
| **Flutter** | 3.16+ | Mobile app (optional) |

### Go Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| gin | v1.9.1 | HTTP framework |
| gorm | v1.25.5 | ORM |
| gorm.io/driver/postgres | v1.5.4 | PostgreSQL driver |
| go-redis/v9 | v9.3.0 | Redis client |
| go-openai | v1.17.0 | OpenAI API |
| chromedp | v0.9.5 | Browser automation |

### Frontend Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| react | ^18.2.0 | UI framework |
| @xyflow/react | ^12.0.0 | Flow editor |
| axios | ^1.6.0 | HTTP client |
| vite | ^5.0.0 | Build tool |
| typescript | ^5.3.0 | Type safety |

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

### Frontend (React)

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Open http://localhost:3000
```

### One-Click Install (Recommended)

```bash
# Run the install script
curl -sSL https://raw.githubusercontent.com/gotonote/corpflow/main/scripts/install.sh | bash

# After installation, edit .env with your API keys
vim corpflow/.env

# Start backend (terminal 1)
cd corpflow && go run cmd/server/main.go

# Start frontend (terminal 2)
cd corpflow/frontend && npm run dev
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
```

---

## How to Use Features

### 1. Chat with AI

1. Go to **Chat** tab
2. Select an agent
3. Type message
4. Get AI response

### 2. Create Workflow

1. Go to **Flows** tab
2. Drag nodes from sidebar
3. Connect them
4. Configure properties
5. Save & Execute

### 3. Manage Agents

1. Go to **Agents** tab
2. Click **+ Add Agent**
3. Set name, model, prompt
4. Save

### 4. View Execution Logs

1. Go to **Logs** tab
2. See all workflow runs
3. Click to see details
4. View step-by-step execution

---

## Connect Mobile App

To connect mobile app to local server:

1. Ensure phone and computer are on the same WiFi
2. Get your computer's IP:
   - Windows: `ipconfig`
   - Mac/Linux: `ifconfig`
3. In mobile app Settings, enter: `http://YOUR_IP:8080`

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

## Troubleshooting

| Problem | Solution |
|---------|----------|
| Can't access localhost:3000 | Check if Docker is running: `docker ps` |
| API calls fail | Verify API Key is configured in Settings |
| Mobile can't connect | Check firewall / ensure same network |
| Flow won't execute | Check all nodes are connected properly |

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/agents` | Create agent |
| GET | `/api/agents` | List agents |
| POST | `/api/flows` | Create flow |
| POST | `/api/flows/:id/execute` | Execute flow |
| POST | `/api/tools/execute` | Execute tool |
| GET | `/api/logs` | Get execution logs |
| POST | `/webhook/feishu` | Feishu webhook |

---

## License

MIT License
