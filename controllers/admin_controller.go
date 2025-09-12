package controllers

import (
	"net/http"
	"strconv"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	unifiedAdminService *services.UnifiedAdminService
}

func NewAdminController(unifiedAdminService *services.UnifiedAdminService) *AdminController {
	return &AdminController{
		unifiedAdminService: unifiedAdminService,
	}
}

// 顯示管理員儀表板
func (c *AdminController) ShowAdminDashboard(ctx *gin.Context) {
	// 獲取用戶統計
	stats, err := c.unifiedAdminService.GetUserStats()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "無法獲取統計數據",
		})
		return
	}

	ctx.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"stats": stats,
	})
}

// 顯示用戶管理頁面
func (c *AdminController) ShowUserManagement(ctx *gin.Context) {
	users, err := c.unifiedAdminService.GetAllUsers()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "無法獲取用戶列表",
		})
		return
	}

	// 轉換為模板友好的結構
	var templateUsers []gin.H
	for _, user := range users {
		templateUsers = append(templateUsers, gin.H{
			"ID":        user.GetID(),
			"Name":      user.GetName(),
			"Email":     user.GetEmail(),
			"Role":      user.GetRole(),
			"IsActive":  user.GetIsActive(),
			"LastLogin": user.GetLastLogin(),
			"LoginCount": user.GetLoginCount(),
			"CreatedAt": user.GetCreatedAt(),
			"UpdatedAt": user.GetUpdatedAt(),
		})
	}

	ctx.HTML(http.StatusOK, "admin_users.html", gin.H{
		"users": templateUsers,
	})
}

// 顯示創建用戶頁面
func (c *AdminController) ShowCreateUser(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin_create_user.html", gin.H{})
}

// 顯示編輯用戶頁面
func (c *AdminController) ShowEditUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	user, err := c.unifiedAdminService.GetUserByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "用戶不存在",
		})
		return
	}

	ctx.HTML(http.StatusOK, "admin_edit_user.html", gin.H{
		"user": user,
	})
}

// API: 獲取所有用戶
func (c *AdminController) GetAllUsers(ctx *gin.Context) {
	users, err := c.unifiedAdminService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "無法獲取用戶列表",
		})
		return
	}

	// 轉換為 JSON 友好的結構
	var jsonUsers []gin.H
	for _, user := range users {
		jsonUsers = append(jsonUsers, gin.H{
			"ID":        user.GetID(),
			"Name":      user.GetName(),
			"Email":     user.GetEmail(),
			"Role":      user.GetRole(),
			"IsActive":  user.GetIsActive(),
			"LastLogin": user.GetLastLogin(),
			"LoginCount": user.GetLoginCount(),
			"CreatedAt": user.GetCreatedAt(),
			"UpdatedAt": user.GetUpdatedAt(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": jsonUsers,
	})
}

// API: 根據角色獲取用戶
func (c *AdminController) GetUsersByRole(ctx *gin.Context) {
	role := ctx.Param("role")
	if role != "customer" && role != "merchant" && role != "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的角色",
		})
		return
	}

	users, err := c.unifiedAdminService.GetUsersByRole(role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "無法獲取用戶列表",
		})
		return
	}

	// 轉換為 JSON 友好的結構
	var jsonUsers []gin.H
	for _, user := range users {
		jsonUsers = append(jsonUsers, gin.H{
			"ID":        user.GetID(),
			"Name":      user.GetName(),
			"Email":     user.GetEmail(),
			"Role":      user.GetRole(),
			"IsActive":  user.GetIsActive(),
			"LastLogin": user.GetLastLogin(),
			"LoginCount": user.GetLoginCount(),
			"CreatedAt": user.GetCreatedAt(),
			"UpdatedAt": user.GetUpdatedAt(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": jsonUsers,
	})
}

// API: 獲取用戶詳情
func (c *AdminController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	user, err := c.unifiedAdminService.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "用戶不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// API: 創建用戶
func (c *AdminController) CreateUser(ctx *gin.Context) {
	var req services.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
			"details": err.Error(),
		})
		return
	}

	user, err := c.unifiedAdminService.CreateUser(&req)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "用戶創建成功",
		"user":    user,
	})
}

// API: 更新用戶
func (c *AdminController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	var req services.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
			"details": err.Error(),
		})
		return
	}

	user, err := c.unifiedAdminService.UpdateUser(id, &req)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "用戶更新成功",
		"user":    user,
	})
}

// API: 更新用戶狀態
func (c *AdminController) UpdateUserStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
		})
		return
	}

	if err := c.unifiedAdminService.UpdateUserStatus(id, req.IsActive); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "用戶狀態更新成功",
	})
}

// API: 更新用戶角色
func (c *AdminController) UpdateUserRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=customer merchant admin"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求數據",
		})
		return
	}

	if err := c.unifiedAdminService.UpdateUserRole(id, req.Role); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "用戶角色更新成功",
	})
}

// API: 刪除用戶
func (c *AdminController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	if err := c.unifiedAdminService.DeleteUser(id); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "用戶刪除成功",
	})
}

// API: 獲取用戶統計
func (c *AdminController) GetUserStats(ctx *gin.Context) {
	stats, err := c.unifiedAdminService.GetUserStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "無法獲取統計數據",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}
