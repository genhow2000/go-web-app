package controllers

import (
	"net/http"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req services.RegisterRequest
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

func (c *AuthController) Login(ctx *gin.Context) {
	var req services.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
			"details": err.Error(),
		})
		return
	}

	response, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "登入成功",
		"token":   response.Token,
		"user":    response.User,
	})
}

func (c *AuthController) ShowLoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{})
}

func (c *AuthController) ShowRegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{})
}

func (c *AuthController) ShowDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "dashboard.html", gin.H{})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}

// MerchantLogin 商戶登入
func (c *AuthController) MerchantLogin(ctx *gin.Context) {
	var req services.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
		})
		return
	}

	response, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 檢查是否為商戶角色
	if response.User.Role != "merchant" && response.User.Role != "user" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "此帳號無權限登入商戶後台",
		})
		return
	}

	// 設置認證 cookie
	ctx.SetCookie("auth_token", response.Token, 3600*24, "/", "", false, true)
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "商戶登入成功",
		"user":    response.User,
	})
}

// AdminLogin 管理員登入
func (c *AuthController) AdminLogin(ctx *gin.Context) {
	var req services.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
		})
		return
	}

	response, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 檢查是否為管理員角色
	if response.User.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "此帳號無權限登入管理員後台",
		})
		return
	}

	// 設置認證 cookie
	ctx.SetCookie("auth_token", response.Token, 3600*24, "/", "", false, true)
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "管理員登入成功",
		"user":    response.User,
	})
}

// ShowMerchantDashboard 顯示商戶後台
func (c *AuthController) ShowMerchantDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "merchant_dashboard.html", gin.H{
		"title": "商戶後台",
	})
}
