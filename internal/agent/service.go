package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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

// ProcessWithVoting 使用智能体投票处理
// TODO: 迁移到多智能体协作模式
func (s *Service) ProcessWithVoting(ctx context.Context, cfg Config, input string) (string, error) {
	// 检查是否启用协作模式
	if cfg.CollaborationMode {
		return s.ProcessWithCollaboration(ctx, cfg, input)
	}
	
	// 降级为普通单智能体处理
	return s.callModel(ctx, cfg.Model, cfg.SystemPrompt, input)
}

// ProcessWithCollaboration 多智能体协作处理
// CEO 分解任务 → Manager 分配任务 → Worker 执行 → Manager 汇总 → CEO 最终决策
// 支持动态模型切换：下级投票换领导模型，领导根据业绩换下属模型
func (s *Service) ProcessWithCollaboration(ctx context.Context, cfg Config, input string) (string, error) {
	// 初始化模型
	ceoModel := cfg.Model
	managerModel := cfg.Model
	workerModel := cfg.Model
	
	availableModels := []string{"gpt-4", "glm-4", "claude-3-sonnet", "kimi"}
	
	// 1. CEO 角色：分析任务，分解为子任务，同时评估是否需要更换模型
	ceoPrompt := cfg.CEOPrompt
	if ceoPrompt == "" {
		ceoPrompt = "你是一个CEO，负责分析用户需求并分解任务。"
	}
	
	// 下级投票决定是否更换 CEO 模型
	ceoModel, _ = s.voteForModel(ctx, availableModels, "CEO", ceoModel)
	
	ceoResponse, err := s.callModel(ctx, ceoModel, ceoPrompt, input)
	if err != nil {
		return "", fmt.Errorf("CEO 分析失败: %v", err)
	}
	
	// 2. Manager 角色：接收CEO分解的子任务，分配给Worker
	managerPrompt := cfg.ManagerPrompt
	if managerPrompt == "" {
		managerPrompt = "你是一个Manager，负责将任务分配给具体的执行者。"
	}
	managerResponse, err := s.callModel(ctx, managerModel, managerPrompt, 
		"CEO分析结果: "+ceoResponse+"\n\n原始任务: "+input)
	if err != nil {
		return "", fmt.Errorf("Manager 分配失败: %v", err)
	}
	
	// 3. Worker 角色：执行具体任务
	workerPrompt := cfg.WorkerPrompt
	if workerPrompt == "" {
		workerPrompt = "你是一个Worker，负责执行具体的任务。"
	}
	workerResponse, err := s.callModel(ctx, workerModel, workerPrompt, 
		"Manager安排: "+managerResponse+"\n\n请执行具体任务并给出结果。")
	if err != nil {
		return "", fmt.Errorf("Worker 执行失败: %v", err)
	}
	
	// 4. 领导根据业绩（响应质量）决定是否更换下属模型
	managerModel = s.evaluateAndSwitchModel(ctx, availableModels, "Manager", managerModel, managerResponse)
	workerModel = s.evaluateAndSwitchModel(ctx, availableModels, "Worker", workerModel, workerResponse)
	
	// 5. CEO 最终决策：汇总所有结果，给出最终答案
	finalPrompt := "你是一个CEO，负责整合所有信息给出最终答案。"
	finalResponse, err := s.callModel(ctx, ceoModel, finalPrompt, 
		"CEO分析: "+ceoResponse+"\n\nManager安排: "+managerResponse+"\n\nWorker执行结果: "+workerResponse+"\n\n原始任务: "+input)
	if err != nil {
		return "", fmt.Errorf("CEO 最终决策失败: %v", err)
	}
	
	// 汇总完整流程（包含模型切换信息）
	result := fmt.Sprintf("【模型状态】\nCEO: %s | Manager: %s | Worker: %s\n\n【任务分析】\n%s\n\n【任务分配】\n%s\n\n【执行结果】\n%s\n\n【最终方案】\n%s", 
		ceoModel, managerModel, workerModel, ceoResponse, managerResponse, workerResponse, finalResponse)
	
	return result, nil
}

