package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// SimulationService 模拟AI服务
type SimulationService struct {
	usageStats AIUsageStats
}

// NewSimulationService 创建模拟服务
func NewSimulationService() *SimulationService {
	return &SimulationService{
		usageStats: AIUsageStats{
			Provider:    "simulation",
			DailyUsage:  0,
			DailyLimit:  999999,
			LastReset:   time.Now(),
			IsExhausted: false,
			ErrorCount:  0,
			LastError:   "",
			LastUsed:    time.Time{},
		},
	}
}

// GenerateResponse 生成模拟回复
func (s *SimulationService) GenerateResponse(ctx context.Context, message, conversationID string, stockContext map[string]interface{}) (string, error) {
	// 更新使用统计
	s.usageStats.DailyUsage++
	s.usageStats.LastUsed = time.Now()

	// 如果有股票上下文，生成專業的股票分析回复
	if stockContext != nil {
		response := s.generateStockAnalysisResponse(message, stockContext)
		return response, nil
	}

	// 生成模拟回复
	response := s.generateSimulatedResponse(message)
	return response, nil
}

// GetServiceName 获取服务名称
func (s *SimulationService) GetServiceName() string {
	return "Simulation Service"
}

// IsAvailable 检查服务是否可用
func (s *SimulationService) IsAvailable(ctx context.Context) bool {
	return !s.usageStats.IsExhausted
}

