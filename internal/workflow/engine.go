package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"agent-flow/internal/agent"
	"agent-flow/internal/store"
)

// NodeType 节点类型
type NodeType string

const (
	NodeTypeTrigger   NodeType = "trigger"
	NodeTypeAgent     NodeType = "agent"
	NodeTypeCondition NodeType = "condition"
	NodeTypeTool      NodeType = "tool"
	NodeTypeLLM       NodeType = "llm"
)

// Node 流程节点
type Node struct {
	ID       string                 `json:"id"`
	Type     NodeType              `json:"type"`
	Data     map[string]interface{} `json:"data"`
	Position Position               `json:"position"`
}

// Position 位置
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Edge 流程连线
type Edge struct {
	ID         string `json:"id"`
	Source     string `json:"source"`
	Target     string `json:"target"`
	SourceHandle string `json:"sourceHandle,omitempty"`
	TargetHandle string `json:"targetHandle,omitempty"`
	Condition  string `json:"condition,omitempty"` // 条件分支
}

// Flow 流程
type Flow struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Nodes   []Node   `json:"nodes"`
	Edges   []Edge   `json:"edges"`
	Enabled bool     `json:"enabled"`
}

// Engine 流程执行引擎
type Engine struct {
	db        *store.Postgres
	redis     *store.Redis
	agentSvc  *agent.Service
	nodeMutex sync.Map // 节点级别锁
}

// NewEngine 创建流程引擎
func NewEngine(db *store.Postgres, redis *store.Redis, agentSvc *agent.Service) *Engine {
	return &Engine{
		db:       db,
		redis:    redis,
		agentSvc: agentSvc,
	}
}

// ExecuteRequest 执行请求
type ExecuteRequest struct {
	FlowID    string                 `json:"flow_id"`
	Input     string                 `json:"input"`
	UserID    string                 `json:"user_id"`
	ChannelID string                 `json:"channel_id"`
	Context   map[string]interface{} `json:"context"`
}

// ExecuteResponse 执行响应
type ExecuteResponse struct {
	FlowID   string                 `json:"flow_id"`
	Output   string                 `json:"output"`
	NodesExec []NodeExecution       `json:"nodes_exec"`
	Context   map[string]interface{} `json:"context"`
}

// NodeExecution 节点执行记录
type NodeExecution struct {
	NodeID   string                 `json:"node_id"`
	NodeType NodeType              `json:"node_type"`
	Input    string                 `json:"input"`
	Output   string                 `json:"output"`
	Error    string                 `json:"error,omitempty"`
	Duration int64                  `json:"duration_ms"`
}

// Execute 执行流程
func (e *Engine) Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	// 获取流程配置
	flow, err := e.getFlow(req.FlowID)
	if err != nil {
		return nil, fmt.Errorf("flow not found: %w", err)
	}

	if !flow.Enabled {
		return nil, fmt.Errorf("flow is disabled")
	}

	// 构建执行图
	graph := e.buildGraph(flow)

	// 找起始节点 (触发器)
	startNodes := e.findStartNodes(flow)
	if len(startNodes) == 0 {
		return nil, fmt.Errorf("no trigger node found")
	}

	// 执行上下文
	execCtx := &ExecutionContext{
		FlowID:    req.FlowID,
		UserID:    req.UserID,
		ChannelID: req.ChannelID,
		Input:     req.Input,
		Context:   req.Context,
		Variables: make(map[string]interface{}),
		Results:   make(map[string]string),
	}

	// 从触发器开始执行
	var results []NodeExecution
	for _, startNode := range startNodes {
		nodeResults, err := e.executeNode(ctx, graph, startNode, execCtx)
		if err != nil {
			log.Printf("Node %s execution error: %v", startNode.ID, err)
			continue
		}
		results = append(results, nodeResults...)
	}

	return &ExecuteResponse{
		FlowID:    req.FlowID,
		Output:    execCtx.Output,
		NodesExec: results,
		Context:   execCtx.Context,
	}, nil
}

