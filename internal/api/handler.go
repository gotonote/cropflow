package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"agent-flow/internal/store"
	"agent-flow/internal/channel"
)

type Handler struct {
	db          *store.Postgres
	redis       *store.Redis
	channelMgr  *channel.Manager
}

func NewHandler(db *store.Postgres, redis *store.Redis, channelMgr *channel.Manager) *Handler {
	return &Handler{
		db:         db,
		redis:      redis,
		channelMgr: channelMgr,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// 智能体管理
		agents := api.Group("/agents")
		{
			agents.GET("", h.ListAgents)
			agents.POST("", h.CreateAgent)
			agents.PUT("/:id", h.UpdateAgent)
			agents.DELETE("/:id", h.DeleteAgent)
		}

		// 流程管理
		flows := api.Group("/flows")
		{
			flows.GET("", h.ListFlows)
			flows.POST("", h.CreateFlow)
			flows.PUT("/:id", h.UpdateFlow)
			flows.DELETE("/:id", h.DeleteFlow)
			flows.POST("/:id/execute", h.ExecuteFlow)
		}

		// 渠道管理
		channels := api.Group("/channels")
		{
			channels.GET("", h.ListChannels)
			channels.POST("", h.CreateChannel)
			channels.PUT("/:id", h.UpdateChannel)
			channels.DELETE("/:id", h.DeleteChannel)
		}

		// 会话管理
		conversations := api.Group("/conversations")
		{
			conversations.GET("", h.ListConversations)
		}
	}
}

// ========== Agent APIs ==========

type CreateAgentRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	ModelProvider string `json:"model_provider"`
	ModelName    string `json:"model_name"`
	ModelConfig  string `json:"model_config"`
	Tools        string `json:"tools"`
}

func (h *Handler) ListAgents(c *gin.Context) {
	agents, err := h.db.ListAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, agents)
}

func (h *Handler) CreateAgent(c *gin.Context) {
	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent := &store.Agent{
		Name:          req.Name,
		Description:   req.Description,
		ModelProvider: req.ModelProvider,
		ModelName:     req.ModelName,
		ModelConfig:   req.ModelConfig,
		Tools:         req.Tools,
	}

	if err := h.db.CreateAgent(agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

func (h *Handler) UpdateAgent(c *gin.Context) {
	id := c.Param("id")
	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, err := h.db.GetAgent(parseUint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	agent.Name = req.Name
	agent.Description = req.Description
	agent.ModelProvider = req.ModelProvider
	agent.ModelName = req.ModelName
	agent.ModelConfig = req.ModelConfig
	agent.Tools = req.Tools

	if err := h.db.UpdateAgent(agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}

func (h *Handler) DeleteAgent(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.DeleteAgent(parseUint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ========== Flow APIs ==========

type CreateFlowRequest struct {
	Name        string `json:"name" binding:"required"`
	Nodes       string `json:"nodes"`
	Edges       string `json:"edges"`
	TriggerType string `json:"trigger_type"`
}

func (h *Handler) ListFlows(c *gin.Context) {
	flows, err := h.db.ListFlows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flows)
}

func (h *Handler) CreateFlow(c *gin.Context) {
	var req CreateFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flow := &store.Flow{
		Name:        req.Name,
		Nodes:       req.Nodes,
		Edges:       req.Edges,
		TriggerType: req.TriggerType,
		Enabled:     true,
	}

	if err := h.db.CreateFlow(flow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, flow)
}

func (h *Handler) UpdateFlow(c *gin.Context) {
	id := c.Param("id")
	var req CreateFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flow, err := h.db.GetFlow(parseUint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "flow not found"})
		return
	}

	flow.Name = req.Name
	flow.Nodes = req.Nodes
	flow.Edges = req.Edges
	flow.TriggerType = req.TriggerType

	if err := h.db.UpdateFlow(flow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flow)
}

func (h *Handler) DeleteFlow(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.DeleteFlow(parseUint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) ExecuteFlow(c *gin.Context) {
	// TODO: 实现流程执行逻辑
	c.JSON(http.StatusOK, gin.H{"message": "flow execution not implemented yet"})
}

// ========== Channel APIs ==========

type CreateChannelRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Config  string `json:"config"`
	Enabled bool   `json:"enabled"`
}

func (h *Handler) ListChannels(c *gin.Context) {
	channels, err := h.db.ListChannels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, channels)
}

func (h *Handler) CreateChannel(c *gin.Context) {
	var req CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ch := &store.Channel{
		Name:    req.Name,
		Type:    req.Type,
		Config:  req.Config,
		Enabled: req.Enabled,
	}

	if err := h.db.CreateChannel(ch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ch)
}

func (h *Handler) UpdateChannel(c *gin.Context) {
	id := c.Param("id")
	var req CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ch, err := h.db.GetChannel(parseUint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
		return
	}

	ch.Name = req.Name
	ch.Type = req.Type
	ch.Config = req.Config
	ch.Enabled = req.Enabled

	if err := h.db.UpdateChannel(ch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ch)
}

func (h *Handler) DeleteChannel(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.DeleteChannel(parseUint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ========== Conversation APIs ==========

func (h *Handler) ListConversations(c *gin.Context) {
	// TODO: 实现会话列表查询
	c.JSON(http.StatusOK, []interface{}{})
}

// ========== 工具函数 ==========

func parseUint(s string) uint {
	var n uint
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + uint(c-'0')
		}
	}
	return n
}
