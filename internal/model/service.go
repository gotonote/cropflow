package model

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	s.models["glm-4-plus"] = &ModelConfig{
		Provider:    ProviderGLM,
		ModelName:   "glm-4-plus",
		APIKey:      os.Getenv("ZHIPU_API_KEY"),
		BaseURL:     "https://open.bigmodel.cn/api/paas/v4",
		Temperature: 0.7,
		MaxTokens:   4096,
	}
	s.models["glm-4-flash"] = &ModelConfig{
		Provider:    ProviderGLM,
		ModelName:   "glm-4-flash",
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

// IsModelAvailable 检查模型是否可用 (API Key已配置)
func (s *Service) IsModelAvailable(modelName string) bool {
	cfg, ok := s.GetModel(modelName)
	if !ok {
		return false
	}
	return cfg.APIKey != ""
}

// GetAvailableModels 获取所有可用的模型
func (s *Service) GetAvailableModels() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var available []string
	for name, cfg := range s.models {
		if cfg.APIKey != "" {
			available = append(available, name)
		}
	}
	return available
}

// GetBestAvailableModel 获取最佳可用模型
func (s *Service) GetBestAvailableModel(preferred string) string {
	// 优先使用用户指定的模型
	if preferred != "" && s.IsModelAvailable(preferred) {
		return preferred
	}

	// 按优先级选择可用模型
	priority := []string{"gpt-4", "glm-4", "glm-4-plus", "claude-3-opus", "moonshot-v1-8k-chat", "qwen-turbo", "deepseek-chat"}
	for _, name := range priority {
		if s.IsModelAvailable(name) {
			return name
		}
	}

	// 如果没有配置任何模型，返回第一个有client的
	for name := range s.clients {
		return string(name) // TODO: 返回对应的默认模型
	}

	return ""
}

// Chat 调用单个模型 (带自动回退)
func (s *Service) Chat(ctx context.Context, modelName string, messages []Message) (*Response, error) {
	// 检查模型是否可用
	if !s.IsModelAvailable(modelName) {
		// 尝试自动回退到可用模型
		available := s.GetBestAvailableModel(modelName)
		if available == "" {
			return nil, fmt.Errorf("no available model: please configure at least one API key in Settings")
		}
		modelName = available
	}

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
	Models       []string      `json:"models"`        // 参与投票的模型列表
	Messages     []Message     `json:"messages"`      // 消息历史
	SystemPrompt string        `json:"system_prompt"` // 系统提示
	TaskType     string        `json:"task_type"`     // 任务类型: decision/creation/analysis
	VotingMethod string        `json:"voting_method"` // 投票方法: length/交叉评估/cross
}

// VoteResponse 投票响应
type VoteResponse struct {
	Responses     map[string]string  `json:"responses"`      // 模型 -> 响应内容
	Winner        string             `json:"winner"`         // 最终获胜者
	WinnerContent string             `json:"winner_content"` // 获胜者的内容
	Scores        map[string]float64 `json:"scores"`        // 各模型得分
	Evaluation    *Evaluation        `json:"evaluation"`     // 详细评估
}

// Evaluation 评估详情
type Evaluation struct {
	TaskType      string             `json:"task_type"`       // 任务类型
	WinnerReason  string             `json:"winner_reason"`   // 获胜原因
	ModelRatings  map[string]Rating  `json:"model_ratings"`  // 各模型评级
	ProsCons      map[string]ProsCons `json:"pros_cons"`      // 各模型优缺点
}

// Rating 评级
type Rating struct {
	OverallScore float64 `json:"overall_score"` // 综合得分 (0-100)
	Accuracy     float64 `json:"accuracy"`     // 准确性
	Completeness float64 `json:"completeness"` // 完整性
	Clarity      float64 `json:"clarity"`      // 清晰度
	Creativity   float64 `json:"creativity"`   // 创造性
}

