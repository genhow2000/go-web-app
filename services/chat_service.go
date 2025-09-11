package services

import (
	"context"
	"fmt"
	"log"
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
}

// NewChatService 创建聊天服务实例
func NewChatService() *ChatService {
	collection := database.GetMongoDBCollection("conversations")
	if collection == nil {
		log.Println("Warning: MongoDB collection is nil, chat service will not work")
	}
	return &ChatService{
		collection: collection,
	}
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
