package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"agent-flow/internal/api"
	"agent-flow/internal/channel"
	"agent-flow/internal/store"
)

func main() {
	// 初始化存储
	db, err := store.NewPostgres(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()

	// 初始化Redis
	redis, err := store.NewRedis(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}
	defer redis.Close()

	// 初始化渠道管理器
	channelMgr := channel.NewManager(db, redis)

	// 路由设置
	r := gin.Default()

	// API路由
	apiHandler := api.NewHandler(db, redis, channelMgr)
	apiHandler.RegisterRoutes(r)

	// Webhook路由 (各渠道消息入口)
	r.POST("/webhook/:channel_type", channelMgr.HandleWebhook)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
