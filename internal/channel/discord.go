package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// DiscordConfig Discord配置
type DiscordConfig struct {
	BotToken   string `json:"bot_token"`
	AppID      string `json:"app_id"`
	PublicKey  string `json:"public_key"`
	GuildID    string `json:"guild_id"`
}

// DiscordAdapter Discord适配器
type DiscordAdapter struct {
	config   *DiscordConfig
	client   *http.Client
	botToken string
}

// NewDiscordAdapter 创建Discord适配器
func NewDiscordAdapter() *DiscordAdapter {
	return &DiscordAdapter{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Init 初始化
func (a *DiscordAdapter) Init(config string) error {
	var cfg DiscordConfig
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return err
	}
	a.config = &cfg
	a.botToken = cfg.BotToken
	return nil
}

// GetChannelType 获取渠道类型
func (a *DiscordAdapter) GetChannelType() ChannelType {
	return ChannelDiscord
}

// ========== Webhook处理 ==========

// ParseWebhook 解析webhook请求
func (a *DiscordAdapter) ParseWebhook(req *http.Request) (*Message, error) {
	// Discord使用Interaction机制
	contentType := req.Header.Get("Content-Type")
	
	if strings.Contains(contentType, "application/json") {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, err
		}

		// 处理PING (Discord交互验证)
		t, _ := payload["type"].(float64)
		if t == 1 {
			// 返回PONG
			return &Message{
				Type:    "ping",
				Content: "",
				Channel: string(ChannelDiscord),
				RawData: payload,
			}, nil
		}

		// 处理消息组件/命令
		if t == 3 { // MESSAGE_COMPONENT
			data, _ := payload["data"].(map[string]interface{})
			customID, _ := data["custom_id"].(string)
			
			return &Message{
				Type:      "component",
				Content:   customID,
				Channel:   string(ChannelDiscord),
				RawData:   payload,
			}, nil
		}

		// 处理Slash命令
		if t == 2 { // APPLICATION_COMMAND
			data, _ := payload["data"].(map[string]interface{})
			name, _ := data["name"].(string)
			
			options, _ := data["options"].([]interface{})
			var args string
			if len(options) > 0 {
				opt := options[0].(map[string]interface{})
				value, _ := opt["value"].(string)
				args = value
			}

			member, _ := payload["member"].(map[string]interface{})
			user, _ := member["user"].(map[string]interface{})
			userID, _ := user["id"].(string)

			guildID, _ := payload["guild_id"].(string)

			return &Message{
				Type:      "command",
				Content:   "/" + name + " " + args,
				UserID:    userID,
				ChannelID: guildID,
				Channel:   string(ChannelDiscord),
				RawData:   payload,
			}, nil
		}
	}

	// 处理普通消息 (通过Messages API)
	return &Message{
		Type:    "unknown",
		Content: "",
		Channel: string(ChannelDiscord),
	}, nil
}

// SendMessage 发送消息
func (a *DiscordAdapter) SendMessage(recipient, text string) error {
	// recipient 可以是 channel_id
	return a.sendMessageToChannel(recipient, text, nil)
}

// SendText 发送文本 (实现Sender接口)
func (a *DiscordAdapter) SendText(userID, text string) error {
	// Discord需要通过Channel发送消息
	// TODO: 获取用户所在的DM Channel
	return nil
}

// sendMessageToChannel 发送消息到频道
func (a *DiscordAdapter) sendMessageToChannel(channelID, content string, embed *DiscordEmbed) error {
	apiURL := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)

	bodyData := map[string]string{
		"content": content,
	}

	if embed != nil {
		bodyData["embeds"] = "[]" // TODO: 处理embed
	}

	body, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bot "+a.botToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord API error: %s", string(respBody))
	}

	return nil
}

// SendEmbed 发送嵌入消息
func (a *DiscordAdapter) SendEmbed(channelID string, embed DiscordEmbed) error {
	apiURL := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)

	body, _ := json.Marshal(map[string]interface{}{
		"embeds": []DiscordEmbed{embed},
	})

	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bot "+a.botToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DiscordEmbed Discord嵌入