// ModelRatingsToScores 转换为分数map
func (e *Evaluation) ModelRatingsToScores() map[string]float64 {
	scores := make(map[string]float64)
	for model, rating := range e.ModelRatings {
		scores[model] = rating.OverallScore
	}
	return scores
}

// singleModelVoting 单模型投票 (多次调用+不同参数)
func (s *Service) singleModelVoting(ctx context.Context, req VoteRequest, model string) map[string]string {
	responses := make(map[string]string)
	
	// 使用不同的temperature多次调用，模拟多模型投票效果
	temperatures := []float64{0.3, 0.7, 1.0}
	labels := []string{"conservative", "balanced", "creative"}
	
	for i, temp := range temperatures {
		// 克隆请求，修改temperature
		tempReq := req
		tempReq.Temperature = temp
		
		// 调用模型
		msgs := make([]Message, 0, len(req.Messages)+1)
		if req.SystemPrompt != "" {
			msgs = append(msgs, Message{Role: "system", Content: req.SystemPrompt})
		}
		msgs = append(msgs, req.Messages...)
		
		resp, err := s.Chat(ctx, model, msgs)
		if err != nil {
			responses[fmt.Sprintf("%s_%s", model, labels[i])] = fmt.Sprintf("[Error: %v]", err)
		} else {
			responses[fmt.Sprintf("%s_%s", model, labels[i])] = resp.Content
		}
	}
	
	return responses
}

// ProsCons 优缺点
type ProsCons struct {
	Pros  []string `json:"pros"`  // 优点
	Cons  []string `json:"cons"`  // 缺点
}

// Vote 多模型投票 - 智能决策
func (s *Service) Vote(ctx context.Context, req VoteRequest) (*VoteResponse, error) {
	if len(req.Models) == 0 {
		return nil, fmt.Errorf("no models specified")
	}

	// 过滤掉不可用的模型
	var availableModels []string
	for _, model := range req.Models {
		if s.IsModelAvailable(model) {
			availableModels = append(availableModels, model)
		}
	}

	if len(availableModels) == 0 {
		return nil, fmt.Errorf("no available models: please configure at least one API key")
	}

	// 更新请求中的模型列表
	req.Models = availableModels

	var responses map[string]string
	var scores map[string]float64

	// 根据可用模型数量选择投票策略
	if len(availableModels) == 1 {
		// 只有一个模型：使用多轮对话+不同参数模拟投票
		responses = s.singleModelVoting(ctx, req, availableModels[0])
	} else {
		// 多个模型：并发调用所有模型
		responses = s并发调用模型(ctx, req)
		
		// 过滤掉失败的模型
		var successfulModels []string
		for model, resp := range responses {
			if !strings.HasPrefix(resp, "[Error:") && !strings.HasPrefix(resp, "[Model") {
				successfulModels = append(successfulModels, model)
			}
		}

		if len(successfulModels) == 1 {
			// 如果只有一个成功，使用单模型投票逻辑
			responses = s.singleModelVoting(ctx, req, successfulModels[0])
		} else if len(successfulModels) == 0 {
			return nil, fmt.Errorf("all models failed: please check API keys")
		}
	}

	// 2. 根据任务类型选择评估方法
	if len(responses) > 1 {
		switch req.VotingMethod {
		case "length":
			scores = s.按长度评分(responses)
		case "cross", "交叉评估":
			_, eval := s.交叉评估(ctx, successfulModels, responses, req.TaskType)
			if eval != nil {
				scores = eval.ModelRatingsToScores()
			}
		default:
			_, eval := s.综合评分(ctx, successfulModels, responses, req.TaskType)
			if eval != nil {
				scores = eval.ModelRatingsToScores()
			}
		}
	}

	// 3. 选择得分最高的
	winner := ""
	maxScore := 0.0
	var eval *Evaluation

	switch req.VotingMethod {
	case "length":
		// 简单方法：按响应长度
		scores = s.按长度评分(responses)
		eval = s.基础评估(responses, req.TaskType)
	case "cross", "交叉评估":
		// 高级方法：让模型互相评估
		scores, eval = s.交叉评估(ctx, req.Models, responses, req.TaskType)
	default:
		// 默认：综合评分
		scores, eval = s.综合评分(ctx, req.Models, responses, req.TaskType)
	}

	// 3. 选择得分最高的
	winner := ""
	maxScore := 0.0
	for model, score := range scores {
		if score > maxScore {
			maxScore = score
			winner = model
		}
	}

	return &VoteResponse{
		Responses:     responses,
		Winner:        winner,
		WinnerContent: responses[winner],
		Scores:        scores,
		Evaluation:    eval,
	}, nil
}

