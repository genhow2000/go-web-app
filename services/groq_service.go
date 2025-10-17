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
func (s *GroqService) GenerateResponse(ctx context.Context, message, conversationID string, stockContext map[string]interface{}) (string, error) {
	// 检查是否超出限制
	if s.usageStats.IsExhausted {
		return "", &AIError{
			Provider:        "groq",
			Message:         "Daily limit exceeded",
			IsQuotaExceeded: true,
		}
	}

	// 构建消息内容，包含增強的股票上下文
	content := s.buildEnhancedPrompt(message, stockContext)

	// 构建请求
	requestBody := map[string]interface{}{
		"model": s.config.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": content,
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

// buildEnhancedPrompt 構建增強的提示詞
func (s *GroqService) buildEnhancedPrompt(message string, stockContext map[string]interface{}) string {
	if stockContext == nil {
		return message
	}
	
	// 提取股票基本資訊
	code, name, market, currentPrice, change := extractStockInfo(stockContext)
	
	// 構建股票資訊字串
	stockInfo := fmt.Sprintf("股票代碼: %s (%s)", code, name)
	if currentPrice > 0 {
		stockInfo += fmt.Sprintf(" 現價: %.2f", currentPrice)
	}
	if change != 0 {
		stockInfo += fmt.Sprintf(" 漲跌: %.2f", change)
	}
	if market != "" {
		stockInfo += fmt.Sprintf(" (%s)", market)
	}
	
	// 檢查是否有查詢指令
	queryInstructions, hasInstructions := stockContext["query_instructions"].(map[string]interface{})
	
	// 構建專門的提示詞
	prompt := fmt.Sprintf("你是專業股票分析師。分析股票：%s\n\n", stockInfo)
	
	if hasInstructions {
		shouldQuery, _ := queryInstructions["should_query_history"].(bool)
		questionType, _ := queryInstructions["question_type"].(string)
		contextNote, _ := queryInstructions["context_note"].(string)
		
		if shouldQuery {
			prompt += "請模擬搜尋台灣證交所、Yahoo Finance、鉅亨網等資料源，提供專業分析：\n\n"
			
			// 根據問題類型提供專門的分析指導
			switch questionType {
			case "investment_advice":
				prompt += "**投資建議分析重點：**\n"
				prompt += "• 基本面分析（營收、獲利、成長性、估值）\n"
				prompt += "• 技術面支撐和阻力位分析\n"
				prompt += "• 產業趨勢和競爭優勢評估\n"
				prompt += "• 市場情緒和資金流向分析\n"
				prompt += "• 給出明確的買入/持有/賣出建議\n\n"
				
			case "technical_analysis":
				prompt += "**技術指標分析重點：**\n"
				prompt += "• RSI相對強弱指標分析（數值、信號、趨勢）\n"
				prompt += "• MACD動量指標分析（金叉死叉、背離）\n"
				prompt += "• KD隨機指標分析（超買超賣、交叉信號）\n"
				prompt += "• 移動平均線系統分析（多空排列、支撐阻力）\n"
				prompt += "• 布林帶通道分析（突破、收斂、擴張）\n"
				prompt += "• 成交量指標分析（量價關係、資金流向）\n"
				prompt += "• 圖表形態識別（頭肩頂底、三角形、旗形等）\n\n"
				
			case "risk_analysis":
				prompt += "**風險評估分析重點：**\n"
				prompt += "• 股價波動性分析（標準差、Beta值、最大回撤）\n"
				prompt += "• 流動性風險評估（成交量、買賣價差）\n"
				prompt += "• 基本面風險因子（財務結構、獲利穩定性）\n"
				prompt += "• 市場風險和系統性風險\n"
				prompt += "• 公司特定風險（管理層、治理結構）\n"
				prompt += "• 產業風險和政策風險\n"
				prompt += "• 風險控制建議和停損點設定\n\n"
				
			case "fundamental_analysis":
				prompt += "**基本面分析重點：**\n"
				prompt += "• 財務報表分析（損益表、資產負債表、現金流量表）\n"
				prompt += "• 獲利能力分析（毛利率、營業利益率、淨利率）\n"
				prompt += "• 成長性分析（營收成長、獲利成長、ROE/ROA）\n"
				prompt += "• 財務結構分析（負債比率、流動比率、速動比率）\n"
				prompt += "• 產業地位和競爭優勢評估\n"
				prompt += "• 管理層品質和公司治理\n"
				prompt += "• 未來展望和成長動能分析\n\n"
				
			default:
				prompt += "• 技術指標分析（RSI、MACD、KD、移動平均線）\n"
				prompt += "• 支撐位和阻力位分析\n"
				prompt += "• 投資風險評估和建議\n"
			}
			
			// 添加專門的上下文說明
			if contextNote != "" {
				prompt += fmt.Sprintf("**分析要求：** %s\n\n", contextNote)
			}
			
			prompt += "• 包含免責聲明\n\n"
		}
	}
	
	prompt += fmt.Sprintf("**用戶問題：** %s", message)
	
	return prompt
}
