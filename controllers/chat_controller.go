package controllers

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-simple-app/database"
	"go-simple-app/models"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatService     *services.ChatService
	rateLimitService *services.RateLimitService
}

// NewChatController 创建聊天控制器
func NewChatController() *ChatController {
	return &ChatController{
		chatService:     services.NewChatService(),
		rateLimitService: services.NewRateLimitService(),
	}
}

// CreateConversation 创建新对话
func (cc *ChatController) CreateConversation(c *gin.Context) {
	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// 获取用户信息（支持匿名用户）
	userID, isAnonymous := cc.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unable to identify user",
		})
		return
	}

	// 检查MongoDB是否可用
	if !database.IsMongoDBConnected() {
		// MongoDB不可用，返回模拟对话ID
		responseData := gin.H{
			"conversation_id": "sim_" + fmt.Sprintf("%d", time.Now().Unix()),
			"success":         true,
			"is_anonymous":    isAnonymous,
			"simulation_mode": true,
		}
		c.JSON(http.StatusOK, responseData)
		return
	}

	response, err := cc.chatService.CreateConversation(userID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create conversation: " + err.Error(),
		})
		return
	}

	// 添加匿名用户标识到响应
	responseData := gin.H{
		"conversation_id": response.ConversationID,
		"success":         response.Success,
		"is_anonymous":    isAnonymous,
	}

	c.JSON(http.StatusOK, responseData)
}