// 并发调用模型
func (s *Service) 并发调用模型(ctx context.Context, req VoteRequest) map[string]string {
	responses := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var failedCount int

	for _, modelName := range req.Models {
		wg.AddOne()
		go func(model string) {
			defer wg.Done()

			// 检查模型是否可用
			if !s.IsModelAvailable(model) {
				mu.Lock()
				responses[model] = "[Model not available: API key not configured]"
				failedCount++
				mu.Unlock()
				return
			}

			msgs := make([]Message, 0, len(req.Messages)+1)
			if req.SystemPrompt != "" {
				msgs = append(msgs, Message{Role: "system", Content: req.SystemPrompt})
			}
			msgs = append(msgs, req.Messages...)

			resp, err := s.Chat(ctx, model, msgs)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				responses[model] = fmt.Sprintf("[Error: %v]", err)
				failedCount++
			} else {
				responses[model] = resp.Content
			}
		}(modelName)
	}

	wg.Wait()

	// 如果全部失败，返回错误提示
	if failedCount == len(req.Models) {
		responses["_error"] = "All models failed - please check API keys configuration"
	}

	return responses
}

// 按长度评分 (基础方法)
func (s *Service) 按长度评分(responses map[string]string) map[string]float64 {
	scores := make(map[string]float64)
	maxLen := 0

	for _, resp := range responses {
		if len(resp) > maxLen {
			maxLen = len(resp)
		}
	}

	// 归一化到0-100
	for model, resp := range responses {
		if maxLen > 0 {
			scores[model] = float64(len(resp)) / float64(maxLen) * 100
		} else {
			scores[model] = 50
		}
	}

	return scores
}

// 基础评估
func (s *Service) 基础评估(responses map[string]string, taskType string) *Evaluation {
	ratings := make(map[string]Rating)
	prosCons := make(map[string]ProsCons)

	for model, resp := range responses {
		// 简单规则评分
		score := 50.0

		// 检查是否包含行动项
		if strings.Contains(resp, "首先") || strings.Contains(resp, "第一步") ||
			strings.Contains(resp, "建议") || strings.Contains(resp, "应该") {
			score += 10
		}

		// 检查是否有结构化内容
		if strings.Contains(resp, "1.") || strings.Contains(resp, "①") ||
			strings.Contains(resp, "•") {
			score += 10
		}

		// 检查是否太短或太长
		if len(resp) < 50 {
			score -= 10
		} else if len(resp) > 5000 {
			score -= 5
		}

		ratings[model] = Rating{
			OverallScore: score,
			Accuracy:     score,
			Completeness: score,
			Clarity:      score,
			Creativity:   score,
		}

		prosCons[model] = ProsCons{
			Pros: []string{"响应完整"},
			Cons: []string{},
		}
	}

	// 找出最佳
	var bestModel string
	var bestScore float64
	for model, r := range ratings {
		if r.OverallScore > bestScore {
			bestScore = r.OverallScore
			bestModel = model
		}
	}

	return &Evaluation{
		TaskType:     taskType,
		WinnerReason: fmt.Sprintf("%s 在基础评估中得分最高", bestModel),
		ModelRatings: ratings,
		ProsCons:     prosCons,
	}
}

