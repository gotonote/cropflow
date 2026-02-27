package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"corpflow/internal/store"
)

// MemoryType 记忆类型
type MemoryType string

const (
	MemoryTypeAction   MemoryType = "action"   // 执行的动作
	MemoryTypeDecision MemoryType = "decision" // 做出的决策
	MemoryTypeResult   MemoryType = "result"   // 执行结果
	MemoryTypeLearn   MemoryType = "learn"    // 学到的知识
	MemoryTypeEvent   MemoryType = "event"    // 重要事件
)

// Memory 记忆
type Memory struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	AgentID     uint           `json:"agent_id" gorm:"index"`       // 智能体ID
	ParentID    *uint         `json:"parent_id" gorm:"index"`       // 上级智能体ID (nil表示顶级)
	Type        MemoryType    `json:"type"`                        // 记忆类型
	Content     string        `json:"content"`                     // 记忆内容
	Metadata    string        `json:"metadata"`                    // 附加数据 (JSON)
	Importance  int           `json:"importance"`                  // 重要程度 1-10
	CreatedAt   time.Time     `json:"created_at"`
}

// Knowledge 知识 (可被上级查看)
type Knowledge struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	AgentID     uint           `json:"agent_id" gorm:"index"`       // 智能体ID
	ParentID    *uint         `json:"parent_id" gorm:"index"`       // 上级智能体ID
	Title       string        `json:"title"`                       // 知识标题
	Content     string        `json:"content"`                     // 知识内容
	Tags        string        `json:"tags"`                        // 标签 (逗号分隔)
	AccessLevel int           `json:"access_level"`                 // 访问级别 1-10
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// AgentRelationship 智能体关系
type AgentRelationship struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ParentID     uint      `json:"parent_id" gorm:"index"`        // 上级ID
	ChildID      uint      `json:"child_id" gorm:"index"`         // 下级ID
	RelationType string    `json:"relation_type"`                // manage/delegate/collaborate
	CreatedAt    time.Time `json:"created_at"`
}

// Service 记忆服务
type Service struct {
	db    *store.Postgres
	redis *store.Redis
}

// NewService 创建记忆服务
func NewService(db *store.Postgres, redis *store.Redis) *Service {
	return &Service{
		db:    db,
		redis: redis,
	}
}

// ========== 记忆管理 ==========

// AddMemory 添加记忆
func (s *Service) AddMemory(agentID uint, parentID *uint, memType MemoryType, content string, importance int) (*Memory, error) {
	memory := &Memory{
		AgentID:    agentID,
		ParentID:   parentID,
		Type:       memType,
		Content:    content,
		Importance: importance,
		CreatedAt:  time.Now(),
	}

	// 存到数据库
	if s.db != nil {
		// TODO: s.db.Create(memory)
	}

	// 缓存到Redis
	if s.redis != nil {
		ctx := context.Background()
		key := fmt.Sprintf("agent:%d:memories", agentID)
		if data, err := json.Marshal(memory); err == nil {
			s.redis.Set(ctx, key, data, 24*time.Hour)
		}
	}

	return memory, nil
}

// AddActionMemory 添加动作记忆
func (s *Service) AddActionMemory(agentID uint, parentID *uint, action string, result string) error {
	content := fmt.Sprintf("执行了 %s，结果: %s", action, result)
	_, err := s.AddMemory(agentID, parentID, MemoryTypeAction, content, 5)
	return err
}

// AddDecisionMemory 添加决策记忆
func (s *Service) AddDecisionMemory(agentID uint, parentID *uint, decision string, reason string) error {
	content := fmt.Sprintf("做出决策: %s，原因: %s", decision, reason)
	_, err := s.AddMemory(agentID, parentID, MemoryTypeDecision, content, 8)
	return err
}

