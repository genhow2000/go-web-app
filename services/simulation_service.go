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
		questionType, _ := queryInstructions["question_type"].(string)
		
		if shouldQuery {
			response += "🔍 **外部資訊搜尋結果：**\n"
			response += "（模擬搜尋台灣證交所、Yahoo Finance、鉅亨網等資料源）\n\n"
			
			// 根據問題類型提供專門的分析
			switch questionType {
			case "investment_advice":
				response += s.generateInvestmentAdviceAnalysis(currentPrice, change)
			case "technical_analysis":
				response += s.generateTechnicalAnalysisDetails(currentPrice)
			case "risk_analysis":
				response += s.generateRiskAnalysisDetails(currentPrice, change)
			case "fundamental_analysis":
				response += s.generateFundamentalAnalysisDetails(currentPrice)
			default:
				// 預設綜合分析
				response += "**歷史股價分析：**\n"
				response += s.generateHistoricalAnalysis(currentPrice)
				response += "**技術指標分析：**\n"
				response += s.generateTechnicalIndicators(currentPrice)
				response += "**支撐位與阻力位：**\n"
				response += s.generateSupportResistance(currentPrice)
			}
		}
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

// generateInvestmentAdviceAnalysis 生成投資建議分析
func (s *SimulationService) generateInvestmentAdviceAnalysis(currentPrice, change float64) string {
	response := "**💡 投資建議分析：**\n\n"
	
	// 基本面分析
	response += "**基本面分析：**\n"
	response += "• 營收成長率：年增8.5%，表現優於同業平均\n"
	response += "• 獲利能力：毛利率32%，營業利益率15%，財務結構穩健\n"
	response += "• 估值水平：本益比18倍，股價淨值比2.1倍，估值合理\n"
	response += "• 產業地位：在相關領域具有技術領先優勢\n\n"
	
	// 技術面分析
	response += "**技術面分析：**\n"
	if change > 0 {
		response += "• 股價突破重要阻力位，技術面轉強\n"
		response += "• 成交量配合上漲，資金流入明顯\n"
		response += "• 短期均線呈多頭排列，趨勢向上\n"
	} else {
		response += "• 股價回調至支撐位附近，技術面偏弱\n"
		response += "• 成交量萎縮，觀望氣氛濃厚\n"
		response += "• 需觀察是否能在支撐位止跌回穩\n"
	}
	response += "\n"
	
	// 投資建議
	response += "**投資建議：**\n"
	if change > 0 {
		response += "• 短期：技術面偏多，可考慮逢低布局\n"
		response += "• 中期：基本面支撐，建議分批進場\n"
		response += "• 風險控制：設置停損點在支撐位下方5%\n"
		response += "• 目標價位：技術面目標價約為現價的110-115%\n"
	} else {
		response += "• 短期：技術面偏弱，建議觀望等待\n"
		response += "• 中期：基本面仍佳，可等待技術面轉強信號\n"
		response += "• 風險控制：跌破支撐位應考慮減碼\n"
		response += "• 進場時機：等待股價站穩支撐位後再考慮\n"
	}
	response += "\n"
	
	return response
}

// generateTechnicalAnalysisDetails 生成技術分析詳細內容
func (s *SimulationService) generateTechnicalAnalysisDetails(currentPrice float64) string {
	response := "**📈 技術指標詳細分析：**\n\n"
	
	// RSI分析
	rsi := 50.0 + (currentPrice-200)/10
	if rsi > 80 { rsi = 80 }
	if rsi < 20 { rsi = 20 }
	
	response += "**RSI相對強弱指標：**\n"
	response += fmt.Sprintf("• 當前RSI(14)：%.1f\n", rsi)
	if rsi > 70 {
		response += "• 信號：超買區域，需注意回調風險\n"
	} else if rsi < 30 {
		response += "• 信號：超賣區域，可能出現反彈\n"
	} else {
		response += "• 信號：中性區域，趨勢持續性較好\n"
	}
	response += "\n"
	
	// MACD分析
	response += "**MACD動量指標：**\n"
	response += "• DIF線：0.85，DEA線：0.72\n"
	response += "• 信號：金叉形成，多頭動能增強\n"
	response += "• 柱狀圖：正值擴大，買盤力道強勁\n\n"
	
	// KD分析
	kValue := rsi + 5
	dValue := rsi - 2
	response += "**KD隨機指標：**\n"
	response += fmt.Sprintf("• K值：%.1f，D值：%.1f\n", kValue, dValue)
	if kValue > 80 && dValue > 80 {
		response += "• 信號：超買區域，短期可能回調\n"
	} else if kValue < 20 && dValue < 20 {
		response += "• 信號：超賣區域，短期可能反彈\n"
	} else {
		response += "• 信號：偏多訊號，趨勢向上\n"
	}
	response += "\n"
	
	// 移動平均線分析
	response += "**移動平均線系統：**\n"
	response += "• 5日均線：股價上方，短期趨勢向上\n"
	response += "• 10日均線：股價上方，中期趨勢向上\n"
	response += "• 20日均線：股價上方，長期趨勢向上\n"
	response += "• 均線排列：多頭排列，趨勢明確\n\n"
	
	// 布林帶分析
	response += "**布林帶通道：**\n"
	response += "• 上軌：股價接近上軌，需注意回調\n"
	response += "• 中軌：股價站穩中軌，趨勢偏多\n"
	response += "• 下軌：支撐位明確，下跌空間有限\n"
	response += "• 通道寬度：正常範圍，波動適中\n\n"
	
	// 成交量分析
	response += "**成交量指標：**\n"
	response += "• 量價關係：價漲量增，量價配合良好\n"
	response += "• 資金流向：主力資金持續流入\n"
	response += "• 換手率：3.2%，流動性充足\n"
	response += "• 量能趨勢：成交量逐步放大\n\n"
	
	return response
}

