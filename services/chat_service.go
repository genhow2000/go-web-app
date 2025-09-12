package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go-simple-app/database"
	"go-simple-app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatService struct {
	collection *mongo.Collection
	aiManager  *AIManager
}

// NewChatService 创建聊天服务实例
func NewChatService() *ChatService {
	collection := database.GetMongoDBCollection("conversations")
	if collection == nil {
		log.Println("Warning: MongoDB collection is nil, chat service will not work")
	}
	return &ChatService{
		collection: collection,
		aiManager:  nil, // 将在外部设置
	}
}

// NewChatServiceWithAI 创建带AI管理器的聊天服务实例
func NewChatServiceWithAI(aiManager *AIManager) *ChatService {
	collection := database.GetMongoDBCollection("conversations")
	if collection == nil {
		log.Println("Warning: MongoDB collection is nil, chat service will not work")
	}
	return &ChatService{
		collection: collection,
		aiManager:  aiManager,
	}
}

// SetAIManager 设置AI管理器
func (s *ChatService) SetAIManager(aiManager *AIManager) {
	s.aiManager = aiManager
}

// CreateConversation 创建新对话
func (s *ChatService) CreateConversation(userID int, title string) (*models.CreateConversationResponse, error) {
	if s.collection == nil {
		return nil, fmt.Errorf("MongoDB not connected")
	}

	conversation := models.Conversation{
		UserID:    userID,
		Title:     title,
		Messages:  []models.Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := s.collection.InsertOne(ctx, conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	conversationID := result.InsertedID.(primitive.ObjectID).Hex()
	
	return &models.CreateConversationResponse{
		ConversationID: conversationID,
		Success:        true,
	}, nil
}

// AddMessage 添加消息到对话
func (s *ChatService) AddMessage(conversationID string, role, content string) (*models.ChatResponse, error) {
	if s.collection == nil {
		return nil, fmt.Errorf("MongoDB not connected")
	}

	objID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	message := models.Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 更新对话，添加新消息
	filter := bson.M{"_id": objID, "is_active": true}
	update := bson.M{
		"$push": bson.M{"messages": message},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to add message: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("conversation not found or inactive")
	}

	return &models.ChatResponse{
		ConversationID: conversationID,
		Message:        message,
		Success:        true,
	}, nil
}

// GetConversation 获取对话详情
func (s *ChatService) GetConversation(conversationID string) (*models.Conversation, error) {
	if s.collection == nil {
		return nil, fmt.Errorf("MongoDB not connected")
	}

	objID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var conversation models.Conversation
	filter := bson.M{"_id": objID, "is_active": true}
	
	err = s.collection.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conversation, nil
}

// GetUserConversations 获取用户的所有对话
func (s *ChatService) GetUserConversations(userID int, limit, offset int) (*models.ConversationListResponse, error) {
	if s.collection == nil {
		return nil, fmt.Errorf("MongoDB not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询条件
	filter := bson.M{
		"user_id":   userID,
		"is_active": true,
	}

	// 计算总数
	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count conversations: %w", err)
	}

	// 查询对话列表
	opts := options.Find().
		SetSort(bson.M{"updated_at": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversations: %w", err)
	}
	defer cursor.Close(ctx)

	var conversations []models.Conversation
	if err = cursor.All(ctx, &conversations); err != nil {
		return nil, fmt.Errorf("failed to decode conversations: %w", err)
	}

	// 转换为摘要格式
	var summaries []models.ConversationSummary
	for _, conv := range conversations {
		lastMessage := ""
		messageCount := len(conv.Messages)
		if messageCount > 0 {
			lastMessage = conv.Messages[messageCount-1].Content
			if len(lastMessage) > 100 {
				lastMessage = lastMessage[:100] + "..."
			}
		}

		summaries = append(summaries, models.ConversationSummary{
			ID:           conv.ID,
			Title:        conv.Title,
			LastMessage:  lastMessage,
			UpdatedAt:    conv.UpdatedAt,
			MessageCount: messageCount,
		})
	}

	return &models.ConversationListResponse{
		Conversations: summaries,
		Total:         int(total),
		Success:       true,
	}, nil
}

// DeleteConversation 删除对话（软删除）
func (s *ChatService) DeleteConversation(conversationID string, userID int) error {
	if s.collection == nil {
		return fmt.Errorf("MongoDB not connected")
	}

	objID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID, "user_id": userID, "is_active": true}
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("conversation not found or not owned by user")
	}

	return nil
}

// CleanupOldConversations 清理旧对话（用于管理512MB限制）
func (s *ChatService) CleanupOldConversations(daysOld int) error {
	if s.collection == nil {
		return fmt.Errorf("MongoDB not connected")
	}

	cutoffDate := time.Now().AddDate(0, 0, -daysOld)
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"updated_at": bson.M{"$lt": cutoffDate},
		"is_active":  true,
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	result, err := s.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to cleanup old conversations: %w", err)
	}

	log.Printf("Cleaned up %d old conversations (older than %d days)", result.ModifiedCount, daysOld)
	return nil
}

// GetDatabaseSize 获取数据库大小（用于监控512MB限制）
func (s *ChatService) GetDatabaseSize() (int64, error) {
	if s.collection == nil {
		return 0, fmt.Errorf("MongoDB not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用聚合管道计算集合大小
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":   nil,
			"size":  bson.M{"$sum": bson.M{"$bsonSize": "$$ROOT"}},
		}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate database size: %w", err)
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, fmt.Errorf("failed to decode size result: %w", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	size, ok := result[0]["size"].(int32)
	if !ok {
		return 0, fmt.Errorf("invalid size result")
	}

	return int64(size), nil
}

// GenerateAIResponse 生成AI回复
func (s *ChatService) GenerateAIResponse(message, conversationID string) (string, error) {
	// 使用AI管理器生成回复
	if s.aiManager != nil {
		ctx := context.Background()
		return s.aiManager.GenerateResponse(ctx, message, conversationID)
	}
	// 如果AI管理器未初始化，返回模拟回复
	return s.getSimulatedAIResponse(message), nil
}

// getSimulatedAIResponse 获取模拟AI回复
func (s *ChatService) getSimulatedAIResponse(message string) string {
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
	} else if strings.Contains(message, "商品") || strings.Contains(message, "產品") {
		return "我們有豐富的商品選擇，包括電子產品、服飾、家居用品等。您可以瀏覽我們的商品分類來找到您需要的商品。"
	} else if strings.Contains(message, "優惠") || strings.Contains(message, "折扣") {
		return "我們經常推出各種優惠活動！目前有新年特惠，全場8折優惠，滿額還免運費。使用優惠碼 NEWYEAR2025 即可享受優惠！"
	} else {
		responses := []string{
			"我了解您的需求，讓我為您推薦一些相關商品。",
			"這是一個很好的問題！根據您的描述，我建議您查看以下分類的商品。",
			"感謝您的詢問！我可以幫您找到最適合的商品。",
			"我明白您想要什麼了，讓我為您搜索相關商品。",
			"好的，我會根據您的需求為您推薦商品。",
			"請告訴我更多關於您需求的細節，我會為您提供更精確的建議。",
			"我們有很多優質商品可以滿足您的需求，讓我為您介紹一下。",
		}
		return responses[time.Now().Unix()%int64(len(responses))]
	}
}