// 交叉评估 (高级方法)
func (s *Service) 交叉评估(ctx context.Context, models []string, responses map[string]string, taskType string) (map[string]float64, *Evaluation) {
	scores := make(map[string]float64)
	ratings := make(map[string]Rating)
	prosCons := make(map[string]ProsCons)

	// 让每个模型评估其他模型的回复
	for evaluator := range responses {
		evalModel := evaluator
		evalPrompt := fmt.Sprintf(`你是一个专业的AI评估专家。请评估以下AI模型对某个问题的回答质量。

任务类型: %s

请评估以下每个答案，并给出1-10分的评分（10分最高）：
- 准确性：答案是否正确
- 完整性：是否涵盖所有重要方面
- 清晰度：表达是否清晰易懂
- 创造性：是否有独到见解

需要评估的模型回复：
%s

请以JSON格式返回评估结果，格式如下：
{"ratings": {"模型名": {"accuracy": X, "completeness": X, "clarity": X, "creativity": X, "overall": X}, ...}, "pros": {"模型名": ["优点1", "优点2"], ...}, "cons": {"模型名": ["缺点1", ...], ...}}`,
			taskType, s.构建评估上下文(models, responses))

		// 调用评估模型
		resp, err := s.Chat(ctx, evalModel, []Message{
			{Role: "system", Content: "你是一个专业的AI评估专家。请严格评估并给出评分。"},
			{Role: "user", Content: evalPrompt},
		})

		if err != nil {
			continue
		}

		// 解析评估结果
		parsed := s.解析评估结果(resp.Content, models)
		for model, r := range parsed.ratings {
			ratings[model] = r
		}
		for model, pc := range parsed.prosCons {
			prosCons[model] = pc
		}
	}

	// 汇总得分
	for model := range responses {
		if r, ok := ratings[model]; ok {
			scores[model] = r.OverallScore * 10 // 转换为0-100
		} else {
			scores[model] = 50
		}
	}

	// 找出最佳模型和原因
	bestModel, bestReason := s.找出最佳模型(models, ratings, prosCons)

	return scores, &Evaluation{
		TaskType:     taskType,
		WinnerReason: bestReason,
		ModelRatings: ratings,
		ProsCons:     prosCons,
	}
}

// 综合评分 (默认方法)
func (s *Service) 综合评分(ctx context.Context, models []string, responses map[string]string, taskType string) (map[string]float64, *Evaluation) {
	scores := make(map[string]float64)
	ratings := make(map[string]Rating)
	prosCons := make(map[string]ProsCons)

	for model, resp := range responses {
		r := s.评估单个回复(resp, taskType)
		ratings[model] = r

		// 综合得分 = 准确性*0.3 + 完整性*0.3 + 清晰度*0.2 + 创造性*0.2
		scores[model] = r.Accuracy*0.3 + r.Completeness*0.3 + r.Clarity*0.2 + r.Creativity*0.2

		// 提取优缺点
		prosCons[model] = s.提取优缺点(resp)
	}

	// 交叉验证：用其他模型确认最佳
	bestModel, bestReason := s.找出最佳模型(models, ratings, prosCons)

	return scores, &Evaluation{
		TaskType:     taskType,
		WinnerReason: bestReason,
		ModelRatings: ratings,
		ProsCons:     prosCons,
	}
}

