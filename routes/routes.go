package routes

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"go-simple-app/controllers"
	"go-simple-app/database"
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

	// 通用文檔頁面
	r.GET("/docs/:type", func(c *gin.Context) {
		docType := c.Param("type")
		titleMap := map[string]string{
			"seeder": "自製 Seeder 系統",
			"migration": "自製 Migration 系統",
			"user-management": "用戶管理系統",
			"auth": "安全認證系統",
			"admin": "管理後台系統",
			"db-management": "資料庫管理系統",
			"monitoring": "系統監控",
			"api": "API 服務",
			"redis": "Redis 快取系統",
			"mongodb": "MongoDB 文檔資料庫",
		}
		
		title, exists := titleMap[docType]
		if !exists {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"title": "文檔不存在",
				"message": "請求的文檔類型不存在",
			})
			return
		}
		
		c.HTML(http.StatusOK, "docs.html", gin.H{
			"title": title,
			"docType": docType,
		})
	})
	
	// API 端點：讀取 Markdown 文件
	r.GET("/api/docs/:type", func(c *gin.Context) {
		docType := c.Param("type")
		filePath := fmt.Sprintf("md/%s.md", docType)
		
		// 讀取 Markdown 文件
		content, err := os.ReadFile(filePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "文檔文件不存在",
			})
			return
		}
		
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, string(content))
	})

	// API 端點：獲取狀態詳情
	r.GET("/api/status/:type", func(c *gin.Context) {
		statusType := c.Param("type")
		
		// 創建監控服務實例
		monitorService := services.NewMonitorService(database.DB)
		
		var response gin.H
		
		switch statusType {
		case "system":
			systemInfo := monitorService.GetSystemInfo()
			response = gin.H{
				"sections": []gin.H{
					{
						"title": "系統信息",
						"items": []gin.H{
							{"label": "服務器狀態", "value": systemInfo["status"], "status": "online"},
							{"label": "運行時間", "value": systemInfo["uptime"]},
							{"label": "CPU 使用率", "value": systemInfo["cpu_usage"]},
							{"label": "記憶體使用", "value": systemInfo["memory_usage"]},
							{"label": "磁盤空間", "value": systemInfo["disk_usage"]},
						},
					},
					{
						"title": "Go 運行時信息",
						"items": []gin.H{
							{"label": "Goroutine 數量", "value": fmt.Sprintf("%d 個", systemInfo["go_routines"])},
							{"label": "GC 次數", "value": fmt.Sprintf("%d 次", systemInfo["gc_count"])},
							{"label": "Go 版本", "value": runtime.Version()},
						},
					},
				},
			}
		case "database":
			dbInfo := monitorService.GetDatabaseInfo()
			tableStats := dbInfo["tables"].(map[string]interface{})
			
			// 構建資料表統計項目
			var tableItems []gin.H
			for label, value := range tableStats {
				tableItems = append(tableItems, gin.H{
					"label": label,
					"value": value,
				})
			}
			
			response = gin.H{
				"sections": []gin.H{
					{
						"title": "資料庫信息",
						"items": []gin.H{
							{"label": "資料庫類型", "value": dbInfo["type"]},
							{"label": "連接狀態", "value": dbInfo["status"], "status": "online"},
							{"label": "資料庫大小", "value": dbInfo["size"]},
							{"label": "遷移版本", "value": dbInfo["migration"].(map[string]interface{})["version"]},
							{"label": "最後更新", "value": dbInfo["last_update"]},
						},
					},
					{
						"title": "資料表統計",
						"items": tableItems,
					},
				},
			}
		case "api":
			apiInfo := monitorService.GetAPIInfo()
			endpoints := apiInfo["endpoints"].([]string)
			
			// 構建 API 端點項目
			var endpointItems []gin.H
			for _, endpoint := range endpoints {
				endpointItems = append(endpointItems, gin.H{
					"label": endpoint,
					"value": "可用",
					"status": "online",
				})
			}
			
			response = gin.H{
				"sections": []gin.H{
					{
						"title": "API 服務信息",
						"items": []gin.H{
							{"label": "API 版本", "value": apiInfo["version"]},
							{"label": "服務狀態", "value": apiInfo["status"], "status": "online"},
							{"label": "響應時間", "value": apiInfo["response_time"]},
							{"label": "請求處理量", "value": apiInfo["request_count"]},
							{"label": "錯誤率", "value": apiInfo["error_rate"]},
						},
					},
					{
						"title": "API 端點",
						"items": endpointItems,
					},
				},
			}
		case "cloud":
			cloudInfo := monitorService.GetCloudInfo()
			envVars := cloudInfo["env_vars"].([]string)
			cicdInfo := cloudInfo["cicd"].(map[string]interface{})
			gitInfo := cloudInfo["git"].(map[string]interface{})
			buildSteps := cicdInfo["build_steps"].([]string)
			
			// 構建環境變數項目
			var envItems []gin.H
			for _, envVar := range envVars {
				envItems = append(envItems, gin.H{
					"label": envVar,
					"value": os.Getenv(envVar),
				})
			}
			
			// 構建 CI/CD 構建步驟項目
			var buildStepItems []gin.H
			for i, step := range buildSteps {
				buildStepItems = append(buildStepItems, gin.H{
					"label": fmt.Sprintf("步驟 %d", i+1),
					"value": step,
					"status": "online",
				})
			}
			
			response = gin.H{
				"sections": []gin.H{
					{
						"title": "雲端部署信息",
						"items": []gin.H{
							{"label": "平台", "value": cloudInfo["platform"]},
							{"label": "區域", "value": cloudInfo["region"]},
							{"label": "部署狀態", "value": cloudInfo["status"], "status": "online"},
							{"label": "實例數", "value": cloudInfo["instances"]},
							{"label": "CPU 分配", "value": cloudInfo["cpu_allocated"]},
							{"label": "記憶體分配", "value": cloudInfo["memory_allocated"]},
						},
					},
					{
						"title": "CI/CD 流水線",
						"items": []gin.H{
							{"label": "CI/CD 平台", "value": cicdInfo["platform"]},
							{"label": "觸發方式", "value": cicdInfo["trigger_type"]},
							{"label": "構建時間", "value": cicdInfo["build_time"]},
							{"label": "部署時間", "value": cicdInfo["deploy_time"]},
							{"label": "最後構建", "value": cicdInfo["last_build"], "status": "online"},
						},
					},
					{
						"title": "Git 版本控制",
						"items": []gin.H{
							{"label": "代碼倉庫", "value": gitInfo["repository"]},
							{"label": "分支", "value": gitInfo["branch"]},
							{"label": "提交哈希", "value": gitInfo["commit_hash"]},
							{"label": "最後提交", "value": gitInfo["last_commit"]},
							{"label": "自動部署", "value": gitInfo["auto_deploy"], "status": "online"},
						},
					},
					{
						"title": "構建步驟",
						"items": buildStepItems,
					},
					{
						"title": "部署配置",
						"items": []gin.H{
							{"label": "容器映像", "value": cloudInfo["container_image"]},
							{"label": "端口", "value": cloudInfo["port"]},
							{"label": "健康檢查", "value": cloudInfo["health_check"]},
							{"label": "自動擴展", "value": cloudInfo["auto_scale"]},
						},
					},
					{
						"title": "環境變數",
						"items": envItems,
					},
				},
			}
		default:
			c.JSON(http.StatusNotFound, gin.H{"error": "未知的狀態類型"})
			return
		}
		
		c.JSON(http.StatusOK, response)
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

	// 註冊路由（保持兼容性）
	auth := r.Group("/auth")
	{
		auth.GET("/register", authController.ShowRegisterPage)
		auth.POST("/register", authController.Register)
		auth.POST("/logout", authController.Logout)
	}

	// 兼容舊路由
	r.GET("/register", authController.ShowRegisterPage)
	r.POST("/register", authController.Register)
	r.POST("/logout", authController.Logout)

	// 受保護的路由
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
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


	return r
}
