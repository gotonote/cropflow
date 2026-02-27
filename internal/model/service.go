package model

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/sashabaranov/go-openai"
)

// ModelProvider 模型供应商
type ModelProvider string

const (
	ProviderOpenAI   ModelProvider = "openai"
	ProviderAnthropic ModelProvider = "anthropic"
	ProviderGLM      ModelProvider = "glm"
	ProviderMiniMax  ModelProvider = "minimax"
	ProviderKimi     ModelProvider = "kimi"
	ProviderQwen     ModelProvider = "qwen"
	ProviderDeepSeek ModelProvider = "deepseek"
	ProviderCustom   ModelProvider = "custom"
)

// ModelConfig 模型配置
type ModelConfig struct {
	Provider    ModelProvider `json:"provider"`
	ModelName   string       `json:"model_name"`
	APIKey      string       `json:"api_key"`
	BaseURL     string       `json:"base_url"`      // 自定义API端点
	MaxTokens   int          `json:"max_tokens"`    // 最大token数
	Temperature float64      `json:"temperature"`   // 温度参数
	TopP        float64      `json:"top_p"`        // top_p采样
}

// Message 消息
type Message struct {
	Role    string `json:"role"`    // system/user/assistant
	Content string `json:"content"`
}

// Request 请求
type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Response 响应
type Response struct {
	Model      string   `json:"model"`
	Content    string   `json:"content"`
	FinishReason string `json:"finish_reason"`
}

// Service 模型服务
type Service struct {
	mu           sync.RWMutex
	models       map[string]*ModelConfig
	defaultModel string
	clients      map[ModelProvider]Client
}

// Client 模型客户端接口
type Client interface {
	Chat(ctx context.Context, req Request) (*Response, error)
}

// NewService 创建模型服务
func NewService() *Service {
	s := &Service{
		models:  make(map[string]*ModelConfig),
		clients: make(map[ModelProvider]Client),
	}

	// 初始化默认模型配置
	s.initDefaultModels()

	// 初始化客户端
	s.initClients()

	return s
}