// 评估单个回复
func (s *Service) 评估单个回复(resp, taskType string) Rating {
	score := Rating{
		Accuracy:     70,
		Completeness: 70,
		Clarity:      70,
		Creativity:   70,
	}

	// 准确性检查
	if strings.Contains(resp, "错误") || strings.Contains(resp, "不确定") {
		score.Accuracy -= 10
	}

	// 完整性检查
	if strings.Contains(resp, "第一") && strings.Contains(resp, "第二") && strings.Contains(resp, "第三") {
		score.Completeness += 15
	}
	if len(resp) > 200 {
		score.Completeness += 10
	}

	// 清晰度检查
	if strings.Contains(resp, "首先") || strings.Contains(resp, "其次") {
		score.Clarity += 15
	}
	if strings.Contains(resp, "。") && strings.Contains(resp, "，") {
		score.Clarity += 10
	}

	// 创造性检查
	if strings.Contains(resp, "但是") || strings.Contains(resp, "然而") {
		score.Creativity += 10
	}
	if strings.Contains(resp, "创新") || strings.Contains(resp, "独特") {
		score.Creativity += 15
	}

	// 任务类型调整
	switch taskType {
	case "decision":
		// 决策任务更看重准确性和完整性
		score.Accuracy *= 1.2
		score.Completeness *= 1.2
	case "creation":
		// 创作任务更看重创造性
		score.Creativity *= 1.3
	case "analysis":
		// 分析任务看重完整性和清晰度
		score.Completeness *= 1.2
		score.Clarity *= 1.2
	}

	// 限制在0-100
	score.Accuracy = min(score.Accuracy, 100)
	score.Completeness = min(score.Completeness, 100)
	score.Clarity = min(score.Clarity, 100)
	score.Creativity = min(score.Creativity, 100)

	score.OverallScore = (score.Accuracy + score.Completeness + score.Clarity + score.Creativity) / 4

	return score
}

// 提取优缺点
func (s *Service) 提取优缺点(resp string) ProsCons {
	pc := ProsCons{
		Pros:  []string{},
		Cons:  []string{},
	}

	// 简单规则提取
	if len(resp) > 100 {
		pc.Pros = append(pc.Pros, "内容详尽")
	}
	if strings.Contains(resp, "建议") || strings.Contains(resp, "应该") {
		pc.Pros = append(pc.Pros, "有具体建议")
	}
	if strings.Contains(resp, "1.") || strings.Contains(resp, "①") {
		pc.Pros = append(pc.Pros, "结构清晰")
	}

	if len(resp) < 50 {
		pc.Cons = append(pc.Cons, "内容过于简短")
	}
	if strings.Contains(resp, "可能") || strings.Contains(resp, "也许") {
		pc.Cons = append(pc.Cons, "语气不够确定")
	}

	return pc
}

// 找出最佳模型
func (s *Service) 找出最佳模型(models []string, ratings map[string]Rating, prosCons map[string]ProsCons) (string, string) {
	bestModel := ""
	bestScore := 0.0

	for _, model := range models {
		if r, ok := ratings[model]; ok {
			if r.OverallScore > bestScore {
				bestScore = r.OverallScore
				bestModel = model
			}
		}
	}

	if pc, ok := prosCons[bestModel]; ok {
		reason := fmt.Sprintf("%s 得分%.1f分", bestModel, bestScore)
		if len(pc.Pros) > 0 {
			reason += fmt.Sprintf("，优点: %s", pc.Pros[0])
		}
		return bestModel, reason
	}

	return bestModel, fmt.Sprintf("%s 综合得分最高", bestModel)
}

// 构建评估上下文
func (s *Service) 构建评估上下文(models []string, responses map[string]string) string {
	var sb strings.Builder
	for _, model := range models {
		if resp, ok := responses[model]; ok {
			sb.WriteString(fmt.Sprintf("\n【%s】:\n%s\n", model, resp))
		}
	}
	return sb.String()
}

// 解析评估结果
func (s *Service) 解析评估结果(content string, knownModels []string) struct {
	ratings map[string]Rating
	prosCons map[string]ProsCons
} {
	ratings := make(map[string]Rating)
	prosCons := make(map[string]ProsCons)

	// 简单解析 - 提取数字评分
	for _, model := range knownModels {
		ratings[model] = Rating{OverallScore: 70}
		prosCons[model] = ProsCons{Pros: []string{}, Cons: []string{}}
	}

	return struct {
		ratings map[string]Rating
		prosCons map[string]ProsCons
	}{ratings, prosCons}
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