// SendMessage 发送消息
func (cc *ChatController) SendMessage(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// 获取用户信息（支持匿名用户）
	userID, isAnonymous := cc.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unable to identify user",
		})
		return
	}

	// 检查匿名用户的请求限制
	if isAnonymous {
		identifier := cc.getAnonymousIdentifier(c)
		allowed, errorMsg := cc.rateLimitService.CheckRateLimit(identifier, true)
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   errorMsg,
				"rate_limit_exceeded": true,
				"is_anonymous": true,
				"usage_stats": cc.rateLimitService.GetUsageStats(identifier, true),
				"register_url": "/customer/register",
				"login_url": "/customer/login",
			})
			return
		}
	}

	// 检查MongoDB是否可用
	if !database.IsMongoDBConnected() {
		// MongoDB不可用，使用模拟模式
		aiResponse := cc.getFallbackResponse(req.Message)
		
		// 返回模拟响应
		responseData := gin.H{
			"success": true,
			"conversation_id": req.ConversationID,
			"user_message": gin.H{
				"content": req.Message,
				"role": "user",
				"timestamp": time.Now(),
			},
			"ai_message": gin.H{
				"content": aiResponse,
				"role": "assistant", 
				"timestamp": time.Now(),
			},
			"is_anonymous": isAnonymous,
			"simulation_mode": true,
		}

		// 如果是匿名用户，添加使用统计
		if isAnonymous {
			identifier := cc.getAnonymousIdentifier(c)
			usageStats := cc.rateLimitService.GetUsageStats(identifier, true)
			responseData["usage_stats"] = usageStats
		}

		c.JSON(http.StatusOK, responseData)
		return
	}

	// 验证对话是否属于当前用户（包括匿名用户）
	conversation, err := cc.chatService.GetConversation(req.ConversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Conversation not found: " + err.Error(),
		})
		return
	}

	if conversation.UserID != userID {
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

	// 调用AI服务生成回复
	aiResponse, err := cc.chatService.GenerateAIResponse(req.Message, req.ConversationID)
	if err != nil {
		// 如果AI服务失败，返回模拟回复
		aiResponse = cc.getFallbackResponse(req.Message)
	}

	// 添加AI回复到对话
	aiMessage, err := cc.chatService.AddMessage(req.ConversationID, "assistant", aiResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to add AI response: " + err.Error(),
		})
		return
	}

	// 获取使用统计
	var usageStats map[string]interface{}
	if isAnonymous {
		identifier := cc.getAnonymousIdentifier(c)
		usageStats = cc.rateLimitService.GetUsageStats(identifier, true)
	}

	// 返回包含AI回复的响应
	responseData := gin.H{
		"success": true,
		"conversation_id": req.ConversationID,
		"user_message": response.Message,
		"ai_message": aiMessage.Message,
		"is_anonymous": isAnonymous,
	}

	// 如果是匿名用户，添加使用统计
	if isAnonymous {
		responseData["usage_stats"] = usageStats
		// 如果接近限制，添加警告
		if dailyCount, ok := usageStats["daily_requests"].(int); ok {
			dailyLimit := usageStats["daily_limit"].(int)
			if dailyCount >= int(float64(dailyLimit)*0.8) { // 达到80%时警告
				responseData["warning"] = "您今日的使用次數即將達到上限，建議註冊會員獲得更多使用次數"
			}
		}
	}

	c.JSON(http.StatusOK, responseData)
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
	userObj, ok := user.(models.UserInterface)
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
	if conversation.UserID != userObj.GetID() {
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
	userObj, ok := user.(models.UserInterface)
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

	response, err := cc.chatService.GetUserConversations(userObj.GetID(), limit, offset)
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
	userObj, ok := user.(models.UserInterface)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid user format",
		})
		return
	}

	err := cc.chatService.DeleteConversation(conversationID, userObj.GetID())
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

	userObj, ok := user.(models.UserInterface)
	if !ok || userObj.GetRole() != "admin" {
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

	userObj, ok := user.(models.UserInterface)
	if !ok || userObj.GetRole() != "admin" {
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

// getUserID 获取用户ID，支持匿名用户
func (cc *ChatController) getUserID(c *gin.Context) (int, bool) {
	// 首先检查是否有认证用户
	if user, exists := c.Get("user"); exists {
		if userObj, ok := user.(models.UserInterface); ok {
			return userObj.GetID(), false
		}
	}

	// 如果没有认证用户，创建匿名用户ID
	anonymousID := cc.getAnonymousUserID(c)
	return anonymousID, true
}

// getAnonymousUserID 生成匿名用户ID
func (cc *ChatController) getAnonymousUserID(c *gin.Context) int {
	// 使用IP地址和User-Agent生成匿名用户ID
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	
	// 生成基于IP和User-Agent的哈希
	hash := md5.Sum([]byte(ip + userAgent + time.Now().Format("2006-01-02")))
	hashStr := fmt.Sprintf("%x", hash)
	
	// 将哈希转换为负数，避免与真实用户ID冲突
	// 取前8位字符转换为数字，然后转为负数
	hashNum := 0
	for i := 0; i < 8; i++ {
		hashNum = hashNum*16 + int(hashStr[i])
		if hashStr[i] >= 'a' {
			hashNum = hashNum - 10 + int(hashStr[i]-'a')
		} else {
			hashNum = hashNum - int(hashStr[i]-'0')
		}
	}
	
	// 确保是负数
	if hashNum > 0 {
		hashNum = -hashNum
	}
	
	return hashNum
}

// getAnonymousIdentifier 获取匿名用户标识符（用于限制检查）
func (cc *ChatController) getAnonymousIdentifier(c *gin.Context) string {
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	return cc.rateLimitService.GenerateAnonymousIdentifier(ip, userAgent)
}

// getFallbackResponse 获取备用回复
func (cc *ChatController) getFallbackResponse(message string) string {
	// 简单的关键词匹配回复
	message = strings.ToLower(message)
	
	if strings.Contains(message, "價格") || strings.Contains(message, "多少錢") {
		return "我們有各種價格區間的商品，從經濟實惠到高端精品都有。您可以在商品頁面查看詳細價格信息。"
	} else if strings.Contains(message, "推薦") || strings.Contains(message, "建議") {
		return "根據您的需求，我推薦您查看我們的精選商品。這些商品都經過嚴格篩選，品質有保證。"
	} else if strings.Contains(message, "配送") || strings.Contains(message, "運費") {
		return "我們提供快速配送服務，24小時內發貨，3-5天送達。滿額還有免運費優惠！"
	} else if strings.Contains(message, "退換") || strings.Contains(message, "售後") {
		return "我們提供7天無理由退換貨服務，讓您買得放心。如有任何問題，我們的客服團隊隨時為您服務。"
	} else if strings.Contains(message, "你好") || strings.Contains(message, "您好") {
		return "您好！我是阿和商城的AI購物助手，很高興為您服務！有什麼可以幫助您的嗎？"
	} else {
		responses := []string{
			"我了解您的需求，讓我為您推薦一些相關商品。",
			"這是一個很好的問題！根據您的描述，我建議您查看以下分類的商品。",
			"感謝您的詢問！我可以幫您找到最適合的商品。",
			"我明白您想要什麼了，讓我為您搜索相關商品。",
			"好的，我會根據您的需求為您推薦商品。",
		}
		return responses[time.Now().Unix()%int64(len(responses))]
	}
}