// initDefaultModels 初始化默认模型
func (s *Service) initDefaultModels() {
	// OpenAI
	s.models["gpt-4"] = &ModelConfig{
		Provider:    ProviderOpenAI,
		ModelName:   "gpt-4",
		APIKey:      os.Getenv("OPENAI_API_KEY"),
		Temperature: 0.7,
		MaxTokens:   4096,
	}
	s.models["gpt-3.5-turbo"] = &ModelConfig{
		Provider:    ProviderOpenAI,
		ModelName:   "gpt-3.5-turbo",
		APIKey:      os.Getenv("OPENAI_API_KEY"),
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	// Anthropic
	s.models["claude-3-opus"] = &ModelConfig{
		Provider:    ProviderAnthropic,
		ModelName:   "claude-3-opus-20240229",
		APIKey:      os.Getenv("ANTHROPIC_API_KEY"),
		Temperature: 0.7,
		MaxTokens:   4096,
	}
	s.models["claude-3-sonnet"] = &ModelConfig{
		Provider:    ProviderAnthropic,
		ModelName:   "claude-3-sonnet-20240229",
		APIKey:      os.Getenv("ANTHROPIC_API_KEY"),
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	// 国产模型 - GLM (智谱)
	s.models["glm-4"] = &ModelConfig{
		Provider:    ProviderGLM,
		ModelName:   "glm-4",
		APIKey:      os.Getenv("ZHIPU_API_KEY"),
		BaseURL:     "https://open.bigmodel.cn/api/paas/v4",
		Temperature: 0.7,
		MaxTokens:   4096,
	}
	s.models["glm-3-turbo"] = &ModelConfig{
		Provider:    ProviderGLM,
		ModelName:   "glm-3-turbo",
		APIKey:      os.Getenv("ZHIPU_API_KEY"),
		BaseURL:     "https://open.bigmodel.cn/api/paas/v4",
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	// MiniMax
	s.models["abab6.5s-chat"] = &ModelConfig{
		Provider:    ProviderMiniMax,
		ModelName:   "abab6.5s-chat",
		APIKey:      os.Getenv("MINIMAX_API_KEY"),
		BaseURL:     "https://api.minimax.chat/v1",
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	// Kimi (月之暗面)
	s.models["moonshot-v1-8k-chat"] = &ModelConfig{
		Provider:    ProviderKimi,
		ModelName:   "moonshot-v1-8k-chat",
		APIKey:      os.Getenv("KIMI_API_KEY"),
		BaseURL:     "https://api.moonshot.cn/v1",
		Temperature: 0.7,
		MaxTokens:   4096,
	}
	s.models["moonshot-v1-32k-chat"] = &ModelConfig{
		Provider:    ProviderKimi,
		ModelName:   "moonshot-v1-32k-chat",
		APIKey:      os.Getenv("KIMI_API_KEY"),
		BaseURL:     "https://api.moonshot.cn/v1",
		Temperature: 0.7,
		MaxTokens:   32768,
	}

	// Qwen (通义千问)
	s.models["qwen-turbo"] = &ModelConfig{
		Provider:    ProviderQwen,
		ModelName:   "qwen-turbo",
		APIKey:      os.Getenv("DASHSCOPE_API_KEY"),
		BaseURL:     "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Temperature: 0.7,
		MaxTokens:   8192,
	}
	s.models["qwen-plus"] = &ModelConfig{
		Provider:    ProviderQwen,
		ModelName:   "qwen-plus",
		APIKey:      os.Getenv("DASHSCOPE_API_KEY"),
		BaseURL:     "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Temperature: 0.7,
		MaxTokens:   32768,
	}
	s.models["qwen-max"] = &ModelConfig{
		Provider:    ProviderQwen,
		ModelName:   "qwen-max",
		APIKey:      os.Getenv("DASHSCOPE_API_KEY"),
		BaseURL:     "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Temperature: 0.7,
		MaxTokens:   8192,
	}

	// DeepSeek
	s.models["deepseek-chat"] = &ModelConfig{
		Provider:    ProviderDeepSeek,
		ModelName:   "deepseek-chat",
		APIKey:      os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL:     "https://api.deepseek.com/v1",
		Temperature: 0.7,
		MaxTokens:   4096,
	}
	s.models["deepseek-coder"] = &ModelConfig{
		Provider:    ProviderDeepSeek,
		ModelName:   "deepseek-coder",
		APIKey:      os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL:     "https://api.deepseek.com/v1",
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	s.defaultModel = "gpt-4"
}

// initClients 初始化客户端
func (s *Service) initClients() {
	// OpenAI客户端
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		s.clients[ProviderOpenAI] = NewOpenAIClient(apiKey, "")
	}

	// Anthropic客户端
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		s.clients[ProviderAnthropic] = NewAnthropicClient(apiKey)
	}

	// 国产模型客户端
	if apiKey := os.Getenv("ZHIPU_API_KEY"); apiKey != "" {
		s.clients[ProviderGLM] = NewGLMClient(apiKey)
	}
	if apiKey := os.Getenv("MINIMAX_API_KEY"); apiKey != "" {
		s.clients[ProviderMiniMax] = NewMiniMaxClient(apiKey)
	}
	if apiKey := os.Getenv("KIMI_API_KEY"); apiKey != "" {
		s.clients[ProviderKimi] = NewKimiClient(apiKey)
	}
	if apiKey := os.Getenv("DASHSCOPE_API_KEY"); apiKey != "" {
		s.clients[ProviderQwen] = NewQwenClient(apiKey)
	}
	if apiKey := os.Getenv("DEEPSEEK_API_KEY"); apiKey != "" {
		s.clients[ProviderDeepSeek] = NewDeepSeekClient(apiKey)
	}
}

// ListModels 列出所有模型
func (s *Service) ListModels() []*ModelConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	configs := make([]*ModelConfig, 0, len(s.models))
	for _, cfg := range s.models {
		configs = append(configs, cfg)
	}
	return configs
}

// GetModel 获取模型配置
func (s *Service) GetModel(modelName string) (*ModelConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cfg, ok := s.models[modelName]
	return cfg, ok
}

// SetModel 设置模型配置
func (s *Service) SetModel(modelName string, config *ModelConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.models[modelName] = config
}

// Chat 调用单个模型
func (s *Service) Chat(ctx context.Context, modelName string, messages []Message) (*Response, error) {
	cfg, ok := s.GetModel(modelName)
	if !ok {
		return nil, fmt.Errorf("model not found: %s", modelName)
	}

	client, ok := s.clients[cfg.Provider]
	if !ok {
		return nil, fmt.Errorf("client not available for provider: %s", cfg.Provider)
	}

	req := Request{
		Model:       cfg.ModelName,
		Messages:    messages,
		Temperature: cfg.Temperature,
		MaxTokens:   cfg.MaxTokens,
	}

	return client.Chat(ctx, req)
}

// ========== 多模型投票 ==========

// VoteRequest 投票请求
type VoteRequest struct {
	Models   []string `json:"models"`   // 参与投票的模型列表
	Messages []Message `json:"messages"` // 消息历史
	SystemPrompt string   `json:"system_prompt"` // 系统提示
}

// VoteResponse 投票响应
type VoteResponse struct {
	Responses map[string]string `json:"responses"` // 模型 -> 响应内容
	Winner    string            `json:"winner"`    // 最终获胜者
	Scores    map[string]int    `json:"scores"`    // 投票得分
}

// Vote 多模型投票
func (s *Service) Vote(ctx context.Context, req VoteRequest) (*VoteResponse, error) {
	if len(req.Models) == 0 {
		return nil, fmt.Errorf("no models specified")
	}

	// 并发调用所有模型
	responses := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, modelName := range req.Models {
		wg.AddOne()
		go func(model string) {
			defer wg.Done()
			
			msgs := make([]Message, 0, len(req.Messages)+1)
			if req.SystemPrompt != "" {
				msgs = append(msgs, Message{Role: "system", Content: req.SystemPrompt})
			}
			msgs = append(msgs, req.Messages...)

			resp, err := s.Chat(ctx, model, msgs)
			mu.Lock()
			if err != nil {
				responses[model] = fmt.Sprintf("Error: %v", err)
			} else {
				responses[model] = resp.Content
			}
			mu.Unlock()
		}(modelName)
	}

	wg.Wait()

	// 简单投票逻辑：选择最长的回复作为最终响应
	// 生产环境可以使用更复杂的投票算法
	winner := ""
	maxLen := 0
	scores := make(map[string]int)

	for model, resp := range responses {
		scores[model] = 1 // 每个模型都有一票
		if len(resp) > maxLen {
			maxLen = len(resp)
			winner = model
		}
	}

	return &VoteResponse{
		Responses: responses,
		Winner:    winner,
		Scores:    scores,
	}, nil
}

// ========== 各个模型的客户端实现 ==========

// OpenAI客户端
type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey, baseURL string) *OpenAIClient {
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	return &OpenAIClient{client: openai.NewClientWithConfig(cfg)}
}

func (c *OpenAIClient) Chat(ctx context.Context, req Request) (*Response, error) {
	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	})
	if err != nil {
		return nil, err
	}

	return &Response{
		Model:      resp.Model,
		Content:    resp.Choices[0].Message.Content,
		FinishReason: string(resp.Choices[0].FinishReason),
	}, nil
}