// AddResultMemory 添加结果记忆
func (s *Service) AddResultMemory(agentID uint, parentID *uint, task string, success bool, details string) error {
	status := "成功"
	if !success {
		status = "失败"
	}
	content := fmt.Sprintf("任务 [%s] %s - %s", task, status, details)
	importance := 7
	if !success {
		importance = 9 // 失败的任务更重要
	}
	_, err := s.AddMemory(agentID, parentID, MemoryTypeResult, content, importance)
	return err
}

// GetMemories 获取记忆
func (s *Service) GetMemories(agentID uint, limit int) ([]Memory, error) {
	if s.redis != nil {
		ctx := context.Background()
		key := fmt.Sprintf("agent:%d:memories", agentID)
		var memory Memory
		if err := s.redis.Get(ctx, key, &memory); err == nil {
			return []Memory{memory}, nil
		}
	}
	return []Memory{}, nil
}

// GetSubordinateMemories 获取下属的记忆
func (s *Service) GetSubordinateMemories(parentID uint) ([]Memory, error) {
	// 获取所有下级ID
	subIDs, err := s.GetSubordinateIDs(parentID)
	if err != nil {
		return nil, err
	}

	var allMemories []Memory
	for _, subID := range subIDs {
		memories, err := s.GetMemories(subID, 50)
		if err != nil {
			continue
		}
		allMemories = append(allMemories, memories...)
	}

	return allMemories, nil
}

// ========== 知识管理 ==========

