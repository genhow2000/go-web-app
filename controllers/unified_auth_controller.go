package controllers

import (
	"net/http"
	"go-simple-app/models"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
)

type UnifiedAuthController struct {
	authService *services.UnifiedAuthService
}

func NewUnifiedAuthController(authService *services.UnifiedAuthService) *UnifiedAuthController {
	return &UnifiedAuthController{
		authService: authService,
	}
}

// 顯示客戶登入頁面
func (c *UnifiedAuthController) ShowCustomerLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "customer_login.html", gin.H{
		"title": "客戶登入",
	})
}

// 客戶登入
func (c *UnifiedAuthController) CustomerLogin(ctx *gin.Context) {
	var req services.UnifiedLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
		})
		return
	}

	// 設置角色為 customer
	req.Role = "customer"

	response, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 設置認證 cookie
	ctx.SetCookie("auth_token", response.Token, 3600*24, "/", "", false, true)
	
	// 如果是頁面請求，跳轉到儀表板
	if ctx.Request.Header.Get("Accept") == "text/html" {
		ctx.Redirect(http.StatusFound, "/customer/dashboard")
		return
	}
	
	// 創建包含role字段的用戶響應
	userResponse := gin.H{
		"id":             response.User.GetID(),
		"name":           response.User.GetName(),
		"email":          response.User.GetEmail(),
		"role":           response.User.GetRole(),
		"is_active":      response.User.GetIsActive(),
		"last_login":     response.User.GetLastLogin(),
		"login_count":    response.User.GetLoginCount(),
		"created_at":     response.User.GetCreatedAt(),
		"updated_at":     response.User.GetUpdatedAt(),
	}
	
	// 如果是客戶，添加客戶特有字段
	if customer, ok := response.User.(*models.Customer); ok {
		userResponse["phone"] = customer.Phone
		userResponse["address"] = customer.Address
		userResponse["birth_date"] = customer.BirthDate
		userResponse["gender"] = customer.Gender
		userResponse["email_verified"] = customer.EmailVerified
		userResponse["profile_data"] = customer.ProfileData
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "客戶登入成功",
		"user":    userResponse,
		"token":   response.Token,
	})
}

// 顯示商戶登入頁面
func (c *UnifiedAuthController) ShowMerchantLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "merchant_login.html", gin.H{
		"title": "商戶登入",
	})
}

// 商戶登入
func (c *UnifiedAuthController) MerchantLogin(ctx *gin.Context) {
	var req services.UnifiedLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
		})
		return
	}

	// 設置角色為 merchant
	req.Role = "merchant"

	response, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 設置認證 cookie
	ctx.SetCookie("auth_token", response.Token, 3600*24, "/", "", false, true)
	
	// 如果是頁面請求，跳轉到儀表板
	if ctx.Request.Header.Get("Accept") == "text/html" {
		ctx.Redirect(http.StatusFound, "/merchant/dashboard")
		return
	}
	
	// 創建包含role字段的用戶響應
	userResponse := gin.H{
		"id":             response.User.GetID(),
		"name":           response.User.GetName(),
		"email":          response.User.GetEmail(),
		"role":           response.User.GetRole(),
		"is_active":      response.User.GetIsActive(),
		"last_login":     response.User.GetLastLogin(),
		"login_count":    response.User.GetLoginCount(),
		"created_at":     response.User.GetCreatedAt(),
		"updated_at":     response.User.GetUpdatedAt(),
	}
	
	// 如果是商戶，添加商戶特有字段
	if merchant, ok := response.User.(*models.Merchant); ok {
		userResponse["business_name"] = merchant.BusinessName
		userResponse["business_license"] = merchant.BusinessLicense
		userResponse["phone"] = merchant.Phone
		userResponse["business_type"] = merchant.BusinessType
		userResponse["is_verified"] = merchant.IsVerified
		userResponse["business_data"] = merchant.BusinessData
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "商戶登入成功",
		"user":    userResponse,
		"token":   response.Token,
	})
}

// 顯示管理員登入頁面
func (c *UnifiedAuthController) ShowAdminLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin_login.html", gin.H{
		"title": "管理員登入",
	})
}

// 管理員登入
func (c *UnifiedAuthController) AdminLogin(ctx *gin.Context) {
	var req services.UnifiedLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
		})
		return
	}

	// 設置角色為 admin
	req.Role = "admin"

	response, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 設置認證 cookie
	ctx.SetCookie("auth_token", response.Token, 3600*24, "/", "", false, true)
	
	// 如果是頁面請求，跳轉到儀表板
	if ctx.Request.Header.Get("Accept") == "text/html" {
		ctx.Redirect(http.StatusFound, "/admin/dashboard")
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "管理員登入成功",
		"user":    response.User,
		"token":   response.Token,
	})
}

// 統一註冊
func (c *UnifiedAuthController) Register(ctx *gin.Context) {
	var req services.UnifiedRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
			"details": err.Error(),
		})
		return
	}

	user, err := c.authService.Register(&req)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "註冊成功",
		"user":    user,
	})
}

// 顯示註冊頁面
func (c *UnifiedAuthController) ShowRegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{})
}

// 顯示商戶註冊頁面
func (c *UnifiedAuthController) ShowMerchantRegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "merchant_register.html", gin.H{
		"title": "商戶註冊 - 我要開店",
	})
}

// 顯示客戶儀表板
func (c *UnifiedAuthController) ShowCustomerDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "customer_dashboard.html", gin.H{
		"title": "客戶儀表板",
	})
}

// 顯示商戶儀表板
func (c *UnifiedAuthController) ShowMerchantDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "merchant_dashboard.html", gin.H{
		"title": "商戶儀表板",
	})
}

// 登出
func (c *UnifiedAuthController) Logout(ctx *gin.Context) {
	// 清除 cookie
	ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}

// 獲取用戶資料
func (c *UnifiedAuthController) GetUserProfile(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未認證的用戶",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
