package channel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"agent-flow/internal/store"
)

// ChannelType 渠道类型
type ChannelType string

const (
	ChannelFeishu   ChannelType = "feishu"
	ChannelTelegram ChannelType = "telegram"
	ChannelDiscord  ChannelType = "discord"
	ChannelWhatsApp ChannelType = "whatsapp"
	ChannelSignal   ChannelType = "signal"
	ChannelWebAPI   ChannelType = "webapi"
)

// Message 统一消息格式
type Message struct {
	Type      string      `json:"type"`       // text/image/file
	Content   string      `json:"content"`    // 文本内容或文件URL
	UserID    string      `json:"user_id"`    // 用户ID
	ChannelID string      `json:"channel_id"` // 渠道ID
	Channel   string      `json:"channel"`    // 渠道类型
	RawData   interface{} `json:"raw_data"`   // 原始数据
}

// Sender 消息发送者接口
type Sender interface {
	SendMessage(msg Message) error
	SendText(userID, text string) error
}

// Manager 渠道管理器
type Manager struct {
	db     *store.Postgres
	redis  *store.Redis
	adapters map[ChannelType]Adapter
}

// NewManager 创建渠道管理器
func NewManager(db *store.Postgres, redis *store.Redis) *Manager {
	m := &Manager{
		db:     db,
		redis:  redis,
		adapters: make(map[ChannelType]Adapter),
	}

	// 注册渠道适配器
	m.adapters[ChannelFeishu] = NewFeishuAdapter()
	m.adapters[ChannelTelegram] = NewTelegramAdapter()
	m.adapters[ChannelDiscord] = NewDiscordAdapter()
	m.adapters[ChannelWhatsApp] = NewWhatsAppAdapter()

	return m
}

// HandleWebhook 处理各渠道webhook
func (m *Manager) HandleWebhook(c *gin.Context) {
	channelType := c.Param("channel_type")

	adapter, ok := m.adapters[ChannelType(channelType)]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported channel type"})
		return
	}

	// 解析消息
	msg, err := adapter.ParseWebhook(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理消息
	response, err := m.processMessage(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 发送回复
	if response != "" {
		if err := adapter.SendMessage(msg.UserID, response); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// processMessage 处理消息（核心逻辑）
func (m *Manager) processMessage(msg *Message) (string, error) {
	// TODO: 
	// 1. 查找或创建会话
	// 2. 获取关联的Agent/Flow
	// 3. 调用AI处理
	// 4. 返回响应

	return "收到消息: " + msg.Content, nil
}

// Adapter 渠道适配器接口
type Adapter interface {
	ParseWebhook(req *http.Request) (*Message, error)
	SendMessage(userID, text string) error
	GetChannelType() ChannelType
}

// ========== 飞书适配器 ==========

type FeishuAdapter struct{}

func NewFeishuAdapter() *FeishuAdapter {
	return &FeishuAdapter{}
}

func (a *FeishuAdapter) GetChannelType() ChannelType {
	return ChannelFeishu
}

type FeishuWebhook struct {
	Type    string `json:"type"`
	Message struct {
		Type string `json:"type"`
		Text string `json:"text"`
		ID   string `json:"id"`
	} `json:"message"`
	Sender struct {
		SenderID struct {
			ID string `json:"id"`
		} `json:"sender_id"`
	} `json:"sender"`
}

func (a *FeishuAdapter) ParseWebhook(req *http.Request) (*Message, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var webhook FeishuWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		return nil, err
	}

	return &Message{
		Type:    webhook.Message.Type,
		Content: webhook.Message.Text,
		UserID:  webhook.Sender.SenderID.ID,
		Channel: string(ChannelFeishu),
		RawData: webhook,
	}, nil
}

func (a *FeishuAdapter) SendMessage(userID, text string) error {
	// TODO: 调用飞书API发送消息
	fmt.Printf("[Feishu] Send to %s: %s\n", userID, text)
	return nil
}

// ========== Telegram适配器 ==========

type TelegramAdapter struct{}

func NewTelegramAdapter() *TelegramAdapter {
	return &TelegramAdapter{}
}

func (a *TelegramAdapter) GetChannelType() ChannelType {
	return ChannelTelegram
}

type TelegramUpdate struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		Text     string `json:"text"`
		Chat     struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		From struct {
			ID int64 `json:"id"`
		} `json:"from"`
	} `json:"message"`
}

func (a *TelegramAdapter) ParseWebhook(req *http.Request) (*Message, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var update TelegramUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		return nil, err
	}

	return &Message{
		Type:    "text",
		Content: update.Message.Text,
		UserID:  fmt.Sprintf("%d", update.Message.From.ID),
		Channel: string(ChannelTelegram),
		RawData: update,
	}, nil
}

func (a *TelegramAdapter) SendMessage(userID, text string) error {
	// TODO: 调用Telegram Bot API发送消息
	fmt.Printf("[Telegram] Send to %s: %s\n", userID, text)
	return nil
}

// ========== Discord适配器 ==========

type DiscordAdapter struct{}

func NewDiscordAdapter() *DiscordAdapter {
	return &DiscordAdapter{}
}

func (a *DiscordAdapter) GetChannelType() ChannelType {
	return ChannelDiscord
}

type DiscordWebhook struct {
	Type        int    `json:"type"`
	ChannelID   string `json:"channel_id"`
	GuildID     string `json:"guild_id"`
	Content     string `json:"content"`
	Member      struct {
		User struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"member"`
}

func (a *DiscordAdapter) ParseWebhook(req *http.Request) (*Message, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var webhook DiscordWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		return nil, err
	}

	return &Message{
		Type:    "text",
		Content: webhook.Content,
		UserID:  webhook.Member.User.ID,
		Channel: string(ChannelDiscord),
		RawData: webhook,
	}, nil
}

func (a *DiscordAdapter) SendMessage(userID, text string) error {
	// TODO: 调用Discord API发送消息
	fmt.Printf("[Discord] Send to %s: %s\n", userID, text)
	return nil
}

// ========== WhatsApp适配器 ==========

type WhatsAppAdapter struct{}

func NewWhatsAppAdapter() *WhatsAppAdapter {
	return &WhatsAppAdapter{}
}

func (a *WhatsAppAdapter) GetChannelType() ChannelType {
	return ChannelWhatsApp
}

type WhatsAppWebhook struct {
	Object string `json:"object"`
	Entry  []struct {
		Changes []struct {
			Value struct {
				Messages []struct {
					From string `json:"from"`
					Text  struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

func (a *WhatsAppAdapter) ParseWebhook(req *http.Request) (*Message, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var webhook WhatsAppWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		return nil, err
	}

	if len(webhook.Entry) == 0 || len(webhook.Entry[0].Changes) == 0 ||
		len(webhook.Entry[0].Changes[0].Value.Messages) == 0 {
		return nil, fmt.Errorf("no message found")
	}

	msg := webhook.Entry[0].Changes[0].Value.Messages[0]

	return &Message{
		Type:    "text",
		Content: msg.Text.Body,
		UserID:  msg.From,
		Channel: string(ChannelWhatsApp),
		RawData: webhook,
	}, nil
}

func (a *WhatsAppAdapter) SendMessage(userID, text string) error {
	// TODO: 调用WhatsApp Business API发送消息
	fmt.Printf("[WhatsApp] Send to %s: %s\n", userID, text)
	return nil
}