// voteForModel 下级投票选择最佳模型
func (s *Service) voteForModel(ctx context.Context, models []string, role, currentModel string) (string, error) {
	if len(models) < 2 {
		return currentModel, nil
	}
	
	// 模拟下级投票：让其他模型评估当前模型
	var bestModel string
	bestScore := 0.0
	
	votePrompt := fmt.Sprintf("作为%s的下级，请投票选择最适合当前任务的AI模型。候选模型: %v。当前模型: %s。请直接返回你认为最佳的模型名称。", role, models, currentModel)
	
	for _, model := range models {
		response, err := s.callModel(ctx, model, "你是一个公正的评审", votePrompt)
		if err != nil {
			continue
		}
		// 简单评分逻辑
		score := float64(len(response)) / 10 // 粗略评分
		if score > bestScore {
			bestScore = score
			bestModel = model
		}
	}
	
	if bestModel == "" {
		return currentModel, nil
	}
	
	return bestModel, nil
}

// evaluateAndSwitchModel 领导根据业绩决定是否更换下属模型
func (s *Service) evaluateAndSwitchModel(ctx context.Context, models []string, role, currentModel, performance string) string {
	if len(models) < 2 {
		return currentModel
	}
	
	// 让领导评估下属表现
	evalPrompt := fmt.Sprintf("作为%s的领导，请评估下属的工作质量(0-100分)，并决定是否需要更换模型。直接返回: 分数,是否换模型(yes/no)", role)
	
	response, err := s.callModel(ctx, currentModel, "你是一个严格的领导", evalPrompt+"\n\n工作成果: "+performance)
	if err != nil {
		return currentModel
	}
	
	// 解析响应，判断是否需要更换
	// 简单实现：如果响应中包含 "yes" 或 "更换" 则更换模型
	if strings.Contains(response, "yes") || strings.Contains(response, "更换") || strings.Contains(response, "换") {
		// 选择一个不同的模型
		for _, model := range models {
			if model != currentModel {
				return model
			}
		}
	}
	
	return currentModel
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
	Model       string                 `json:"model"`        // 默认模型
	Provider    string                 `json:"provider"`    // openai/anthropic/custom
	SystemPrompt string                `json:"system_prompt"`
	Tools       []string               `json:"tools"`       // 工具列表
	Temperature float64                `json:"temperature"`
	MaxTokens   int                    `json:"max_tokens"`
	
	// 多智能体协作配置
	CollaborationMode bool              `json:"collaboration_mode"`  // 是否启用多智能体协作
	CEOPrompt         string            `json:"ceo_prompt"`          // CEO 角色提示
	ManagerPrompt     string            `json:"manager_prompt"`      // Manager 角色提示
	WorkerPrompt      string            `json:"worker_prompt"`       // Worker 角色提示
	AvailableModels   []string          `json:"available_models"`    // 可用模型列表（用于动态切换）
	EnableModelVote   bool              `json:"enable_model_vote"`   // 启用下级投票换领导模型
	EnablePerfSwitch  bool              `json:"enable_perf_switch"`  // 启用领导根据业绩换下属模型
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

// ========== 定时绩效评估 ==========

// PerformanceEvalConfig 绩效评估定时任务配置
type PerformanceEvalConfig struct {
	EvalInterval time.Duration // 评估间隔，如 1小时、每天
	Models       []string      // 可用模型列表
	Workers      []string      // 待评估的下属列表: ["Manager", "Worker"]
	CompanyGoal  string        // 公司目标/业绩指标
}

// StartPerformanceEvalScheduler 启动绩效评估定时调度器
func (s *Service) StartPerformanceEvalScheduler(cfg PerformanceEvalConfig) {
	go func() {
		ticker := time.NewTicker(cfg.EvalInterval)
		defer ticker.Stop()
		
		for range ticker.C {
			s.runCompanyPerformanceEval(cfg)
		}
	}()
}

// runCompanyPerformanceEval 从公司整体业绩出发进行绩效评估
// 流程：评估公司业绩 → 反向追溯各层级贡献 → 调整模型配置
func (s *Service) runCompanyPerformanceEval(cfg PerformanceEvalConfig) {
	ctx := context.Background()
	
	// 1. 获取公司整体业绩
	companyPerformance := s.getCompanyPerformance()
	
	// 2. 根据公司业绩，反向追溯各层级的贡献和表现
	fmt.Printf("\n========== 公司绩效评估 ==========\n")
	fmt.Printf("公司整体业绩: %s\n", companyPerformance)
	fmt.Printf("目标: %s\n\n", cfg.CompanyGoal)
	
	// 3. CEO 评估：基于公司业绩，决定是否调整战略/模型
	ceoModel := s.defaultModel
	ceoEvaluation := s.evaluateRolePerformance(ctx, "CEO", ceoModel, companyPerformance, cfg.CompanyGoal)
	if ceoEvaluation.shouldSwitch {
		ceoModel = ceoEvaluation.newModel
		fmt.Printf("→ CEO 模型调整为: %s (原因: %s)\n", ceoModel, ceoEvaluation.reason)
	}
	
	// 4. Manager 评估：CEO 基于公司业绩评估 Manager
	managerModel := s.defaultModel
	managerEvaluation := s.evaluateRolePerformance(ctx, "Manager", managerModel, companyPerformance, cfg.CompanyGoal)
	if managerEvaluation.shouldSwitch {
		managerModel = managerEvaluation.newModel
		fmt.Printf("→ Manager 模型调整为: %s (原因: %s)\n", managerModel, managerEvaluation.reason)
	}
	
	// 5. Worker 评估：Manager 基于公司业绩评估 Worker
	workerModel := s.defaultModel
	workerEvaluation := s.evaluateRolePerformance(ctx, "Worker", workerModel, companyPerformance, cfg.CompanyGoal)
	if workerEvaluation.shouldSwitch {
		workerModel = workerEvaluation.newModel
		fmt.Printf("→ Worker 模型调整为: %s (原因: %s)\n", workerModel, workerEvaluation.reason)
	}
	
	fmt.Printf("=====================================\n\n")
	
	// TODO: 保存评估结果到数据库
}

// EvaluationResult 评估结果
type EvaluationResult struct {
	shouldSwitch bool
	newModel     string
	reason       string
	score        int
}

// evaluateRolePerformance 评估某个角色的表现（从公司业绩出发）
func (s *Service) evaluateRolePerformance(ctx context.Context, role, currentModel, companyPerformance, goal string) EvaluationResult {
	evalPrompt := fmt.Sprintf(`你作为公司的%s，需要根据公司整体业绩来评估自己的工作表现和模型配置是否合适。

公司整体业绩/成果:
%s

公司目标:
%s

请从以下角度评估:
1. 自己的工作是否对公司业绩有贡献?
2. 当前使用的AI模型(%s)是否最适合当前任务?
3. 是否需要更换模型来提升业绩?

请直接返回以下格式(每行一个):
分数(0-100):
是否更换模型(yes/no):
更换原因(一句话):
推荐模型(如果是yes):
`, role, companyPerformance, goal, currentModel)

	response, err := s.callModel(ctx, currentModel, "你是一个追求业绩的领导者", evalPrompt)
	if err != nil {
		return EvaluationResult{shouldSwitch: false, newModel: currentModel, reason: "评估失败"}
	}
	
	// 解析响应
	result := EvaluationResult{
		shouldSwitch: strings.Contains(response, "yes"),
		newModel:     currentModel,
		reason:       "保持当前模型",
	}
	
	// 提取推荐模型
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "推荐模型") || strings.Contains(line, "gpt-") || strings.Contains(line, "glm-") || strings.Contains(line, "claude") {
			for _, m := range []string{"gpt-4", "glm-4", "claude-3-sonnet", "kimi"} {
				if strings.Contains(line, m) {
					result.newModel = m
					break
				}
			}
		}
		if strings.Contains(line, "分数") {
			fmt.Sscanf(line, "分数%d", &result.score)
		}
		if strings.Contains(line, "更换原因") {
			result.reason = strings.Split(line, ":")[1]
		}
	}
	
	return result
}

// getCompanyPerformance 获取公司整体业绩
func (s *Service) getCompanyPerformance() string {
	// TODO: 从数据库获取公司实际业绩数据
	// 这里返回模拟数据
	return "本月完成产品上线3个，用户增长20%，收入增长15%"
}

// getWorkerPerformance 获取下属的工作表现数据
func (s *Service) getWorkerPerformance(worker string) string {
	// TODO: 从数据库或记忆服务获取下属最近的工作记录
	// 这里返回模拟数据
	return fmt.Sprintf("%s 最近完成了代码审查、测试生成等任务，工作质量良好", worker)
}

// getLeaderRole 获取下属对应的领导角色
func (s *Service) getLeaderRole(worker string) string {
	switch worker {
	case "Manager":
		return "CEO"
	case "Worker":
		return "Manager"
	default:
		return "CEO"
	}
}
