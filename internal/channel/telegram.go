package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TelegramConfig Telegram配置
type TelegramConfig struct {
	BotToken string `json:"bot_token"`
	WebhookURL string `json:"webhook_url"`
}

// TelegramAdapter Telegram适配器
type TelegramAdapter struct {
	config  *TelegramConfig
	client  *http.Client
	botToken string
}

// NewTelegramAdapter 创建Telegram适配器
func NewTelegramAdapter() *TelegramAdapter {
	return &TelegramAdapter{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Init 初始化
func (a *TelegramAdapter) Init(config string) error {
	var cfg TelegramConfig
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return err
	}
	a.config = &cfg
	a.botToken = cfg.BotToken
	return nil
}

// GetChannelType 获取渠道类型
func (a *TelegramAdapter) GetChannelType() ChannelType {
	return ChannelTelegram
}

// ========== Webhook处理 ==========

// ParseWebhook 解析webhook请求
func (a *TelegramAdapter) ParseWebhook(req *http.Request) (*Message, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var update struct {
		UpdateID int64 `json:"update_id"`
		Message  struct {
			MessageID int    `json:"message_id"`
			Text      string `json:"text"`
			Chat      struct {
				ID   int64  `json:"id"`
				Type string `json:"type"`
			} `json:"chat"`
			From struct {
				ID        int64  `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Date     int    `json:"date"`
			Entities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
			} `json:"entities"`
		} `json:"message"`
		CallbackQuery *struct {
			ID   string `json:"id"`
			Data string `json:"data"`
			From struct {
				ID        int64  `json:"id"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Message *struct {
				MessageID int `json:"message_id"`
				Chat      struct {
					ID int64 `json:"id"`
				} `json:"chat"`
			} `json:"message"`
		} `json:"callback_query"`
	}

	if err := json.Unmarshal(body, &update); err != nil {
		return nil, err
	}

	// 处理普通消息
	if update.Message.Text != "" {
		// 处理命令
		text := update.Message.Text
		var msgType string
		var content string

		if strings.HasPrefix(text, "/") {
			msgType = "command"
			content = text
		} else {
			msgType = "text"
			content = text
		}

		return &Message{
			Type:      msgType,
			Content:   content,
			UserID:    strconv.FormatInt(update.Message.From.ID, 10),
			ChannelID: strconv.FormatInt(update.Message.Chat.ID, 10),
			Channel:   string(ChannelTelegram),
			RawData:   update,
		}, nil
	}

	// 处理回调查询 (Inline Keyboard)
	if update.CallbackQuery != nil {
		return &Message{
			Type:      "callback",
			Content:   update.CallbackQuery.Data,
			UserID:    strconv.FormatInt(update.CallbackQuery.From.ID, 10),
			ChannelID: strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10),
			Channel:   string(ChannelTelegram),
			RawData:   update,
		}, nil
	}

	return nil, nil
}

// SendMessage 发送消息
func (a *TelegramAdapter) SendMessage(recipient, text string) error {
	return a.sendMessageWithOptions(recipient, text, nil)
}

// SendText 发送文本 (实现Sender接口)
func (a *TelegramAdapter) SendText(userID, text string) error {
	return a.sendMessageWithOptions(userID, text, nil)
}

// sendMessageWithOptions 发送带选项的消息
func (a *TelegramAdapter) sendMessageWithOptions(chatID, text string, opts *SendMessageOptions) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", a.botToken)

	reqBody := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
		"parse_mode": "Markdown",
	}

	if opts != nil {
		if opts.ReplyMarkup != nil {
			reqBody["reply_markup"] = opts.ReplyMarkup
		}
		if opts.InlineKeyboard {
			reqBody["reply_markup"] = map[string]interface{}{
				"inline_keyboard": opts.InlineKeyboardButtons,
			}
		}
	}

	body, _ := json.Marshal(reqBody)
	
	resp, err := a.client.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.OK {
		return fmt.Errorf("telegram API error: %s", result.Description)
	}

	return nil
}

// SendMessageOptions 消息选项
type SendMessageOptions struct {
	ReplyKeyboard         [][]KeyboardButton
	InlineKeyboard        bool
	InlineKeyboardButtons [][]InlineKeyboardButton
	ReplyMarkup           map[string]interface{}
}

// KeyboardButton 键盘按钮
type KeyboardButton struct {
	Text string `json:"text"`
}

// InlineKeyboardButton Inline键盘按钮
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

// SendInlineKeyboard 发送Inline键盘
func (a *TelegramAdapter) SendInlineKeyboard(chatID string, buttons [][]InlineKeyboardButton, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", a.botToken)

	reqBody := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
		"reply_markup": map[string]interface{}{
			"inline_keyboard": buttons,
		},
	}

	body, _ := json.Marshal(reqBody)
	
	resp, err := a.client.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// SendPhoto 发送图片
func (a *TelegramAdapter) SendPhoto(chatID, photoURL, caption string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", a.botToken)

	reqBody := map[string]interface{}{
		"chat_id": chatID,
		"photo":   photoURL,
		"caption": caption,
	}

	body, _ := json.Marshal(reqBody)
	
	resp, err := a.client.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// EditMessageText 编辑消息
func (a *TelegramAdapter) EditMessageText(chatID, messageID, text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", a.botToken)

	reqBody := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
	}

	body, _ := json.Marshal(reqBody)
	
	resp, err := a.client.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// AnswerCallbackQuery 回答回调查询
func (a *TelegramAdapter) AnswerCallbackQuery(callbackID, text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", a.botToken)

	reqBody := map[string]interface{}{
		"callback_query_id": callbackID,
		"text":              text,
	}

	body, _ := json.Marshal(reqBody)
	
	resp, err := a.client.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ========== Bot API 方法 ==========

// SetWebhook 设置Webhook
func (a *TelegramAdapter) SetWebhook(webhookURL string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", a.botToken)

	reqBody := map[string]interface{}{
		"url": webhookURL,
	}

	body, _ := json.Marshal(reqBody)
	
	resp, err := a.client.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteWebhook 删除Webhook
func (a *TelegramAdapter) DeleteWebhook() error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/deleteWebhook", a.botToken)
	
	resp, err := a.client.Post(apiURL, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetMe 获取Bot信息
func (a *TelegramAdapter) GetMe() (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", a.botToken)
	
	resp, err := a.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		OK   bool                   `json:"ok"`
		Result map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

// GetUpdates 获取更新
func (a *TelegramAdapter) GetUpdates(offset int64, limit int) ([]map[string]interface{}, error) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&limit=%d", 
		a.botToken, offset, limit)
	
	resp, err := a.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		OK   bool                      `json:"ok"`
		Result []map[string]interface{} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

// ========== 便捷方法 ==========

// NewTelegramHandler 创建Telegram处理器
func NewTelegramHandler(botToken string) *TelegramAdapter {
	adapter := NewTelegramAdapter()
	adapter.botToken = botToken
	return adapter
}

// BindRoutes 绑定路由
func (a *TelegramAdapter) BindRoutes(r *gin.Engine, path string) {
	r.POST(path, func(c *gin.Context) {
		msg, err := a.ParseWebhook(c.Request)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if msg == nil {
			c.JSON(200, gin.H{"status": "ok"})
			return
		}

		// TODO: 处理消息
		
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// 环境变量获取配置
func GetTelegramConfigFromEnv() (botToken string) {
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	return
}

// 兼容性问题修复
var _ Adapter = (*TelegramAdapter)(nil)

func init() {
	var _ Adapter = (*TelegramAdapter)(nil)
}
