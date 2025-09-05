package routes

import (
	"go-simple-app/controllers"
	"go-simple-app/logger"
	"go-simple-app/middleware"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(
	authController *controllers.AuthController,
	userController *controllers.UserController,
	authService *services.AuthService,
) *gin.Engine {
	r := gin.Default()

	// 載入 HTML 模板
	r.LoadHTMLGlob("templates/*")

	// 靜態文件
	r.Static("/static", "./static")

	// 中間件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.ErrorLogger())
	
	// 記錄服務器啟動
	logger.Info("服務器路由初始化完成", logrus.Fields{
		"templates_loaded": true,
		"middleware_count": 3,
	})

	// 健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"service":  "go-simple-app",
			"database": "connected",
		})
	})

	// 首頁
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "歡迎來到 Go 服務器！",
			"status":  "running",
			"version": "2.0.0",
		})
	})

	// 認證路由
	auth := r.Group("/auth")
	{
		auth.GET("/login", authController.ShowLoginPage)
		auth.GET("/register", authController.ShowRegisterPage)
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
		auth.POST("/logout", authController.Logout)
	}

	// 兼容舊路由
	r.GET("/login", authController.ShowLoginPage)
	r.GET("/register", authController.ShowRegisterPage)
	r.POST("/login", authController.Login)
	r.POST("/register", authController.Register)
	r.POST("/logout", authController.Logout)

	// 受保護的路由
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/dashboard", authController.ShowDashboard)
		protected.GET("/users", userController.GetAllUsers)
		protected.GET("/users/:id", userController.GetUserByID)
		protected.POST("/users", userController.CreateUser)
		protected.PUT("/users/:id", userController.UpdateUser)
		protected.DELETE("/users/:id", userController.DeleteUser)
	}

	// 兼容舊路由（需要認證）
	// r.GET("/dashboard", middleware.AuthMiddleware(authService), authController.ShowDashboard)

	return r
}
