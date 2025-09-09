package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 設置日誌
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Go 簡化版應用程序啟動")

	// 創建 Gin 路由器
	r := gin.Default()

	// 中間件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 基本路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "歡迎來到 Go 簡化版服務器！",
			"status":  "running",
			"version": "1.0.2",
			"features": []string{
				"無資料庫版本",
				"CI/CD 測試",
				"基本 API",
			},
		})
	})

	// 健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"service":  "go-simple-app",
			"database": "none",
			"uptime":   "running",
		})
	})

	// API 路由
	api := r.Group("/api")
	{
		api.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"message": "API 正常運行",
			})
		})

		api.GET("/info", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"app_name": "Go Simple App",
				"version": "1.0.0",
				"environment": "production",
				"features": []string{
					"REST API",
					"健康檢查",
					"日誌記錄",
				},
			})
		})
	}

	// 獲取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 啟動服務器
	logrus.Info("服務器準備啟動", logrus.Fields{
		"port": port,
		"mode": gin.Mode(),
	})

	log.Printf("服務器啟動在 :%s", port)
	log.Fatal(r.Run(":" + port))
}