// ExecutionContext 执行上下文
type ExecutionContext struct {
	FlowID    string
	UserID    string
	ChannelID string
	Input     string
	Output    string
	Context   map[string]interface{}
	Variables map[string]interface{}
	Results   map[string]string // 节点ID -> 输出
	mu        sync.RWMutex
}

func (ec *ExecutionContext) SetVar(key string, value interface{}) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.Variables[key] = value
}

func (ec *ExecutionContext) GetVar(key string) interface{} {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return ec.Variables[key]
}

func (ec *ExecutionContext) SetResult(nodeID, result string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.Results[nodeID] = result
}

func (ec *ExecutionContext) GetResult(nodeID string) string {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return ec.Results[nodeID]
}

// NodeGraph 节点图
type NodeGraph map[string][]string // nodeID -> 子节点IDs

// buildGraph 构建执行图
func (e *Engine) buildGraph(flow *Flow) NodeGraph {
	graph := make(NodeGraph)
	for _, edge := range flow.Edges {
		graph[edge.Source] = append(graph[edge.Source], edge.Target)
	}
	return graph
}

// findStartNodes 找起始节点
func (e *Engine) findStartNodes(flow *Flow) []Node {
	var triggers []Node
	for _, node := range flow.Nodes {
		if node.Type == NodeTypeTrigger {
			triggers = append(triggers, node)
		}
	}
	return triggers
}

// executeNode 执行单个节点
func (e *Engine) executeNode(ctx context.Context, graph NodeGraph, node Node, execCtx *ExecutionContext) ([]NodeExecution, error) {
	// 获取节点锁 (防止并发执行同一节点)
	lockKey := fmt.Sprintf("flow:%s:node:%s", execCtx.FlowID, node.ID)
	execCtx.mu.Lock()
	// TODO: 使用Redis分布式锁
	execCtx.mu.Unlock()

	var result string
	var err error
	var duration int64

	switch node.Type {
	case NodeTypeTrigger:
		result, err = e.executeTrigger(node, execCtx)
	case NodeTypeAgent:
		result, err = e.executeAgent(node, execCtx)
	case NodeTypeCondition:
		result, err = e.executeCondition(node, graph, execCtx)
	case NodeTypeTool:
		result, err = e.executeTool(node, execCtx)
	case NodeTypeLLM:
		result, err = e.executeLLM(node, execCtx)
	default:
		err = fmt.Errorf("unknown node type: %s", node.Type)
	}

	execCtx.SetResult(node.ID, result)

	execution := NodeExecution{
		NodeID:   node.ID,
		NodeType: node.Type,
		Input:    execCtx.GetVar("node_input_" + node.ID).(string),
		Output:   result,
		Duration: duration,
	}
	if err != nil {
		execution.Error = err.Error()
	}

	var results []NodeExecution
	results = append(results, execution)

	// 执行子节点
	children := graph[node.ID]
	for _, childID := range children {
		childNode := e.findNode(flow, childID)
		if childNode == nil {
			continue
		}

		// 条件分支检查
		edge := e.findEdge(flow, node.ID, childID)
		if edge != nil && edge.Condition != "" {
			if !e.evaluateCondition(edge.Condition, execCtx) {
				continue
			}
		}

		childResults, err := e.executeNode(ctx, graph, *childNode, execCtx)
		if err != nil {
			log.Printf("Child node %s error: %v", childID, err)
			continue
		}
		results = append(results, childResults...)
	}

	return results, nil
}

// findNode 查找节点
func (e *Engine) findNode(flow *Flow, nodeID string) *Node {
	for i := range flow.Nodes {
		if flow.Nodes[i].ID == nodeID {
			return &flow.Nodes[i]
		}
	}
	return nil
}

// findEdge 查找连线
func (e *Engine) findEdge(flow *Flow, source, target string) *Edge {
	for i := range flow.Edges {
		if flow.Edges[i].Source == source && flow.Edges[i].Target == target {
			return &flow.Edges[i]
		}
	}
	return nil
}

