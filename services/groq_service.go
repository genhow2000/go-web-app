package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-simple-app/config"
)

// GroqService Groq API服务
type GroqService struct {
	config     config.GroqConfig
	usageStats AIUsageStats
	client     *http.Client
}

// NewGroqService 创建Groq服务
func NewGroqService(cfg config.GroqConfig) *GroqService {
	return &GroqService{
		config: cfg,
		usageStats: AIUsageStats{
			Provider:    "groq",
			DailyUsage:  0,
			DailyLimit:  cfg.DailyLimit,
			LastReset:   time.Now(),
			IsExhausted: false,
			ErrorCount:  0,
			LastError:   "",
			LastUsed:    time.Time{},
		},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateResponse 生成回复
func (s *GroqService) GenerateResponse(ctx context.Context, message, conversationID string) (string, error) {
	// 检查是否超出限制
	if s.usageStats.IsExhausted {
		return "", &AIError{
			Provider:        "groq",
			Message:         "Daily limit exceeded",
			IsQuotaExceeded: true,
		}
	}

	// 构建请求
	requestBody := map[string]interface{}{
		"model": s.config.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": message,
			},
		},
		"max_tokens":   s.config.MaxTokens,
		"temperature":  s.config.Temperature,
		"stream":       false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", &AIError{
			Provider:       "groq",
			Message:        fmt.Sprintf("Failed to marshal request: %v", err),
			IsNetworkError: true,
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", &AIError{
			Provider:       "groq",
			Message:        fmt.Sprintf("Failed to create request: %v", err),
			IsNetworkError: true,
		}
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		s.usageStats.ErrorCount++
		s.usageStats.LastError = err.Error()
		return "", &AIError{
			Provider:       "groq",
			Message:        fmt.Sprintf("Request failed: %v", err),
			IsNetworkError: true,
		}
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		s.usageStats.ErrorCount++
		errorMsg := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		s.usageStats.LastError = errorMsg
		
		if resp.StatusCode == 429 {
			return "", &AIError{
				Provider:      "groq",
				Message:       errorMsg,
				IsRateLimited: true,
			}
		}
		
		return "", &AIError{
			Provider: "groq",
			Message:  errorMsg,
		}
	}

	// 解析响应
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		s.usageStats.ErrorCount++
		s.usageStats.LastError = err.Error()
		return "", &AIError{
			Provider: "groq",
			Message:  fmt.Sprintf("Failed to decode response: %v", err),
		}
	}

	// 更新使用统计
	s.usageStats.DailyUsage++
	s.usageStats.LastUsed = time.Now()

	// 检查是否接近限制
	if s.usageStats.DailyUsage >= int(float64(s.usageStats.DailyLimit)*0.9) {
		s.usageStats.IsExhausted = true
	}

	// 返回回复
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", &AIError{
		Provider: "groq",
		Message:  "No response generated",
	}
}

// GetServiceName 获取服务名称
func (s *GroqService) GetServiceName() string {
	return "Groq API"
}

// IsAvailable 检查服务是否可用
func (s *GroqService) IsAvailable(ctx context.Context) bool {
	// 检查API密钥是否存在
	if s.config.APIKey == "" {
		return false
	}
	
	// 检查是否超出限制
	if s.usageStats.IsExhausted {
		return false
	}
	
	return true
}

// GetUsageStats 获取使用统计
func (s *GroqService) GetUsageStats() map[string]interface{} {
	return map[string]interface{}{
		"provider":        s.usageStats.Provider,
		"daily_usage":     s.usageStats.DailyUsage,
		"daily_limit":     s.usageStats.DailyLimit,
		"last_reset":      s.usageStats.LastReset,
		"is_exhausted":    s.usageStats.IsExhausted,
		"error_count":     s.usageStats.ErrorCount,
		"last_error":      s.usageStats.LastError,
		"last_used":       s.usageStats.LastUsed,
		"usage_percentage": float64(s.usageStats.DailyUsage) / float64(s.usageStats.DailyLimit) * 100,
	}
}
