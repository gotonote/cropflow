# AgentFlow - 多智能体协作编排平台

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react" alt="React">
  <img src="https://img.shields.io/badge/TypeScript-5.3+-3178C6?style=flat&logo=typescript" alt="TypeScript">
  <img src="https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker" alt="Docker">
</p>

## 概述

AgentFlow 是一个支持**多通讯平台**的智能体协作编排系统。核心特性：

- 🎨 **可视化流程编排** - 拖拽式工作流设计，支持流程图连线
- 📱 **多渠道接入** - 飞书、Telegram、Discord、WhatsApp、Webhook
- 🤖 **多智能体协作** - 支持串行、并行、条件分支
- 🔧 **工具系统** - 浏览器控制、文件读写、搜索等

## 快速开始

### 1. 环境要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+

### 2. 本地开发

```bash
# 克隆项目
git clone https://github.com/your-repo/agent-flow.git
cd agent-flow

# 复制配置
cp .env.example .env

# 编辑 .env 配置必要参数
# DATABASE_URL=postgres://user:pass@localhost:5432/agentflow
# REDIS_URL=localhost:6379
# FEISHU_APP_ID=你的飞书AppID
# FEISHU_APP_SECRET=你的飞书AppSecret
# OPENAI_API_KEY=sk-xxx

# 启动后端
go mod tidy
go run cmd/server/main.go

# 启动前端 (新终端)
cd frontend
npm install
npm run dev
```

### 3. Docker 部署

```bash
# 一键启动
docker-compose up -d

# 查看日志
docker-compose logs -f
```

访问 http://localhost:3000 打开管理界面。

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        React Flow 画布                       │
│    (可视化编排智能体协作关系，支持拖拽节点、连线、条件分支)      │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                     Go API Gateway                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  渠道适配器  │  │  流程引擎   │  │   智能体服务        │  │
│  │  飞书/TG    │  │  节点执行   │  │   大模型调用        │  │
│  │  Discord   │  │  上下文管理  │  │   工具系统          │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                   PostgreSQL + Redis                         │
│       (流程配置/智能体配置/会话历史/实时缓存)                  │
└─────────────────────────────────────────────────────────────┘
```

## 核心模块

### 1. 渠道层 (internal/channel)

| 文件 | 功能 |
|------|------|
| `manager.go` | 渠道管理器，接收消息、路由处理 |
| `feishu.go` | 飞书适配器 |
| `telegram.go` | Telegram适配器 (待完善) |

### 2. 流程引擎 (internal/workflow)

| 文件 | 功能 |
|------|------|
| `engine.go` | 流程执行引擎，节点调度、上下文管理 |

### 3. 智能体服务 (internal/agent)

| 文件 | 功能 |
|------|------|
| `service.go` | 智能体核心，大模型调用、工具执行 |

### 4. API层 (internal/api)

| 文件 | 功能 |
|------|------|
| `handler.go` | REST API处理器 |

## API文档

### 创建智能体
```bash
POST /api/agents
{
  "name": "客服助手",
  "description": "处理客户咨询",
  "model_provider": "openai",
  "model_name": "gpt-4",
  "tools": ["search", "calculate"]
}
```

### 创建流程
```bash
POST /api/flows
{
  "name": "用户咨询流程",
  "nodes": [...],  // React Flow nodes
  "edges": [...]   // React Flow edges
}
```

### 执行流程
```bash
POST /api/flows/:id/execute
{
  "input": "用户问题",
  "user_id": "user_123"
}
```

### 飞书Webhook
```
POST /webhook/feishu
```

## 配置说明

### 环境变量

| 变量 | 必填 | 说明 |
|------|------|------|
| `DATABASE_URL` | ✅ | PostgreSQL连接字符串 |
| `REDIS_URL` | ✅ | Redis地址 |
| `PORT` | - | 服务端口，默认8080 |
| `OPENAI_API_KEY` | - | OpenAI API Key |
| `ANTHROPIC_API_KEY` | - | Anthropic API Key |
| `FEISHU_APP_ID` | - | 飞书应用ID |
| `FEISHU_APP_SECRET` | - | 飞书应用密钥 |

## 飞书配置指南

1. 登录 [飞书开放平台](https://open.feishu.cn/)
2. 创建企业应用
3. 开启"机器人"能力
4. 配置"事件订阅"：
   - `im.message.receive_v1` - 接收消息
   - `im.chat.message.v1` - 群消息
5. 获取 `App ID` 和 `App Secret`
6. 设置"重定向URL"为你的服务器地址

## 界面预览

### 流程编排
- 左侧：节点库（智能体、条件分支、工具）
- 中间：React Flow 画布，支持拖拽连线
- 右侧：节点属性配置

### 智能体管理
- 创建/编辑/删除智能体
- 配置模型、工具、System Prompt

### 渠道配置
- 飞书/Telegram/Discord/WhatsApp
- 渠道启用/禁用
- 密钥配置

## 路线图

- [ ] 流程执行引擎完善
- [ ] 大模型多供应商支持（OpenAI/Claude/通义/文心）
- [ ] 工具系统完善（浏览器控制、代码执行）
- [ ] Telegram/Discord适配器完善
- [ ] Web管理界面增强
- [ ] 用户认证系统
- [ ] 多租户支持

## 许可证

MIT License