type DiscordEmbed struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Color       int                     `json:"color"`
	Fields      []DiscordEmbedField     `json:"fields,omitempty"`
	Footer      *DiscordEmbedFooter    `json:"footer,omitempty"`
	Thumbnail   *DiscordEmbedThumbnail `json:"thumbnail,omitempty"`
	Image       *DiscordEmbedImage     `json:"image,omitempty"`
	Author      *DiscordEmbedAuthor    `json:"author,omitempty"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type DiscordEmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

type DiscordEmbedThumbnail struct {
	URL string `json:"url"`
}

type DiscordEmbedImage struct {
	URL string `json:"url"`
}

type DiscordEmbedAuthor struct {
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
	URL    string `json:"url"`
}

// SendComponent 发送组件消息
func (a *DiscordAdapter) SendComponent(channelID string, content string, components []DiscordComponent) error {
	apiURL := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)

	body, _ := json.Marshal(map[string]interface{}{
		"content":    content,
		"components": components,
	})

	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bot "+a.botToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DiscordComponent Discord消息组件
type DiscordComponent struct {
	Type     int                      `json:"type"` // 1=action_row, 2=button, 3=select
	Style    int                      `json:"style,omitempty"`
	Label    string                   `json:"label,omitempty"`
	CustomID string                   `json:"custom_id,omitempty"`
	URL      string                   `json:"url,omitempty"`
	Disabled bool                     `json:"disabled,omitempty"`
	Options  []DiscordSelectOption    `json:"options,omitempty"`
}

type DiscordSelectOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

// CreateResponse 创建交互响应
func (a *DiscordAdapter) CreateResponse(interactionID, interactionToken, responseType int, data map[string]interface{}) error {
	apiURL := fmt.Sprintf("https://discord.com/api/v10/interactions/%d/%s/callback", 
		interactionID, interactionToken)

	body, _ := json.Marshal(map[string]interface{}{
		"type": responseType, // 4 = ChannelMessageWithSource
		"data": data,
	})

	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bot "+a.botToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ========== Slash命令注册 ==========

// RegisterCommands 注册Slash命令
func (a *DiscordAdapter) RegisterCommands(commands []DiscordCommand) error {
	apiURL := fmt.Sprintf("https://discord.com/api/v10/applications/%s/commands", a.config.AppID)

	for _, cmd := range commands {
		body, _ := json.Marshal(cmd)
		
		req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bot "+a.botToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			respBody, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to register command %s: %s", cmd.Name, string(respBody))
		}
	}

	return nil
}

// DiscordCommand Slash命令定义
type DiscordCommand struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Options     []DiscordCommandOption    `json:"options,omitempty"`
}

type DiscordCommandOption struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Type        int                      `json:"type"` // 1=sub_command, 2=sub_command_group, 3=string, 4=integer, 5=boolean, 6=user, 7=channel, 8=role, 9=mentionable, 10=number
	Required    bool                     `json:"required"`
	Choices     []DiscordCommandChoice   `json:"choices,omitempty"`
}

type DiscordCommandChoice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ========== 便捷方法 ==========

// NewDiscordHandler 创建Discord处理器
func NewDiscordHandler(botToken, appID string) *DiscordAdapter {
	adapter := NewDiscordAdapter()
	adapter.botToken = botToken
	adapter.config = &DiscordConfig{
		BotToken: botToken,
		AppID:   appID,
	}
	return adapter
}

// BindRoutes 绑定路由
func (a *DiscordAdapter) BindRoutes(r *gin.Engine, path string) {
	r.POST(path, func(c *gin.Context) {
		msg, err := a.ParseWebhook(c.Request)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		
		if msg.Type == "ping" {
			c.JSON(200, gin.H{"type": 1}) // PONG
			return
		}

		// TODO: 处理消息
		
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// 环境变量获取配置
func GetDiscordConfigFromEnv() (botToken, appID string) {
	botToken = os.Getenv("DISCORD_BOT_TOKEN")
	appID = os.Getenv("DISCORD_APP_ID")
	return
}

// 兼容性问题修复
var _ Adapter = (*DiscordAdapter)(nil)

func init() {
	var _ Adapter = (*DiscordAdapter)(nil)
}
