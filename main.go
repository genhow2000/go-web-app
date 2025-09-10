package main

import (
	"log"
	"os"
	"go-simple-app/config"
	"go-simple-app/controllers"
	"go-simple-app/database"
	"go-simple-app/logger"
	"go-simple-app/models"
	"go-simple-app/routes"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日誌系統
	logger.Info("Go應用程序啟動", logrus.Fields{
		"version": "2.0.0",
		"service": "go-simple-app",
	})

	// 載入配置
	cfg := config.Load()
	logger.Info("配置載入完成", logrus.Fields{
		"server_host": cfg.Server.Host,
		"server_port": cfg.Server.Port,
		"db_path":     cfg.Database.Path,
	})

	// 初始化SQLite資料庫
	if err := database.Init(); err != nil {
		logger.Fatal("SQLite資料庫初始化失敗", err, logrus.Fields{
			"db_type": "sqlite",
		})
	}
	defer database.Close()

	// 初始化 Repository
	userRepo := models.NewUserRepository(database.DB)
	logger.Info("Repository初始化完成")

	// 初始化 Service
	authService := services.NewAuthService(userRepo, &cfg.JWT)
	userService := services.NewUserService(userRepo)
	adminService := services.NewAdminService(userRepo)
	logger.Info("Service層初始化完成")

	// 初始化 Controller
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	adminController := controllers.NewAdminController(adminService)
	logger.Info("Controller層初始化完成")

	// 設置路由
	router := routes.SetupRoutes(authController, userController, adminController, authService)

	// 設置 Gin 模式
	if cfg.Server.Host == "0.0.0.0" {
		gin.SetMode(gin.ReleaseMode)
		logger.Info("Gin設置為Release模式")
	}

	// 啟動服務器
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}
	
	addr := ":" + port
	logger.Info("服務器準備啟動", logrus.Fields{
		"address": addr,
		"mode":    gin.Mode(),
	})
	
	log.Printf("服務器啟動在 %s", addr)
	log.Fatal(router.Run(addr))
}
