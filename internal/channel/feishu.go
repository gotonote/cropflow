package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// FeishuConfig 飞书配置
type FeishuConfig struct {
	AppID        string `json:"app_id"`
	AppSecret    string `json:"app_secret"`
	Verification string `json:"verification"` // 验证Token
	EncryptKey   string `json:"encrypt_key"`   // 加密Key
}

// FeishuAdapter 飞书适配器
type FeishuAdapter struct {
	config *FeishuConfig
	client *http.Client
	appID  string
	secret string
}

// NewFeishuAdapter 创建飞书适配器
func NewFeishuAdapter() *FeishuAdapter {
	return &FeishuAdapter{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Init 初始化 (从DB加载配置)
func (a *FeishuAdapter) Init(config string) error {
	var cfg FeishuConfig
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return err
	}
	a.config = &cfg
	a.appID = cfg.AppID
	a.secret = cfg.AppSecret
	return nil
}

// GetChannelType 获取渠道类型
func (a *FeishuAdapter) GetChannelType() ChannelType {
	return ChannelFeishu
}

// ========== Webhook处理 ==========

// ParseWebhook 解析webhook请求
func (a *FeishuAdapter) ParseWebhook(req *http.Request) (*Message, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// 飞书v2消息格式
	var payload struct {
		Type    int `json:"type"`
		Schema  string `json:"schema"`
		Header  struct {
			EventID   string `json:"event_id"`
			EventType string `json:"event_type"`
			CreateTime string `json:"create_time"`
			Token     string `json:"token"`
			AppID     string `json:"app_id"`
		} `json:"header"`
		Event struct {
			Message struct {
				MessageID string `json:"message_id"`
				RootID    string `json:"root_id"`
				ParentID  string `json:"parent_id"`
				CreateTime string `json:"create_time"`
				ChatID    string `json:"chat_id"`
				ChatType  string `json:"chat_type"`
				MessageType string `json:"msg_type"`
				Content   string `json:"content"`
				Mentions  []struct {
					ID      string `json:"id"`
					IDType  string `json:"id_type"`
					Name    string `json:"name"`
					Key     string `json:"key"`
				} `json:"mentions"`
				Sender struct {
					SenderID struct {
						UnionID string `json:"union_id"`
						OpenID  string `json:"open_id"`
						UserID  string `json:"user_id"`
					} `json:"sender_id"`
					SenderType string `json:"sender_type"`
				} `json:"sender"`
			} `json:"message"`
		} `json:"event"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	// 处理验证请求 (URL验证)
	if payload.Header.EventType == "url_verification" {
		return &Message{
			Type:    "url_verification",
			Content: a.config.Verification,
			Channel: string(ChannelFeishu),
			RawData: payload,
		}, nil
	}

	// 处理事件回调验证
	if payload.Header.EventType == "callback_verification" {
		return nil, nil
	}

	// 解析消息内容
	content := payload.Event.Message.Content
	var msgContent string
	
	// 文本消息
	if payload.Event.Message.MessageType == "text" {
		var textContent struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal([]byte(content), &textContent); err == nil {
			msgContent = textContent.Text
		}
	}

	return &Message{
		Type:      payload.Event.Message.MessageType,
		Content:   msgContent,
		UserID:    payload.Event.Message.Sender.SenderID.OpenID,
		ChannelID: payload.Event.Message.ChatID,
		Channel:   string(ChannelFeishu),
		RawData:   payload,
	}, nil
}

// SendMessage 发送消息
func (a *FeishuAdapter) SendMessage(recipient, text string) error {
	// 获取tenant_access_token
	token, err := a.getTenantAccessToken()
	if err != nil {
		return err
	}

	// 发送消息API
	url := "https://open.feishu.cn/open-apis/im/v1/messages"
	
	// 根据recipient类型决定发送方式
	reqBody := map[string]interface{}{
		"receive_id": recipient,
		"msg_type":   "text",
		"content": map[string]string{
			"text": text,
		},
	}
	
	body, _ := json.Marshal(reqBody)
	
	httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	
	resp, err := a.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("feishu API error: %s", string(respBody))
	}

	return nil
}

// SendMessageWithCard 使用卡片消息
func (a *FeishuAdapter) SendMessageWithCard(recipient, title, content string) error {
	token, err := a.getTenantAccessToken()
	if err != nil {
		return err
	}

	url := "https://open.feishu.cn/open-apis/im/v1/messages"

	// 卡片消息模板
	card := map[string]interface{}{
		"config": map[string]interface{}{
			"wide_screen_mode": true,
		},
		"header": map[string]interface{}{
			"title": map[string]string{
				"tag":  "plain_text",
				"content": title,
			},
			"template": "blue",
		},
		"elements": []map[string]interface{}{
			{
				"tag":  "markdown",
				"content": content,
			},
		},
	}

	reqBody := map[string]interface{}{
		"receive_id": recipient,
		"msg_type":   "interactive",
		"content": map[string]string{
			"card": fmt.Sprintf(`%s`, card), // 实际应JSON序列化
		},
	}

	body, _ := json.Marshal(reqBody)

	httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ========== 认证相关 ==========

// getTenantAccessToken 获取tenant_access_token
func (a *FeishuAdapter) getTenantAccessToken() (string, error) {
	// TODO: 缓存token
	url := "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	
	reqBody := map[string]string{
		"app_id":     a.appID,
		"app_secret": a.secret,
	}
	
	body, _ := json.Marshal(reqBody)
	
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Code              int    `json:"code"`
		Msg               string `json:"msg"`
		TenantAccessToken string `json:"tenant_access_token"`
		Expire            int    `json:"expire"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Code != 0 {
		return "", fmt.Errorf("feishu auth error: %s", result.Msg)
	}

	return result.TenantAccessToken, nil
}

// SendText 发送文本消息 (实现Sender接口)
func (a *FeishuAdapter) SendText(userID, text string) error {
	return a.SendMessage(userID, text)
}

// ========== 飞书Bot消息处理 ==========

// FeishuHandler 飞书消息处理器
type FeishuHandler struct {
	adapter    *FeishuAdapter
	callbackURL string
}

// NewFeishuHandler 创建飞书处理器
func NewFeishuHandler(appID, appSecret, callbackURL string) *FeishuHandler {
	adapter := NewFeishuAdapter()
	adapter.appID = appID
	adapter.secret = appSecret

	return &FeishuAdapter{
		adapter:    adapter,
		callbackURL: callbackURL,
	}
}

// RegisterEvents 注册事件回调
func (h *FeishuHandler) RegisterEvents(appID, appSecret string) error {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v4/hook/%s", appID)
	
	// 设置回调URL
	reqBody := map[string]string{
		"url": h.callbackURL,
	}
	
	body, _ := json.Marshal(reqBody)
	
	httpReq, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	
	// 需要签名验证
	sign := h.generateSign(appSecret)
	httpReq.Header.Set("Authorization", "Bearer "+sign)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// generateSign 生成签名
func (h *FeishuHandler) generateSign(secret string) string {
	// TODO: 实现飞书签名算法
	// timestamp + "\n" + secret  -> HmacSHA256 -> Base64
	return ""
}

// ========== 接收消息处理 ==========

// FeishuMessageHandler 飞书消息处理函数类型
type FeishuMessageHandler func(userID, messageID, messageType, content string) string

// HandleTextMessage 处理文本消息
func HandleTextMessage(userID, messageID, content string) string {
	// 这里会调用Agent/Flow处理
	// TODO: 集成Agent服务
	return "收到消息: " + content
}

// ========== 工具函数 ==========

// SendMessageToChat 发送消息到群聊
func (a *FeishuAdapter) SendMessageToChat(chatID, text string) error {
	token, err := a.getTenantAccessToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=chat_id")

	reqBody := map[string]interface{}{
		"receive_id": chatID,
		"msg_type":   "text",
		"content": map[string]string{
			"text": text,
		},
	}

	body, _ := json.Marshal(reqBody)

	httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ReplyMessage 回复消息
func (a *FeishuAdapter) ReplyMessage(messageID, text string) error {
	token, err := a.getTenantAccessToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages/%s/reply", messageID)

	reqBody := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": text,
		},
	}

	body, _ := json.Marshal(reqBody)

	httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetUserInfo 获取用户信息
func (a *FeishuAdapter) GetUserInfo(openID string) (map[string]interface{}, error) {
	token, err := a.getTenantAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/auth/v3/user_info/get?open_id=%s", openID)

	httpReq, _ := http.NewRequest("GET", url, nil)
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}

// 环境变量获取飞书配置
func GetFeishuConfigFromEnv() (appID, appSecret string) {
	appID = os.Getenv("FEISHU_APP_ID")
	appSecret = os.Getenv("FEISHU_APP_SECRET")
	return
}

// 兼容性问题修复
var _ Adapter = (*FeishuAdapter)(nil)

func init() {
	// 确保实现了接口
	var _ Adapter = (*FeishuAdapter)(nil)
}

// gin绑定 (修复类型错误)
func (a *FeishuAdapter) BindRoutes(r *gin.Engine) {
	r.POST("/webhook/feishu", func(c *gin.Context) {
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