// generateRiskAnalysisDetails 生成風險分析詳細內容
func (s *SimulationService) generateRiskAnalysisDetails(currentPrice, change float64) string {
	response := "**⚠️ 風險評估詳細分析：**\n\n"
	
	// 股價波動性分析
	response += "**股價波動性分析：**\n"
	response += "• 日波動率：2.8%，屬於中等波動水平\n"
	response += "• Beta值：1.15，系統性風險略高於大盤\n"
	response += "• 最大回撤：過去一年最大回撤15.2%\n"
	response += "• 波動區間：股價在支撐阻力位間震盪\n\n"
	
	// 流動性風險
	response += "**流動性風險評估：**\n"
	response += "• 日均成交量：2.5萬張，流動性充足\n"
	response += "• 買賣價差：0.1%，交易成本較低\n"
	response += "• 大單交易：機構投資人參與度高\n"
	response += "• 流動性評級：A級，流動性風險較低\n\n"
	
	// 基本面風險
	response += "**基本面風險因子：**\n"
	response += "• 財務結構：負債比率45%，財務結構穩健\n"
	response += "• 獲利穩定性：過去四季獲利穩定成長\n"
	response += "• 現金流：經營現金流為正，資金充裕\n"
	response += "• 信用評級：BBB+，信用風險較低\n\n"
	
	// 市場風險
	response += "**市場風險分析：**\n"
	response += "• 系統性風險：受大盤影響程度中等\n"
	response += "• 利率風險：對利率變化敏感度較低\n"
	response += "• 匯率風險：外銷比重30%，匯率影響有限\n"
	response += "• 政策風險：受政策影響程度中等\n\n"
	
	// 公司特定風險
	response += "**公司特定風險：**\n"
	response += "• 管理層風險：管理層經驗豐富，治理良好\n"
	response += "• 技術風險：技術領先優勢明顯\n"
	response += "• 競爭風險：市場地位穩固，競爭優勢明顯\n"
	response += "• 營運風險：營運效率持續改善\n\n"
	
	// 風險控制建議
	response += "**風險控制建議：**\n"
	response += "• 停損點設定：建議設在支撐位下方3-5%\n"
	response += "• 倉位控制：建議單一持股不超過總資產10%\n"
	response += "• 分散投資：建議搭配其他產業股票\n"
	response += "• 定期檢視：建議每月檢視持股狀況\n\n"
	
	return response
}

// generateFundamentalAnalysisDetails 生成基本面分析詳細內容
func (s *SimulationService) generateFundamentalAnalysisDetails(currentPrice float64) string {
	response := "**📊 基本面分析詳細內容：**\n\n"
	
	// 財務報表分析
	response += "**財務報表分析：**\n"
	response += "• 營收：年增8.5%，成長動能穩健\n"
	response += "• 毛利率：32.1%，較去年同期提升1.2%\n"
	response += "• 營業利益率：15.3%，獲利能力持續改善\n"
	response += "• 淨利率：12.8%，稅後獲利表現優異\n\n"
	
	// 獲利能力分析
	response += "**獲利能力分析：**\n"
	response += "• ROE（股東權益報酬率）：18.5%，高於同業平均\n"
	response += "• ROA（總資產報酬率）：12.3%，資產運用效率佳\n"
	response += "• EPS（每股盈餘）：8.5元，較去年同期成長15%\n"
	response += "• 本益比：18倍，估值合理\n\n"
	
	// 成長性分析
	response += "**成長性分析：**\n"
	response += "• 營收成長：過去三年複合成長率12%\n"
	response += "• 獲利成長：過去三年複合成長率15%\n"
	response += "• 未來展望：預估明年營收成長10-15%\n"
	response += "• 成長動能：新產品線貢獻持續增加\n\n"
	
	// 財務結構分析
	response += "**財務結構分析：**\n"
	response += "• 負債比率：45%，財務結構穩健\n"
	response += "• 流動比率：2.1，短期償債能力佳\n"
	response += "• 速動比率：1.8，流動性充足\n"
	response += "• 利息保障倍數：8.5倍，利息負擔輕微\n\n"
	
	// 產業地位分析
	response += "**產業地位和競爭優勢：**\n"
	response += "• 市場地位：在相關領域市占率排名前三\n"
	response += "• 技術優勢：擁有核心技術專利，進入門檻高\n"
	response += "• 品牌價值：品牌知名度高，客戶忠誠度佳\n"
	response += "• 成本優勢：規模經濟效應明顯，成本控制佳\n\n"
	
	// 管理層品質
	response += "**管理層品質和公司治理：**\n"
	response += "• 管理層經驗：平均從業經驗超過15年\n"
	response += "• 公司治理：董事會結構完善，獨立董事比例適當\n"
	response += "• 資訊透明度：財務資訊揭露完整，透明度高\n"
	response += "• 股東權益：重視股東權益，股利政策穩定\n\n"
	
	// 未來展望
	response += "**未來展望和成長動能：**\n"
	response += "• 產業趨勢：所屬產業前景看好，需求穩定成長\n"
	response += "• 新產品：預計明年推出3款新產品\n"
	response += "• 市場擴張：計劃進軍東南亞市場\n"
	response += "• 研發投入：研發費用占營收比重8%，持續創新\n\n"
	
	return response
}