// GetUsageStats 获取使用统计
func (s *SimulationService) GetUsageStats() map[string]interface{} {
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

// generateSimulatedResponse 生成模拟回复
func (s *SimulationService) generateSimulatedResponse(message string) string {
	message = strings.ToLower(message)
	
	// 关键词匹配回复
	responses := map[string][]string{
		"你好": {
			"您好！我是阿和商城的AI購物助手，很高興為您服務！有什麼可以幫助您的嗎？",
			"你好！歡迎來到阿和商城，我是您的專屬購物助手，有什麼需要幫助的嗎？",
			"您好！很高興見到您，我是AI助手，隨時為您提供購物建議和幫助！",
		},
		"產品": {
			"我們有很多優質商品可以滿足您的需求，讓我為您介紹一下。",
			"我們提供各種類型的產品，包括電子產品、服裝、家居用品等，您對哪類產品感興趣？",
			"我們有豐富的產品選擇，從時尚服飾到智能家電，應有盡有！",
		},
		"優惠": {
			"我們經常推出各種優惠活動，包括折扣、滿減、贈品等，請關注我們的促銷信息！",
			"目前有多個優惠活動正在進行中，您可以查看我們的優惠頁面了解更多詳情。",
			"我們會定期推出限時優惠，建議您關注我們的官方通知！",
		},
		"配送": {
			"我們提供快速配送服務，一般情況下1-3個工作日內送達。",
			"配送時間根據您的地區而定，我們會盡快為您安排發貨。",
			"我們與多家物流公司合作，確保您的訂單安全快速送達。",
		},
		"退換": {
			"我們提供7天無理由退換貨服務，讓您購物更放心。",
			"如果商品有問題，請聯繫我們的客服，我們會為您處理退換貨事宜。",
			"我們有完善的售後服務，確保您的購物體驗。",
		},
		"推薦": {
			"根據您的需求，我推薦以下熱銷商品：智能手機、筆記本電腦、時尚服飾等。",
			"我們有很多熱門商品，包括最新款的電子產品和時尚單品。",
			"讓我為您推薦一些性價比很高的商品，相信您會喜歡的！",
		},
	}
	
	// 查找匹配的關鍵詞
	for keyword, responseList := range responses {
		if strings.Contains(message, keyword) {
			rand.Seed(time.Now().UnixNano())
			return responseList[rand.Intn(len(responseList))]
		}
	}
	
	// 默認回復
	defaultResponses := []string{
		"感謝您的詢問！我是阿和商城的AI助手，很樂意為您提供幫助。",
		"您好！我是購物助手，有什麼關於商品的問題都可以問我。",
		"很高興為您服務！請告訴我您需要什麼幫助，我會盡力協助您。",
		"我是阿和商城的智能助手，隨時為您提供購物建議和產品信息。",
		"歡迎來到阿和商城！我是您的專屬購物助手，有什麼可以幫助您的嗎？",
	}
	
	rand.Seed(time.Now().UnixNano())
	return defaultResponses[rand.Intn(len(defaultResponses))]
}

// generateStockAnalysisResponse 生成股票分析回复
func (s *SimulationService) generateStockAnalysisResponse(message string, stockContext map[string]interface{}) string {
	// 提取股票基本信息
	stockCode, stockName, market, currentPrice, change := extractStockInfo(stockContext)
	
	// 檢查是否有查詢指令
	queryInstructions, hasInstructions := stockContext["query_instructions"].(map[string]interface{})
	
	// 構建專業的股票分析回复
	response := "📊 **股票分析報告**\n\n"
	
	// 基本資訊
	response += "**股票資訊：**\n"
	response += fmt.Sprintf("• 股票代碼：%s\n", stockCode)
	response += fmt.Sprintf("• 股票名稱：%s\n", stockName)
	response += fmt.Sprintf("• 現價：%.2f 元\n", currentPrice)
	response += fmt.Sprintf("• 漲跌：%.2f 元\n", change)
	response += fmt.Sprintf("• 市場：%s\n\n", market)
	
	// 如果有查詢指令，模擬搜尋外部資訊
	if hasInstructions {
		shouldQuery, _ := queryInstructions["should_query_history"].(bool)
		if shouldQuery {
			response += "🔍 **外部資訊搜尋結果：**\n"
			response += "（模擬搜尋台灣證交所、Yahoo Finance、鉅亨網等資料源）\n\n"
			
			// 模擬搜尋到的歷史數據
			response += "**歷史股價分析：**\n"
			response += s.generateHistoricalAnalysis(currentPrice)
			
			// 模擬技術指標分析
			response += "**技術指標分析：**\n"
			response += s.generateTechnicalIndicators(currentPrice)
			
			// 模擬支撐阻力分析
			response += "**支撐位與阻力位：**\n"
			response += s.generateSupportResistance(currentPrice)
		}
	}
	
	// 根據用戶問題提供具體分析
	message = strings.ToLower(message)
	if strings.Contains(message, "技術指標") || strings.Contains(message, "技術分析") {
		response += "**技術指標詳細分析：**\n"
		response += "• RSI相對強弱指標顯示股價處於健康上漲狀態\n"
		response += "• MACD指標出現金叉，短期動能轉強\n"
		response += "• 布林帶顯示股價接近上軌，需注意回調風險\n"
		response += "• 成交量配合股價上漲，量價關係良好\n\n"
	} else if strings.Contains(message, "投資建議") || strings.Contains(message, "值得買") {
		response += "**投資建議：**\n"
		response += "• 短期：技術面偏多，可考慮逢低布局\n"
		response += "• 中期：需關注基本面變化，建議分批進場\n"
		response += "• 風險：注意市場波動，設置停損點\n"
		response += "• 建議：可參考專業分析師報告，謹慎投資\n\n"
	} else if strings.Contains(message, "風險") || strings.Contains(message, "風險評估") {
		response += "**風險評估：**\n"
		response += "• 市場風險：整體市場波動可能影響個股表現\n"
		response += "• 流動性風險：成交量需持續關注\n"
		response += "• 基本面風險：需關注公司財報和產業動態\n"
		response += "• 技術面風險：股價接近阻力位，可能面臨回調\n\n"
	}
	
	// 添加免責聲明
	response += "⚠️ **免責聲明：**\n"
	response += "以上分析僅供參考，不構成投資建議。投資有風險，入市需謹慎。\n"
	response += "建議投資前諮詢專業理財顧問，並充分了解相關風險。"
	
	return response
}

// generateHistoricalAnalysis 生成歷史分析
func (s *SimulationService) generateHistoricalAnalysis(currentPrice float64) string {
	// 基於當前價格生成動態分析
	basePrice := currentPrice
	if basePrice == 0 {
		basePrice = 200 // 預設價格
	}
	
	weekChange := basePrice * 0.02   // 1週變化約2%
	monthChange := basePrice * 0.05  // 1個月變化約5%
	yearRange := basePrice * 0.3     // 1年波動範圍約30%
	
	return fmt.Sprintf("• 最近1週：股價呈現上漲趨勢，成交量放大 (%.2f元)\n", basePrice+weekChange) +
		fmt.Sprintf("• 最近1個月：整體表現優於大盤，技術指標偏多 (%.2f元)\n", basePrice+monthChange) +
		fmt.Sprintf("• 最近3個月：股價在支撐位附近震盪，等待突破\n") +
		fmt.Sprintf("• 最近1年：股價波動範圍在 %.0f-%.0f 元之間\n\n", basePrice-yearRange, basePrice+yearRange)
}

// generateTechnicalIndicators 生成技術指標
func (s *SimulationService) generateTechnicalIndicators(currentPrice float64) string {
	// 基於當前價格生成動態技術指標
	rsi := 50.0 + (currentPrice-200)/10 // RSI 基於價格變化
	if rsi > 80 { rsi = 80 }
	if rsi < 20 { rsi = 20 }
	
	kValue := rsi + 5
	dValue := rsi - 2
	
	return fmt.Sprintf("• RSI (14)：%.1f - 處於中性偏多區域\n", rsi) +
		fmt.Sprintf("• MACD：金叉形成，多頭動能增強\n") +
		fmt.Sprintf("• KD指標：K值 %.1f，D值 %.1f - 偏多訊號\n", kValue, dValue) +
		fmt.Sprintf("• 移動平均線：股價站上5日、10日、20日均線\n\n")
}

// generateSupportResistance 生成支撐阻力位
func (s *SimulationService) generateSupportResistance(currentPrice float64) string {
	// 基於當前價格生成動態支撐阻力位
	basePrice := currentPrice
	if basePrice == 0 {
		basePrice = 200
	}
	
	support1 := basePrice * 0.95
	support2 := basePrice * 0.90
	resistance1 := basePrice * 1.05
	resistance2 := basePrice * 1.10
	
	return fmt.Sprintf("• 支撐位：%.0f元、%.0f元\n", support1, support2) +
		fmt.Sprintf("• 阻力位：%.0f元、%.0f元\n", resistance1, resistance2) +
		fmt.Sprintf("• 當前股價接近阻力位，需觀察突破情況\n\n")
}