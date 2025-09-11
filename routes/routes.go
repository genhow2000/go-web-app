package routes

import (
	"net/http"
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
	adminController *controllers.AdminController,
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

	// 系統統計
	r.GET("/api/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"system": gin.H{
				"status":    "running",
				"version":   "2.0.0",
				"uptime":    "24小時",
				"platform":  "Google Cloud Run",
				"region":    "asia-east1",
			},
			"database": gin.H{
				"type":      "SQLite",
				"status":    "connected",
				"migrations": "001_add_role_and_status",
			},
			"features": []string{
				"用戶管理系統",
				"JWT 認證",
				"管理後台",
				"資料庫管理",
				"API 服務",
				"Redis 快取",
				"MongoDB 文檔資料庫",
				"雲端部署",
			},
			"tech_stack": gin.H{
				"backend":  []string{"Go 1.21", "Gin", "JWT", "bcrypt"},
				"frontend": []string{"HTML5", "CSS3", "JavaScript"},
				"database": []string{"SQLite", "資料庫遷移"},
				"deploy":   []string{"Docker", "Google Cloud Run", "CI/CD"},
			},
		})
	})

	// 首頁
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "homepage.html", gin.H{
			"version":           "2.0.0",
			"uptime":            "24小時",
			"migration_version": "001",
		})
	})

	// 商戶登入路由
	merchant := r.Group("/merchant")
	{
		merchant.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "merchant_login.html", gin.H{
				"title": "商戶登入",
			})
		})
		merchant.POST("/login", authController.MerchantLogin)
		merchant.GET("/dashboard", authController.ShowMerchantDashboard)
	}

	// 管理員登入路由
	adminAuth := r.Group("/admin")
	{
		adminAuth.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin_login.html", gin.H{
				"title": "管理員登入",
			})
		})
		adminAuth.POST("/login", authController.AdminLogin)
	}

	// 認證路由（保持兼容性）
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

	// 管理員專用路由
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService))
	admin.Use(middleware.AdminMiddleware())
	{
		// 管理員頁面
		admin.GET("/dashboard", adminController.ShowAdminDashboard)
		admin.GET("/users", adminController.ShowUserManagement)
		admin.GET("/users/create", adminController.ShowCreateUser)
		admin.GET("/users/:id/edit", adminController.ShowEditUser)

		// 管理員 API
		adminAPI := admin.Group("/api")
		{
			adminAPI.GET("/users", adminController.GetAllUsers)
			adminAPI.GET("/users/role/:role", adminController.GetUsersByRole)
			adminAPI.GET("/users/:id", adminController.GetUserByID)
			adminAPI.POST("/users", adminController.CreateUser)
			adminAPI.PUT("/users/:id", adminController.UpdateUser)
			adminAPI.PUT("/users/:id/status", adminController.UpdateUserStatus)
			adminAPI.PUT("/users/:id/role", adminController.UpdateUserRole)
			adminAPI.DELETE("/users/:id", adminController.DeleteUser)
			adminAPI.GET("/stats", adminController.GetUserStats)
		}
	}

	// 資料庫管理登入/登出（不需要認證）
	dbController := controllers.NewDBController()
	r.GET("/admin/db/login", dbController.ShowDBLogin)
	r.POST("/admin/db/login", dbController.DBLogin)
	r.POST("/admin/db/logout", dbController.DBLogout)

	// 資料庫管理路由（獨立認證）
	db := r.Group("/admin/db")
	db.Use(middleware.DBAuthMiddleware())
	{
		// 資料庫管理頁面
		db.GET("/", dbController.ShowDBManager)

		// 資料庫管理 API
		dbAPI := db.Group("/api")
		{
			dbAPI.GET("/tables", dbController.GetTables)
			dbAPI.GET("/tables/:table/data", dbController.GetTableData)
			dbAPI.POST("/query", dbController.ExecuteQuery)
			dbAPI.GET("/stats", dbController.GetDBStats)
		}
	}

	// 兼容舊路由（需要認證）
	// r.GET("/dashboard", middleware.AuthMiddleware(authService), authController.ShowDashboard)

	return r
}
