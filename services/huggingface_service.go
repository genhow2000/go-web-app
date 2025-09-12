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

// HuggingFaceService Hugging Face Inference API服务
type HuggingFaceService struct {
	config     config.HuggingFaceConfig
	usageStats AIUsageStats
	client     *http.Client
}

// NewHuggingFaceService 创建HuggingFace服务
func NewHuggingFaceService(cfg config.HuggingFaceConfig) *HuggingFaceService {
	return &HuggingFaceService{
		config: cfg,
		usageStats: AIUsageStats{
			Provider:    "huggingface",
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
func (s *HuggingFaceService) GenerateResponse(ctx context.Context, message, conversationID string) (string, error) {
	// 检查是否超出限制
	if s.usageStats.IsExhausted {
		return "", &AIError{
			Provider:        "huggingface",
			Message:         "Daily limit exceeded",
			IsQuotaExceeded: true,
		}
	}

	// 构建请求
	requestBody := map[string]interface{}{
		"inputs": message,
		"parameters": map[string]interface{}{
			"max_length": s.config.MaxTokens,
			"temperature": s.config.Temperature,
			"do_sample": true,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", &AIError{
			Provider:       "huggingface",
			Message:        fmt.Sprintf("Failed to marshal request: %v", err),
			IsNetworkError: true,
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", &AIError{
			Provider:       "huggingface",
			Message:        fmt.Sprintf("Failed to create request: %v", err),
			IsNetworkError: true,
		}
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if s.config.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.config.APIToken)
	}

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		s.usageStats.ErrorCount++
		s.usageStats.LastError = err.Error()
		return "", &AIError{
			Provider:       "huggingface",
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
				Provider:      "huggingface",
				Message:       errorMsg,
				IsRateLimited: true,
			}
		}
		
		return "", &AIError{
			Provider: "huggingface",
			Message:  errorMsg,
		}
	}

	// 解析响应
	var response []struct {
		GeneratedText string `json:"generated_text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		s.usageStats.ErrorCount++
		s.usageStats.LastError = err.Error()
		return "", &AIError{
			Provider: "huggingface",
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
	if len(response) > 0 {
		return response[0].GeneratedText, nil
	}

	return "", &AIError{
		Provider: "huggingface",
		Message:  "No response generated",
	}
}

// GetServiceName 获取服务名称
func (s *HuggingFaceService) GetServiceName() string {
	return "Hugging Face Inference API"
}

// IsAvailable 检查服务是否可用
func (s *HuggingFaceService) IsAvailable(ctx context.Context) bool {
	// 检查API URL是否存在
	if s.config.APIURL == "" {
		return false
	}
	
	// 检查是否超出限制
	if s.usageStats.IsExhausted {
		return false
	}
	
	return true
}

// GetUsageStats 获取使用统计
func (s *HuggingFaceService) GetUsageStats() map[string]interface{} {
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
