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

// GeminiService Google Gemini API服务
type GeminiService struct {
	config     config.GeminiConfig
	usageStats AIUsageStats
	client     *http.Client
}

// NewGeminiService 创建Gemini服务
func NewGeminiService(cfg config.GeminiConfig) *GeminiService {
	return &GeminiService{
		config: cfg,
		usageStats: AIUsageStats{
			Provider:    "gemini",
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
func (s *GeminiService) GenerateResponse(ctx context.Context, message, conversationID string, stockContext map[string]interface{}) (string, error) {
	// 检查是否超出限制
	if s.usageStats.IsExhausted {
		return "", &AIError{
			Provider:        "gemini",
			Message:         "Daily limit exceeded",
			IsQuotaExceeded: true,
		}
	}

	// 构建消息内容，包含簡化的股票上下文
	content := message
	if stockContext != nil {
		stockInfo := ""
		if code, ok := stockContext["code"].(string); ok {
			stockInfo += fmt.Sprintf("股票代碼: %s", code)
		}
		if name, ok := stockContext["name"].(string); ok {
			stockInfo += fmt.Sprintf(" (%s)", name)
		}
		if currentPrice, ok := stockContext["current_price"].(float64); ok && currentPrice > 0 {
			stockInfo += fmt.Sprintf(" 現價: %.2f", currentPrice)
		}
		if change, ok := stockContext["change"].(float64); ok {
			stockInfo += fmt.Sprintf(" 漲跌: %.2f", change)
		}
		if market, ok := stockContext["market"].(string); ok {
			stockInfo += fmt.Sprintf(" (%s)", market)
		}
		
		if stockInfo != "" {
			content = fmt.Sprintf("股票: %s\n問題: %s", stockInfo, message)
		}
	}

	// 构建请求
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": content,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": s.config.MaxTokens,
			"temperature":     s.config.Temperature,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", &AIError{
			Provider:       "gemini",
			Message:        fmt.Sprintf("Failed to marshal request: %v", err),
			IsNetworkError: true,
		}
	}

	// 创建请求URL
	url := fmt.Sprintf("%s?key=%s", s.config.APIURL, s.config.APIKey)

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", &AIError{
			Provider:       "gemini",
			Message:        fmt.Sprintf("Failed to create request: %v", err),
			IsNetworkError: true,
		}
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		s.usageStats.ErrorCount++
		s.usageStats.LastError = err.Error()
		return "", &AIError{
			Provider:       "gemini",
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
				Provider:      "gemini",
				Message:       errorMsg,
				IsRateLimited: true,
			}
		}
		
		return "", &AIError{
			Provider: "gemini",
			Message:  errorMsg,
		}
	}

	// 解析响应
	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			TotalTokenCount int `json:"totalTokenCount"`
		} `json:"usageMetadata"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		s.usageStats.ErrorCount++
		s.usageStats.LastError = err.Error()
		return "", &AIError{
			Provider: "gemini",
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
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		return response.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", &AIError{
		Provider: "gemini",
		Message:  "No response generated",
	}
}

// GetServiceName 获取服务名称
func (s *GeminiService) GetServiceName() string {
	return "Google Gemini API"
}

// IsAvailable 检查服务是否可用
func (s *GeminiService) IsAvailable(ctx context.Context) bool {
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
func (s *GeminiService) GetUsageStats() map[string]interface{} {
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
