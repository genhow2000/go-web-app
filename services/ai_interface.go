package services

import (
	"context"
	"time"
)

// AIService 定义AI服务的通用接口
type AIService interface {
	// GenerateResponse 生成AI回复
	GenerateResponse(ctx context.Context, message, conversationID string, stockContext map[string]interface{}) (string, error)
	
	// GetServiceName 获取服务名称
	GetServiceName() string
	
	// IsAvailable 检查服务是否可用
	IsAvailable(ctx context.Context) bool
	
	// GetUsageStats 获取使用统计
	GetUsageStats() map[string]interface{}
}

// AIUsageStats AI服务使用统计
type AIUsageStats struct {
	Provider       string    `json:"provider"`
	DailyUsage     int       `json:"daily_usage"`
	DailyLimit     int       `json:"daily_limit"`
	LastReset      time.Time `json:"last_reset"`
	IsExhausted    bool      `json:"is_exhausted"`
	ErrorCount     int       `json:"error_count"`
	LastError      string    `json:"last_error"`
	LastUsed       time.Time `json:"last_used"`
}

// AIError AI服务错误
type AIError struct {
	Provider        string `json:"provider"`
	Message         string `json:"message"`
	IsQuotaExceeded bool   `json:"is_quota_exceeded"`
	IsRateLimited   bool   `json:"is_rate_limited"`
	IsNetworkError  bool   `json:"is_network_error"`
}

// Error 实现error接口
func (e *AIError) Error() string {
	return e.Message
}

// IsQuotaExceededError 检查是否为配额超限错误
func (e *AIError) IsQuotaExceededError() bool {
	return e.IsQuotaExceeded
}

// IsRateLimitedError 检查是否为速率限制错误
func (e *AIError) IsRateLimitedError() bool {
	return e.IsRateLimited
}

// IsNetworkErrorType 检查是否为网络错误
func (e *AIError) IsNetworkErrorType() bool {
	return e.IsNetworkError
}