package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"corpflow/internal/memory"
)

// Service 智能体服务
type Service struct {
	openaiClient *openai.Client
	anthropicKey string
	defaultModel string
	memorySvc    *memory.Service
}

// NewService 创建智能体服务
func NewService(memSvc *memory.Service) *Service {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	var client *openai.Client
	if openaiKey != "" {
		client = openai.NewClient(openaiKey)
	}

	return &Service{
		openaiClient: client,
		anthropicKey: os.Getenv("ANTHROPIC_API_KEY"),
		defaultModel: "gpt-4",
		memorySvc:    memSvc,
	}
}

// Process 处理用户消息
func (s *Service) Process(ctx context.Context, input, userID string) (string, error) {
	// 简单实现：直接调用大模型
	return s.CallLLM(ctx, s.defaultModel, input)
}

// ProcessWithAgent 使用指定智能体处理 (带记忆)
func (s *Service) ProcessWithAgent(ctx context.Context, agentID, input, userID string) (string, error) {
	// 获取智能体配置
	// TODO: 从DB获取
	
	// 获取相关记忆
	contextInfo := s.GetContextForAgent(agentID, input)
	
	// 构建系统prompt
	systemPrompt := "你是一个AI助手，请帮助用户解决问题。"
	if contextInfo != "" {
		systemPrompt += "\n\n相关背景信息:\n" + contextInfo
	}

	response := s.callModel(ctx, s.defaultModel, systemPrompt, input)
	
	// 记录到记忆
	if s.memorySvc != nil {
		_ = s.memorySvc.AddActionMemory(agentID, nil, input, response)
	}
	
	return response, nil
}

// GetContextForAgent 获取智能体的上下文记忆
func (s *Service) GetContextForAgent(agentID uint, currentInput string) string {
	if s.memorySvc == nil {
		return ""
	}

	var context strings.Builder

	// 1. 获取自己的记忆
	memories, _ := s.memorySvc.GetMemories(agentID, 10)
	if len(memories) > 0 {
		context.WriteString("【我的历史行为】\n")
		for _, mem := range memories {
			context.WriteString(fmt.Sprintf("- %s\n", mem.Content))
		}
	}

	// 2. 获取下属的记忆 (如果我是上级)
	subMemories, _ := s.memorySvc.GetSubordinateMemories(agentID)
	if len(subMemories) > 0 {
		context.WriteString("\n【下属的行为】\n")
		for _, mem := range subMemories {
			context.WriteString(fmt.Sprintf("- %s\n", mem.Content))
		}
	}

	// 3. 获取相关知识
	knowledge, _ := s.memorySvc.SearchKnowledge(agentID, currentInput)
	if len(knowledge) > 0 {
		context.WriteString("\n【相关知识】\n")
		for _, k := range knowledge {
			context.WriteString(fmt.Sprintf("- %s: %s\n", k.Title, k.Content))
		}
	}

	return context.String()
}

// SetParent 设置上级
func (s *Service) SetParent(agentID, parentID uint) error {
	if s.memorySvc == nil {
		return nil
	}
	return s.memorySvc.SetRelationship(parentID, agentID, "manage")
}

// AddSubordinate 添加下属
func (s *Service) AddSubordinate(agentID, subordinateID uint) error {
	if s.memorySvc == nil {
		return nil
	}
	return s.memorySvc.SetRelationship(agentID, subordinateID, "manage")
}

// GetSubordinates 获取下属列表
func (s *Service) GetSubordinates(agentID uint) ([]uint, error) {
	if s.memorySvc == nil {
		return []uint{}, nil
	}
	return s.memorySvc.GetSubordinateIDs(agentID)
}

// GenerateReport 生成工作报告
func (s *Service) GenerateReport(agentID uint, period string) (*memory.Report, error) {
	if s.memorySvc == nil {
		return nil, fmt.Errorf("memory service not initialized")
	}
	return s.memorySvc.GenerateReport(agentID, period)
}

// RecordDecision 记录决策
func (s *Service) RecordDecision(agentID uint, decision, reason string) error {
	if s.memorySvc == nil {
		return nil
	}
	return s.memorySvc.AddDecisionMemory(agentID, nil, decision, reason)
}

// RecordResult 记录执行结果
func (s *Service) RecordResult(agentID uint, task string, success bool, details string) error {
	if s.memorySvc == nil {
		return nil
	}
	return s.memorySvc.AddResultMemory(agentID, nil, task, success, details)
}

// CallLLM 调用大模型
func (s *Service) CallLLM(ctx context.Context, model, prompt string) (string, error) {
	systemPrompt := "你是一个AI助手，请简洁地回答用户问题。"
	return s.callModel(ctx, model, systemPrompt, prompt)
}