// Anthropic客户端 (简化实现)
type AnthropicClient struct {
	apiKey string
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{apiKey: apiKey}
}

func (c *AnthropicClient) Chat(ctx context.Context, req Request) (*Response, error) {
	// TODO: 实现Anthropic API调用
	return &Response{Model: req.Model, Content: "Anthropic API not implemented"}, nil
}

// GLM客户端 (智谱)
type GLMClient struct {
	apiKey string
}

func NewGLMClient(apiKey string) *GLMClient {
	return &GLMClient{apiKey: apiKey}
}

func (c *GLMClient) Chat(ctx context.Context, req Request) (*Response, error) {
	// TODO: 实现智谱API调用
	return &Response{Model: req.Model, Content: "GLM API not implemented"}, nil
}

// MiniMax客户端
type MiniMaxClient struct {
	apiKey string
}

func NewMiniMaxClient(apiKey string) *MiniMaxClient {
	return &MiniMaxClient{apiKey: apiKey}
}

func (c *MiniMaxClient) Chat(ctx context.Context, req Request) (*Response, error) {
	// TODO: 实现MiniMax API调用
	return &Response{Model: req.Model, Content: "MiniMax API not implemented"}, nil
}

// Kimi客户端
type KimiClient struct {
	apiKey string
}

func NewKimiClient(apiKey string) *KimiClient {
	return &KimiClient{apiKey: apiKey}
}

func (c *KimiClient) Chat(ctx context.Context, req Request) (*Response, error) {
	// TODO: 实现Kimi API调用
	return &Response{Model: req.Model, Content: "Kimi API not implemented"}, nil
}

// Qwen客户端
type QwenClient struct {
	apiKey string
}

func NewQwenClient(apiKey string) *QwenClient {
	return &QwenClient{apiKey: apiKey}
}

func (c *QwenClient) Chat(ctx context.Context, req Request) (*Response, error) {
	// TODO: 实现通义千问API调用
	return &Response{Model: req.Model, Content: "Qwen API not implemented"}, nil
}

// DeepSeek客户端
type DeepSeekClient struct {
	apiKey string
}

func NewDeepSeekClient(apiKey string) *DeepSeekClient {
	return &DeepSeekClient{apiKey: apiKey}
}

func (c *DeepSeekClient) Chat(ctx context.Context, req Request) (*Response, error) {
	// TODO: 实现DeepSeek API调用
	return &Response{Model: req.Model, Content: "DeepSeek API not implemented"}, nil
}

// ========== 工具函数 ==========

// ToOpenAIMessages 转换消息格式
func ToOpenAIMessages(messages []Message) []openai.ChatCompletionMessage {
	result := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		result[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	return result
}

// GetEnvWithDefault 获取环境变量默认值
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
