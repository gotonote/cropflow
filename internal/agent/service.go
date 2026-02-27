package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// Service 智能体服务
type Service struct {
	openaiClient *openai.Client
	anthropicKey string
	defaultModel string
}

// NewService 创建智能体服务
func NewService() *Service {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	var client *openai.Client
	if openaiKey != "" {
		client = openai.NewClient(openaiKey)
	}

	return &Service{
		openaiClient: client,
		anthropicKey: os.Getenv("ANTHROPIC_API_KEY"),
		defaultModel: "gpt-4",
	}
}

// Process 处理用户消息
func (s *Service) Process(ctx context.Context, input, userID string) (string, error) {
	// 简单实现：直接调用大模型
	return s.CallLLM(ctx, s.defaultModel, input)
}

// ProcessWithAgent 使用指定智能体处理
func (s *Service) ProcessWithAgent(ctx context.Context, agentID, input, userID string) (string, error) {
	// TODO: 从DB获取智能体配置
	// 构建系统prompt
	systemPrompt := "你是一个AI助手，请帮助用户解决问题。"

	return s.callModel(ctx, s.defaultModel, systemPrompt, input)
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
