package main

import (
	"fmt"
	"log"
	"os"
	"tems-web-api/config"
	"tems-web-api/models"
	"tems-web-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 初始化全局配置
	cfg, err := config.InitLoadGlobalConfig()
	if err != nil {
		log.Fatal("Failed to load global config:", err)
	}

	// 初始化数据库连接
	if err := config.InitDB(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表结构
	if err := models.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化默认管理员用户
	if err := models.InitManagerUser(); err != nil {
		log.Fatal("Failed to init manager user:", err)
	}

	// 设置Gin模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 获取服务器配置
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server starting on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
