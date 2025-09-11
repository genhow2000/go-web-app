package controllers

import (
	"net/http"
	"strconv"

	"go-simple-app/models"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatService *services.ChatService
}

// NewChatController 创建聊天控制器
func NewChatController() *ChatController {
	return &ChatController{
		chatService: services.NewChatService(),
	}
}

// CreateConversation 创建新对话
func (cc *ChatController) CreateConversation(c *gin.Context) {
	// 获取用户信息（从JWT token中解析）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// 转换用户对象
	userObj, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid user format",
		})
		return
	}

	response, err := cc.chatService.CreateConversation(userObj.ID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create conversation: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SendMessage 发送消息
func (cc *ChatController) SendMessage(c *gin.Context) {
	// 获取用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// 转换用户对象
	userObj, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid user format",
		})
		return
	}

	// 验证对话是否属于当前用户
	conversation, err := cc.chatService.GetConversation(req.ConversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Conversation not found: " + err.Error(),
		})
		return
	}

	if conversation.UserID != userObj.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Access denied to this conversation",
		})
		return
	}

	// 添加用户消息
	response, err := cc.chatService.AddMessage(req.ConversationID, "user", req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to send message: " + err.Error(),
		})
		return
	}

	// TODO: 这里应该调用AI服务生成回复
	// 暂时返回用户消息的确认
	c.JSON(http.StatusOK, response)
}

// GetConversation 获取对话详情
func (cc *ChatController) GetConversation(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Conversation ID is required",
		})
		return
	}

	// 获取用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// 转换用户对象
	userObj, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid user format",
		})
		return
	}

	conversation, err := cc.chatService.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Conversation not found: " + err.Error(),
		})
		return
	}

	// 验证对话是否属于当前用户
	if conversation.UserID != userObj.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Access denied to this conversation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"conversation": conversation,
	})
}

// GetUserConversations 获取用户的所有对话
func (cc *ChatController) GetUserConversations(c *gin.Context) {
	// 获取用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// 转换用户对象
	userObj, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid user format",
		})
		return
	}

	// 获取分页参数
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// 限制最大限制
	if limit > 100 {
		limit = 100
	}

	response, err := cc.chatService.GetUserConversations(userObj.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get conversations: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteConversation 删除对话
func (cc *ChatController) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Conversation ID is required",
		})
		return
	}

	// 获取用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// 转换用户对象
	userObj, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid user format",
		})
		return
	}

	err := cc.chatService.DeleteConversation(conversationID, userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete conversation: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Conversation deleted successfully",
	})
}

// GetDatabaseStatus 获取数据库状态（管理员功能）
func (cc *ChatController) GetDatabaseStatus(c *gin.Context) {
	// 检查是否为管理员
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Admin access required",
		})
		return
	}

	userObj, ok := user.(*models.User)
	if !ok || userObj.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Admin access required",
		})
		return
	}

	// 获取数据库大小
	size, err := cc.chatService.GetDatabaseSize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get database size: " + err.Error(),
		})
		return
	}

	// 计算使用百分比（512MB = 536,870,912 bytes）
	maxSize := int64(536870912) // 512MB
	usagePercent := float64(size) / float64(maxSize) * 100

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"database_size": size,
		"max_size":      maxSize,
		"usage_percent": usagePercent,
		"status":        "healthy",
	})
}

// CleanupOldData 清理旧数据（管理员功能）
func (cc *ChatController) CleanupOldData(c *gin.Context) {
	// 检查是否为管理员
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Admin access required",
		})
		return
	}

	userObj, ok := user.(*models.User)
	if !ok || userObj.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Admin access required",
		})
		return
	}

	// 获取清理天数参数
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	err = cc.chatService.CleanupOldConversations(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to cleanup old data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Old data cleaned up successfully",
		"days":    days,
	})
}