// executeTrigger 执行触发器
func (e *Engine) executeTrigger(node Node, execCtx *ExecutionContext) (string, error) {
	triggerType, _ := node.Data["triggerType"].(string)
	execCtx.SetVar("node_input_"+node.ID, execCtx.Input)
	
	switch triggerType {
	case "用户消息", "message":
		return execCtx.Input, nil
	case "定时", "schedule":
		return "triggered", nil
	case "webhook":
		return "webhook triggered", nil
	default:
		return execCtx.Input, nil
	}
}

// executeAgent 执行智能体节点
func (e *Engine) executeAgent(node Node, execCtx *ExecutionContext) (string, error) {
	agentID, _ := node.Data["agentId"].(string)
	prevResult := e.getPreviousNode(execCtx, node.ID)
	input := execCtx.GetResult(prevResult)
	
	execCtx.SetVar("node_input_"+node.ID, input)

	if agentID == "" {
		// 使用默认Agent
		return e.agentSvc.Process(context.Background(), input, execCtx.UserID)
	}

	return e.agentSvc.ProcessWithAgent(context.Background(), agentID, input, execCtx.UserID)
}

// executeCondition 执行条件分支
func (e *Engine) executeCondition(node Node, graph NodeGraph, execCtx *ExecutionContext) (string, error) {
	condition, _ := node.Data["condition"].(string)
	execCtx.SetVar("node_input_"+node.ID, execCtx.Input)

	result := e.evaluateCondition(condition, execCtx)
	return fmt.Sprintf("%t", result), nil
}

// evaluateCondition 评估条件
func (e *Engine) evaluateCondition(condition string, execCtx *ExecutionContext) bool {
	if condition == "" {
		return true
	}

	// 简单条件解析 (后续可扩展)
	// 例如: "input contains hello" -> 检查输入是否包含hello
	//       "variable_exists:user_id" -> 检查变量是否存在
	
	switch {
	case condition == "always":
		return true
	case condition == "never":
		return false
	default:
		// 默认为true，继续执行
		return true
	}
}

// executeTool 执行工具节点
func (e *Engine) executeTool(node Node, execCtx *ExecutionContext) (string, error) {
	toolType, _ := node.Data["toolType"].(string)
	toolName, _ := node.Data["toolName"].(string)
	input := execCtx.GetResult(e.getPreviousNode(execCtx, node.ID))
	
	execCtx.SetVar("node_input_"+node.ID, input)

	// TODO: 调用实际工具
	log.Printf("[Tool] Executing %s (%s) with input: %s", toolName, toolType, input)

	return fmt.Sprintf("tool %s executed", toolName), nil
}

// executeLLM 执行大模型节点
func (e *Engine) executeLLM(node Node, execCtx *ExecutionContext) (string, error) {
	prompt, _ := node.Data["prompt"].(string)
	model, _ := node.Data["model"].(string)
	input := execCtx.GetResult(e.getPreviousNode(execCtx, node.ID))
	
	execCtx.SetVar("node_input_"+node.ID, input)

	// 构建完整prompt
	fullPrompt := prompt + "\n\n输入: " + input

	return e.agentSvc.CallLLM(ctx, model, fullPrompt)
}

// getPreviousNode 获取上一节点
func (e *Engine) getPreviousNode(execCtx *ExecutionContext, nodeID string) string {
	// 简化实现：查找最近的结果
	for k, v := range execCtx.Results {
		if k != nodeID {
			return k
		}
	}
	return ""
}

// getFlow 获取流程配置
func (e *Engine) getFlow(flowID string) (*Flow, error) {
	// 从DB或缓存获取
	id := 1 // TODO: 解析flowID
	flowData, err := e.db.GetFlow(id)
	if err != nil {
		return nil, err
	}

	var flow Flow
	if err := json.Unmarshal([]byte(flowData.Nodes), &flow.Nodes); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(flowData.Edges), &flow.Edges); err != nil {
		return nil, err
	}
	flow.ID = flowData.Name
	flow.Name = flowData.Name
	flow.Enabled = flowData.Enabled

	return &flow, nil
}
