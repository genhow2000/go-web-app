package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message 表示聊天中的单条消息
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Role      string             `bson:"role" json:"role"`           // "user" 或 "assistant"
	Content   string             `bson:"content" json:"content"`     // 消息内容
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"` // 消息时间
}

// Conversation 表示一次完整的对话
type Conversation struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      int                `bson:"user_id" json:"user_id"`           // 关联SQLite用户ID
	Title       string             `bson:"title" json:"title"`               // 对话标题
	Messages    []Message          `bson:"messages" json:"messages"`         // 消息列表
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`     // 更新时间
	IsActive    bool               `bson:"is_active" json:"is_active"`       // 是否活跃
}

// ChatRequest 表示发送消息的请求
type ChatRequest struct {
	ConversationID string                 `json:"conversation_id" binding:"required"`
	Message        string                 `json:"message" binding:"required"`
	StockContext   map[string]interface{} `json:"stock_context,omitempty"`
}

// ChatResponse 表示聊天响应
type ChatResponse struct {
	ConversationID string    `json:"conversation_id"`
	Message        Message   `json:"message"`
	Success        bool      `json:"success"`
	Error          string    `json:"error,omitempty"`
}

// ConversationListResponse 表示对话列表响应
type ConversationListResponse struct {
	Conversations []ConversationSummary `json:"conversations"`
	Total         int                   `json:"total"`
	Success       bool                  `json:"success"`
	Error         string                `json:"error,omitempty"`
}

// ConversationSummary 表示对话摘要
type ConversationSummary struct {
	ID        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	LastMessage string           `json:"last_message"`
	UpdatedAt time.Time          `json:"updated_at"`
	MessageCount int             `json:"message_count"`
}

// CreateConversationRequest 表示创建对话的请求
type CreateConversationRequest struct {
	Title string `json:"title" binding:"required"`
}

// CreateConversationResponse 表示创建对话的响应
type CreateConversationResponse struct {
	ConversationID string `json:"conversation_id"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
}
