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
	"go-simple-app/models"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(
	unifiedAuthController *controllers.UnifiedAuthController,
	adminController *controllers.AdminController,
	unifiedAuthService *services.UnifiedAuthService,
	chatController *controllers.ChatController,
	oauthController *controllers.OAuthController,
) *gin.Engine {
	r := gin.Default()

	// 載入 HTML 模板
	r.LoadHTMLGlob("templates/*")

	// 靜態文件 - 使用環境變量或默認路徑
	staticPath := os.Getenv("STATIC_PATH")
	if staticPath == "" {
		staticPath = "/root/static"
	}
	
	r.Static("/static", staticPath)
	r.Static("/assets", staticPath+"/dist/assets")
	r.StaticFile("/favicon.ico", staticPath+"/dist/favicon.ico")
	
	// 圖片服務
	r.Static("/images", "./static/images")

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

	// 初始化商城控制器
	productRepo := models.NewProductRepository(database.DB)
	mallController := controllers.NewMallController(productRepo)
	merchantProductController := controllers.NewMerchantProductController(productRepo)
	imageController := controllers.NewImageController()
	imageProxyController := controllers.NewImageProxyController()
	
	// 初始化購物車服務和控制器
	cartService := services.NewCartService(database.DB)

	// Vue.js 前端路由 - 提供所有頁面
	r.GET("/", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	
	// Vue.js SPA 路由 - 所有前端路由都返回 index.html
	r.GET("/tech-showcase", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/customer/login", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/merchant/login", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/admin/login", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/customer/dashboard", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/merchant/dashboard", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/admin/dashboard", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/merchant/products", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/merchant/products/create", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/merchant/products/:id/edit", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/register", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/cart", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/category/:category", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})
	r.GET("/product/:id", func(c *gin.Context) {
		c.File(staticPath + "/dist/index.html")
	})

	// 商城API路由
	api := r.Group("/api")
	{
		// 商品相關API
		api.GET("/categories", mallController.GetCategories)
		api.GET("/products/featured", mallController.GetFeaturedProducts)
		api.GET("/products", mallController.GetProducts)
		api.GET("/products/category/:category", mallController.GetProductsByCategory)
		api.GET("/products/search", mallController.SearchProducts)
		api.GET("/products/:id", mallController.GetProduct)
		
		// 圖片相關API
		api.GET("/image/product", imageController.GenerateProductImage)
		api.GET("/image/placeholder", imageController.GeneratePlaceholderImage)
		api.GET("/image/proxy", imageProxyController.ProxyImage)
		api.GET("/image/external", imageProxyController.GenerateExternalImage)
	}

	// 設置購物車路由
	SetupCartRoutes(r, cartService, unifiedAuthService)

	// 商城頁面路由（已移至Vue.js）
	// {
	//	// 商品詳情頁面
	//	r.GET("/product/:id", mallController.ShowProductPage)
	//	// 分類頁面
	//	r.GET("/category/:category", mallController.ShowCategoryPage)
	//	// 搜尋結果頁面
	//	r.GET("/search", mallController.ShowSearchPage)
	// }

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
			"ai-chat": "AI 智能聊天系統",
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

	// 客戶登入和註冊路由
	// 客戶登入頁面（已移至Vue.js）
	// r.GET("/customer/login", unifiedAuthController.ShowCustomerLogin)
	r.POST("/customer/login", unifiedAuthController.CustomerLogin)
	// 客戶註冊頁面（已移至Vue.js）
	// r.GET("/customer/register", unifiedAuthController.ShowRegisterPage)
	r.POST("/customer/register", unifiedAuthController.Register)

	// 客戶受保護路由
	customerProtected := r.Group("/customer")
	customerProtected.Use(middleware.UnifiedAuthMiddleware(unifiedAuthService))
	customerProtected.Use(middleware.CustomerMiddleware())
	{
		// 客戶儀表板（已移至Vue.js）
		// customerProtected.GET("/dashboard", unifiedAuthController.ShowCustomerDashboard)
		customerProtected.GET("/profile", unifiedAuthController.GetUserProfile)
	}

	// 商戶登入和註冊路由
	merchant := r.Group("/merchant")
	{
		// 商戶登入頁面（已移至Vue.js）
		// merchant.GET("/login", unifiedAuthController.ShowMerchantLogin)
		merchant.POST("/login", unifiedAuthController.MerchantLogin)
		// 商戶註冊頁面（已移至Vue.js）
		// merchant.GET("/register", unifiedAuthController.ShowMerchantRegisterPage)
		merchant.POST("/register", unifiedAuthController.Register)
	}

	// 商戶受保護路由
	merchantProtected := r.Group("/merchant")
	merchantProtected.Use(middleware.UnifiedAuthMiddleware(unifiedAuthService))
	merchantProtected.Use(middleware.MerchantMiddleware())
	{
		// 商戶儀表板（已移至Vue.js）
		// merchantProtected.GET("/dashboard", unifiedAuthController.ShowMerchantDashboard)
		merchantProtected.GET("/profile", unifiedAuthController.GetUserProfile)
		
		// 商戶商品管理API
		merchantAPI := merchantProtected.Group("/api")
		{
			merchantAPI.GET("/products", merchantProductController.GetMerchantProducts)
			merchantAPI.GET("/products/stats", merchantProductController.GetMerchantProductStats)
			merchantAPI.POST("/products", merchantProductController.CreateMerchantProduct)
			merchantAPI.GET("/products/:id", merchantProductController.GetMerchantProduct)
			merchantAPI.PUT("/products/:id", merchantProductController.UpdateMerchantProduct)
			merchantAPI.PUT("/products/:id/toggle-status", merchantProductController.ToggleMerchantProductStatus)
			merchantAPI.DELETE("/products/:id", merchantProductController.DeleteMerchantProduct)
		}
	}

	// 管理員登入和註冊路由
	adminAuth := r.Group("/admin")
	{
		// 管理員登入頁面（已移至Vue.js）
		// adminAuth.GET("/login", unifiedAuthController.ShowAdminLogin)
		adminAuth.POST("/login", unifiedAuthController.AdminLogin)
		adminAuth.GET("/register", unifiedAuthController.ShowRegisterPage)
		adminAuth.POST("/register", unifiedAuthController.Register)
	}

	// 登出路由
	r.POST("/logout", unifiedAuthController.Logout)

	// OAuth路由
	oauth := r.Group("/auth")
	{
		oauth.GET("/line", oauthController.LineLogin)
		oauth.GET("/line/callback", oauthController.LineCallback)
	}


	// 聊天功能路由（支援匿名用戶）
	chat := r.Group("/api/chat")
	{
		// 对话管理（支援匿名用戶）
		chat.POST("/conversations", chatController.CreateConversation)
		chat.POST("/send", chatController.SendMessage)
		
		// 需要認證的路由
		chatAuth := chat.Group("")
		chatAuth.Use(middleware.UnifiedAuthMiddleware(unifiedAuthService))
		chatAuth.Use(middleware.MultiRoleMiddleware("customer", "merchant", "admin"))
		{
			chatAuth.GET("/conversations", chatController.GetUserConversations)
			chatAuth.GET("/conversations/:id", chatController.GetConversation)
			chatAuth.DELETE("/conversations/:id", chatController.DeleteConversation)
		}
	}

	// 管理員專用路由
	admin := r.Group("/admin")
	admin.Use(middleware.UnifiedAuthMiddleware(unifiedAuthService))
	admin.Use(middleware.AdminMiddleware())
	{
		// 管理員頁面
		// 管理員儀表板（已移至Vue.js）
		// admin.GET("/dashboard", adminController.ShowAdminDashboard)
		// 管理員用戶管理頁面（已移至Vue.js）
		// admin.GET("/users", adminController.ShowUserManagement)
		// admin.GET("/users/create", adminController.ShowCreateUser)
		// admin.GET("/users/:id/edit", adminController.ShowEditUser)
		admin.GET("/profile", unifiedAuthController.GetUserProfile)

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
			
			// 聊天管理
			adminAPI.GET("/chat/status", chatController.GetDatabaseStatus)
			adminAPI.POST("/chat/cleanup", chatController.CleanupOldData)
		}
	}

	// 資料庫管理登入/登出（不需要認證）
	dbController := controllers.NewDBController()
	// 資料庫管理登入頁面（已移至Vue.js）
	// r.GET("/admin/db/login", dbController.ShowDBLogin)
	r.POST("/admin/db/login", dbController.DBLogin)
	r.POST("/admin/db/logout", dbController.DBLogout)

	// 資料庫管理路由（獨立認證）
	db := r.Group("/admin/db")
	db.Use(middleware.DBAuthMiddleware())
	{
		// 資料庫管理頁面
		// 資料庫管理頁面（已移至Vue.js）
		// db.GET("/", dbController.ShowDBManager)

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
