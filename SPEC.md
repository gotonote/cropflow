# AgentFlow - 多智能体协作编排平台

## 项目概述

- **项目名称**: AgentFlow
- **类型**: 多通讯平台智能体协作编排系统
- **核心功能**: 可视化流程编排 + 多渠道消息接入 + 智能体任务分发
- **目标用户**: 需要构建AI智能体产品的开发者/企业

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | React 18 + TypeScript + React Flow |
| 后端 | Go (Gin框架) |
| 数据库 | PostgreSQL + Redis |
| 部署 | Docker |

## 核心功能

### 1. 可视化流程编排 (Flow Editor)
- 拖拽式智能体节点创建
- 节点间连线定义调用关系
- 支持串行、并行、条件分支
- 支持文字描述调整结构

### 2. 多渠道接入 (Multi-Channel)
- 飞书
- Telegram
- Discord
- WhatsApp
- Signal
- Web API

### 3. 智能体管理 (Agent Management)
- 创建/编辑/删除智能体
- 配置大模型参数
- 绑定工具能力

### 4. 消息引擎 (Message Engine)
- 消息接收与转发
- 对话上下文管理
- 会话状态追踪

## API设计

### 渠道管理
- `POST /api/channels` - 添加渠道
- `GET /api/channels` - 渠道列表
- `DELETE /api/channels/:id` - 删除渠道

### 智能体管理
- `POST /api/agents` - 创建智能体
- `GET /api/agents` - 智能体列表
- `PUT /api/agents/:id` - 更新智能体
- `DELETE /api/agents/:id` - 删除智能体

### 流程编排
- `POST /api/flows` - 创建流程
- `GET /api/flows` - 流程列表
- `PUT /api/flows/:id` - 更新流程
- `POST /api/flows/:id/execute` - 执行流程

### 消息
- `POST /api/webhook/:channel` - 渠道webhook入口

## 数据模型

### Agent
```
id, name, description, model_provider, model_name, 
model_config(json), tools(json), created_at, updated_at
```

### Flow
```
id, name, nodes(json), edges(json), 
trigger_type, created_at, updated_at
```

### Channel
```
id, name, type(feishu/telegram/discord/whatsapp), 
config(json), enabled, created_at
```

### Conversation
```
id, channel_id, user_id, messages(json),
context(json), created_at, updated_at
```

## 目录结构

```
agent-flow/
├── cmd/server/          # 入口文件
├── internal/
│   ├── api/            # HTTP处理器
│   ├── channel/        # 渠道适配器
│   ├── workflow/       # 流程执行引擎
│   ├── agent/          # 智能体核心
│   └── store/          # 数据层
├── frontend/           # React前端
└── docker-compose.yml  # 部署配置
```

## 开发阶段

### Phase 1: 基础框架
- [x] 项目结构搭建
- [ ] Go后端基础服务
- [ ] React前端基础

### Phase 2: 核心功能
- [ ] 流程编排API
- [ ] React Flow画布
- [ ] 智能体CRUD

### Phase 3: 渠道集成
- [ ] 飞书适配器
- [ ] Telegram适配器
- [ ] Webhook对接

### Phase 4: 高级功能
- [ ] 流程执行引擎
- [ ] 消息上下文
- [ ] 工具系统