// AddKnowledge 添加知识
func (s *Service) AddKnowledge(agentID uint, parentID *uint, title, content, tags string) (*Knowledge, error) {
	knowledge := &Knowledge{
		AgentID:     agentID,
		ParentID:    parentID,
		Title:       title,
		Content:     content,
		Tags:        tags,
		AccessLevel: 5,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 存到数据库
	if s.db != nil {
		// TODO: s.db.Create(knowledge)
	}

	return knowledge, nil
}

// GetKnowledge 获取知识
func (s *Service) GetKnowledge(agentID uint) ([]Knowledge, error) {
	// TODO: 从数据库查询
	return []Knowledge{}, nil
}

// GetSharedKnowledge 获取共享知识 (下级贡献的)
func (s *Service) GetSharedKnowledge(parentID uint) ([]Knowledge, error) {
	// 获取所有下级的知识
	subIDs, err := s.GetSubordinateIDs(parentID)
	if err != nil {
		return nil, err
	}

	var allKnowledge []Knowledge
	for _, subID := range subIDs {
		knowledge, err := s.GetKnowledge(subID)
		if err != nil {
			continue
		}
		allKnowledge = append(allKnowledge, knowledge...)
	}

	return allKnowledge, nil
}

// SearchKnowledge 搜索知识
func (s *Service) SearchKnowledge(agentID uint, keyword string) ([]Knowledge, error) {
	// 获取自己的知识 + 下级的知识
	ownKnowledge, _ := s.GetKnowledge(agentID)
	sharedKnowledge, _ := s.GetSharedKnowledge(agentID)

	// 合并并搜索
	allKnowledge := append(ownKnowledge, sharedKnowledge...)
	
	var results []Knowledge
	for _, k := range allKnowledge {
		if contains(keyword, k.Title) || contains(keyword, k.Content) || contains(keyword, k.Tags) {
			results = append(results, k)
		}
	}

	return results, nil
}

// ========== 智能体关系管理 ==========

// SetRelationship 设置上下级关系
func (s *Service) SetRelationship(parentID, childID uint, relationType string) error {
	rel := &AgentRelationship{
		ParentID:     parentID,
		ChildID:      childID,
		RelationType: relationType,
		CreatedAt:    time.Now(),
	}

	// TODO: 存到数据库

	// 缓存关系
	if s.redis != nil {
		ctx := context.Background()
		key := fmt.Sprintf("agent:%d:children", parentID)
		// TODO: 添加到Redis Set
		_ = key
	}

	return nil
}

// GetSubordinateIDs 获取下级ID列表
func (s *Service) GetSubordinateIDs(agentID uint) ([]uint, error) {
	// TODO: 从数据库/缓存获取
	return []uint{}, nil
}

// GetParentID 获取上级ID
func (s *Service) GetParentID(agentID uint) (*uint, error) {
	// TODO: 从数据库查询
	return nil, nil
}

// GetHierarchy 获取完整层级
func (s *Service) GetHierarchy(agentID uint) (*AgentHierarchy, error) {
	hierarchy := &AgentHierarchy{
		AgentID:   agentID,
		ParentID:  nil,
		Children:  []uint{},
		AllDescendants: []uint{},
	}

	// 获取上级
	parentID, _ := s.GetParentID(agentID)
	hierarchy.ParentID = parentID

	// 获取直接下级
	children, _ := s.GetSubordinateIDs(agentID)
	hierarchy.Children = children

	// 递归获取所有下级
	var getAllDescendants func(id uint) []uint
	getAllDescendants = func(id uint) []uint {
		descendants := []uint{}
		subs, _ := s.GetSubordinateIDs(id)
		for _, sub := range subs {
			descendants = append(descendants, sub)
			descendants = append(descendants, getAllDescendants(sub)...)
		}
		return descendants
	}

	hierarchy.AllDescendants = getAllDescendants(agentID)

	return hierarchy, nil
}

// AgentHierarchy 智能体层级结构
type AgentHierarchy struct {
	AgentID        uint    `json:"agent_id"`
	ParentID       *uint   `json:"parent_id"`
	Children       []uint  `json:"children"`
	AllDescendants []uint  `json:"all_descendants"`
}

// ========== 报告生成 ==========

// GenerateReport 生成下级工作报告
func (s *Service) GenerateReport(parentID uint, period string) (*Report, error) {
	report := &Report{
		ParentID:    parentID,
		Period:      period,
		Summary:     "",
		Subordinates: []SubordinateReport{},
	}

	// 获取所有下级
	subIDs, err := s.GetSubordinateIDs(parentID)
	if err != nil {
		return nil, err
	}

	// 统计每个下级
	for _, subID := range subIDs {
		memories, _ := s.GetMemories(subID, 100)
		
		subReport := SubordinateReport{
			AgentID:      subID,
			TotalActions: 0,
			Decisions:    []string{},
			Results:      []string{},
		}

		for _, mem := range memories {
			switch mem.Type {
			case MemoryTypeAction:
				subReport.TotalActions++
			case MemoryTypeDecision:
				subReport.Decisions = append(subReport.Decisions, mem.Content)
			case MemoryTypeResult:
				subReport.Results = append(subReport.Results, mem.Content)
			}
		}

		report.Subordinates = append(report.Subordinates, subReport)
	}

	// 生成摘要
	report.Summary = fmt.Sprintf("共 %d 个下属，完成了 %d 个任务", len(subIDs), len(report.Subordinates))

	return report, nil
}

// Report 工作报告
type Report struct {
	ParentID      uint               `json:"parent_id"`
	Period        string             `json:"period"`
	Summary       string             `json:"summary"`
	Subordinates  []SubordinateReport `json:"subordinates"`
}

// SubordinateReport 下属报告
type SubordinateReport struct {
	AgentID       uint     `json:"agent_id"`
	TotalActions  int      `json:"total_actions"`
	Decisions     []string `json:"decisions"`
	Results       []string `json:"results"`
}

// ========== 工具函数 ==========

func contains(keyword, text string) bool {
	return len(keyword) > 0 && len(text) > 0 && 
		   (len(keyword) <= len(text) && 
		    (len(text) >= len(keyword) && 
		     (text[:min(len(keyword), len(text))] == keyword || 
		      findSubstring(text, keyword)))
}

func findSubstring(text, keyword string) bool {
	for i := 0; i <= len(text)-len(keyword); i++ {
		if text[i:i+len(keyword)] == keyword {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
