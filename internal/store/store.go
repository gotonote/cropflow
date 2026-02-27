package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	db.AutoMigrate(
		&Agent{},
		&Flow{},
		&Channel{},
		&Conversation{},
	)

	return &Postgres{db: db}, nil
}

func (p *Postgres) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Agent 智能体
type Agent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	ModelProvider string   `gorm:"size:50" json:"model_provider"`
	ModelName   string    `gorm:"size:100" json:"model_name"`
	ModelConfig string    `gorm:"type:jsonb" json:"model_config"` // JSON存储
	Tools       string    `gorm:"type:jsonb" json:"tools"`         // JSON存储
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Agent) TableName() string {
	return "agents"
}

// Flow 流程
type Flow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Nodes       string    `gorm:"type:jsonb" json:"nodes"`        // React Flow nodes
	Edges       string    `gorm:"type:jsonb" json:"edges"`        // React Flow edges
	TriggerType string    `gorm:"size:50" json:"trigger_type"`    // manual/webhook/schedule
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Flow) TableName() string {
	return "flows"
}

// Channel 渠道
type Channel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Type      string    `gorm:"size:50;not null" json:"type"` // feishu/telegram/discord/whatsapp
	Config    string    `gorm:"type:jsonb" json:"config"`      // 渠道配置(JSON)
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Channel) TableName() string {
	return "channels"
}

// Conversation 会话
type Conversation struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ChannelID   uint      `gorm:"index" json:"channel_id"`
	ChannelType string    `gorm:"size:50" json:"channel_type"`
	UserID      string    `gorm:"size:255;index" json:"user_id"`
	AgentID     uint      `gorm:"index" json:"agent_id"`
	Messages    string    `gorm:"type:jsonb" json:"messages"`     // 消息历史
	Context     string    `gorm:"type:jsonb" json:"context"`      // 额外上下文
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Conversation) TableName() string {
	return "conversations"
}

// Redis
type Redis struct {
	client *redis.Client
}

func NewRedis(addr string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{client: client}, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *Redis) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// 便捷方法
func (p *Postgres) CreateAgent(agent *Agent) error {
	return p.db.Create(agent).Error
}

func (p *Postgres) GetAgent(id uint) (*Agent, error) {
	var agent Agent
	err := p.db.First(&agent, id).Error
	return &agent, err
}

func (p *Postgres) ListAgents() ([]Agent, error) {
	var agents []Agent
	err := p.db.Find(&agents).Error
	return agents, err
}

func (p *Postgres) UpdateAgent(agent *Agent) error {
	return p.db.Save(agent).Error
}

func (p *Postgres) DeleteAgent(id uint) error {
	return p.db.Delete(&Agent{}, id).Error
}

func (p *Postgres) CreateFlow(flow *Flow) error {
	return p.db.Create(flow).Error
}

func (p *Postgres) GetFlow(id uint) (*Flow, error) {
	var flow Flow
	err := p.db.First(&flow, id).Error
	return &flow, err
}

func (p *Postgres) ListFlows() ([]Flow, error) {
	var flows []Flow
	err := p.db.Find(&flows).Error
	return flows, err
}

func (p *Postgres) UpdateFlow(flow *Flow) error {
	return p.db.Save(flow).Error
}

func (p *Postgres) DeleteFlow(id uint) error {
	return p.db.Delete(&Flow{}, id).Error
}

func (p *Postgres) CreateChannel(channel *Channel) error {
	return p.db.Create(channel).Error
}

func (p *Postgres) GetChannel(id uint) (*Channel, error) {
	var ch Channel
	err := p.db.First(&ch, id).Error
	return &ch, err
}

func (p *Postgres) ListChannels() ([]Channel, error) {
	var channels []Channel
	err := p.db.Find(&channels).Error
	return channels, err
}

func (p *Postgres) UpdateChannel(channel *Channel) error {
	return p.db.Save(channel).Error
}

func (p *Postgres) DeleteChannel(id uint) error {
	return p.db.Delete(&Channel{}, id).Error
}

func (p *Postgres) CreateConversation(conv *Conversation) error {
	return p.db.Create(conv).Error
}

func (p *Postgres) GetConversation(channelID uint, userID string) (*Conversation, error) {
	var conv Conversation
	err := p.db.Where("channel_id = ? AND user_id = ?", channelID, userID).
		Order("created_at DESC").First(&conv).Error
	return &conv, err
}

func (p *Postgres) UpdateConversation(conv *Conversation) error {
	return p.db.Save(conv).Error
}

// 格式化DSN
func FormatDSN(host, port, user, password, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