// callModel 调用模型
func (s *Service) callModel(ctx context.Context, model, systemPrompt, userPrompt string) (string, error) {
	if s.openaiClient == nil {
		return "⚠️ 请配置 OPENAI_API_KEY 环境变量", nil
	}

	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: userPrompt},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from model")
	}

	return resp.Choices[0].Message.Content, nil
}

// ========== 智能体配置 ==========

// Config 智能体配置
type Config struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Model       string                 `json:"model"`
	Provider    string                 `json:"provider"` // openai/anthropic/custom
	SystemPrompt string                `json:"system_prompt"`
	Tools       []string               `json:"tools"`    // 工具列表
	Temperature float64                `json:"temperature"`
	MaxTokens   int                    `json:"max_tokens"`
}

// Tool 工具定义
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Message 消息
type Message struct {
	Role    string `json:"role"`    // system/user/assistant
	Content string `json:"content"`
}

// Conversation 会话
type Conversation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Messages  []Message `json:"messages"`
	Context   map[string]interface{} `json:"context"`
	AgentID   string    `json:"agent_id"`
	CreatedAt int64     `json:"created_at"`
}

// AddMessage 添加消息
func (c *Conversation) AddMessage(role, content string) {
	c.Messages = append(c.Messages, Message{
		Role:    role,
		Content: content,
	})
}

// ToOpenAIMessages 转换为OpenAI消息格式
func (c *Conversation) ToOpenAIMessages() []openai.ChatCompletionMessage {
	msgs := make([]openai.ChatCompletionMessage, len(c.Messages))
	for i, m := range c.Messages {
		msgs[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	return msgs
}

// ========== 工具服务 ==========

// ToolService 工具服务
type ToolService struct {
	tools map[string]ToolHandler
}

// ToolHandler 工具处理器
type ToolHandler func(ctx context.Context, params map[string]interface{}) (string, error)

// NewToolService 创建工具服务
func NewToolService() *ToolService {
	ts := &ToolService{
		tools: make(map[string]ToolHandler),
	}

	// 注册内置工具
	ts.Register("search_web", ts.handleSearchWeb)
	ts.Register("calculate", ts.handleCalculate)
	ts.Register("get_time", ts.handleGetTime)
	ts.Register("fetch_url", ts.handleFetchURL)

	return ts
}

// Register 注册工具
func (ts *ToolService) Register(name string, handler ToolHandler) {
	ts.tools[name] = handler
}

// Execute 执行工具
func (ts *ToolService) Execute(ctx context.Context, name string, params map[string]interface{}) (string, error) {
	handler, ok := ts.tools[name]
	if !ok {
		return "", fmt.Errorf("tool not found: %s", name)
	}
	return handler(ctx, params)
}

// ListTools 列出所有工具
func (ts *ToolService) ListTools() []Tool {
	tools := make([]Tool, 0, len(ts.tools))
	for name, handler := range ts.tools {
		// 通过反射或注册时存储description
		tools = append(tools, Tool{
			Name:        name,
			Description: fmt.Sprintf("Tool: %s", name),
		})
	}
	_ = handler // 避免编译错误
	return tools
}

// ========== 内置工具实现 ==========

func (ts *ToolService) handleSearchWeb(ctx context.Context, params map[string]interface{}) (string, error) {
	query, _ := params["query"].(string)
	if query == "" {
		return "", fmt.Errorf("missing query parameter")
	}
	// TODO: 集成搜索API
	return fmt.Sprintf("搜索结果: %s (请配置搜索API)", query), nil
}

func (ts *ToolService) handleCalculate(ctx context.Context, params map[string]interface{}) (string, error) {
	expr, _ := params["expression"].(string)
	if expr == "" {
		return "", fmt.Errorf("missing expression parameter")
	}
	// 简单计算器 (生产环境使用专门库)
	// TODO: 使用govaluate或其他库
	return fmt.Sprintf("计算结果: %s", expr), nil
}

func (ts *ToolService) handleGetTime(ctx context.Context, params map[string]interface{}) (string, error) {
	return "当前时间获取成功", nil
}

func (ts *ToolService) handleFetchURL(ctx context.Context, params map[string]interface{}) (string, error) {
	url, _ := params["url"].(string)
	if url == "" {
		return "", fmt.Errorf("missing url parameter")
	}
	// TODO: 实际获取URL内容
	return fmt.Sprintf("获取内容: %s", url), nil
}

// ========== 工具转换为OpenAI格式 ==========

// ToOpenAITools 转换为OpenAI工具格式
func ToOpenAITools(tools []Tool) []openai.Tool {
	result := make([]openai.Tool, len(tools))
	for i, t := range tools {
		params, _ := json.Marshal(t.Parameters)
		result[i] = openai.Tool{
			Type: "function",
			Function: &openai.FunctionDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  json.RawMessage(params),
			},
		}
	}
	return result
}
